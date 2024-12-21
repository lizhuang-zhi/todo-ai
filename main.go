package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"todo-ai/common"
	"todo-ai/common/consts"
	"todo-ai/core"
	"todo-ai/core/logger"
	"todo-ai/events"
	"todo-ai/router"

	"github.com/rs/xid"
	"github.com/spf13/pflag"
)

const srvName = "todo-ai"

var (
	LDFlagVersion   string = "dev" // 打包版本号
	LDFlagBuildTime string = ""    // 打包时间
	LDFlagGOVersion string = ""    // Golang版本号
)

// @title Ai Dev Mate
// @version 0.0.1
// @description This is a sample Server pets
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name x-token
// @BasePath /
func main() {
	// 定义命令行参数
	version := pflag.BoolP("version", "v", false, "print version")                                   // 查看版本信息
	configFile := pflag.StringP("config", "c", "./configs/local/config.yaml", "choose config file.") // 指定配置文件
	srvID := pflag.StringP("srvID", "s", "", "unique srvID")                                         // 指定SrvID
	pflag.Parse()

	// 输出版本信息
	if *version {
		fmt.Printf("Version=%s\n", LDFlagVersion)
		fmt.Printf("BuildTime=%s\n", LDFlagBuildTime)
		fmt.Printf("GoVersion=%s\n", LDFlagGOVersion)
		os.Exit(0)
	}

	// 初始化配置信息
	if err := common.InitConfig(*configFile); err != nil {
		panic(fmt.Errorf("init config err:%v", err))
	}

	err := NewServer().Start(*srvID)
	if err != nil {
		logger.Errorf("server run err:%v", err)
	}
}

type Server struct {
	ctx     context.Context
	closeFn func()
	server  core.Server
}

func NewServer() *Server {
	return &Server{
		server: core.NewServer("0.0.0.0:"+common.Config.System.TCPPort, router.Routers()),
	}
}

// Start Server
func (s *Server) Start(srvID string) error {
	// 如果没有指定SrvID，如果是K8S环境，读取环境变量，如果没有，则随机
	if srvID == "" {
		srvID = os.Getenv("K8S_POD_NAME")
		if srvID == "" {
			srvID = srvName + "-" + xid.New().String()
		}
	} else {
		srvID = srvName + "-" + srvID
	}

	// 初始化日志信息
	if err := logger.Open(&logger.Config{
		Name:       srvName,
		Level:      common.Config.Logging.Level,
		Path:       common.Config.Logging.Path,
		Encoding:   common.Config.Logging.Encoding,
		Color:      common.Config.Logging.Color,
		MaxSize:    common.Config.Logging.MaxSize,
		MaxBackups: common.Config.Logging.MaxBackups,
		MaxAge:     common.Config.Logging.MaxAge,
		Buffer:     common.Config.Logging.Buffer,
	}); err != nil {
		panic(fmt.Errorf("init log err:%v", err))
	}

	s.ctx, s.closeFn = context.WithCancel(context.Background())
	s.listenSignal()

	// start mongodb
	mongodbInstance := common.Config.MongoDB.Instance
	db, err := core.NewMongoDB(mongodbInstance)
	if err != nil {
		return err
	}
	common.Mgo = db
	if err = s.initDB(); err != nil {
		return err
	}

	if err = s.initUUID(); err != nil {
		return err
	}

	// 初始化Dify Workflow
	err = events.InitWorkflowCfg()
	if err != nil {
		panic(fmt.Errorf("init workflow cfg err:%v", err))
	}

	s.server.Start()
	logger.Infof("%v start, srvID: %v", srvName, srvID)
	<-s.ctx.Done()
	return nil
}

// db初始化
func (s *Server) initDB() error {
	if err := common.Mgo.CreateIndexes(context.Background(), consts.CollectionUser,
		core.NewIndex(true, "user_id"),
	); err != nil {
		return err
	}

	if err := common.Mgo.CreateIndexes(context.Background(), consts.CollectionTask,
		core.NewIndex(true, "task_id"),
	); err != nil {
		return err
	}

	return nil
}

func (s *Server) initUUID() error {
	common.UserUUID = core.NewUUID(common.Mgo, consts.CollectionCount, consts.UserUUID)
	// 设置初始值为1
	err := common.UserUUID.Init(1)
	if err != nil {
		return err
	}

	common.TaskUUID = core.NewUUID(common.Mgo, consts.CollectionCount, consts.TaskUUID)
	// 设置初始值为1
	err = common.TaskUUID.Init(1)
	if err != nil {
		return err
	}

	return nil
}

// Stop Server
func (s *Server) Stop() {
	s.server.Stop()
	s.closeFn()
}

// listenSignal
func (s *Server) listenSignal() {
	terminate := make(chan os.Signal, 1)
	signal.Notify(terminate, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		signal := <-terminate
		logger.Infof("Recv signal %v", signal)
		s.Stop()
	}()
}
