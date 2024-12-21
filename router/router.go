package router

import (
	"fmt"
	"net/http"
	"sync"
	"todo-ai/core/logger"
	"todo-ai/core/shttp"
	"todo-ai/middleware"

	"github.com/gin-gonic/gin"
)

type (
	Handler   func(ctx *gin.Context) (data interface{}, err error)
	SLHandler func(ctx *gin.Context) (data string, err error)
)

var ginNewOnce sync.Once

// 初始化总路由
func Routers() *gin.Engine {
	writer := shttp.NewLogger("todo-ai")

	ginNewOnce.Do(func() {
		gin.SetMode(gin.ReleaseMode)
		gin.DefaultWriter = writer
		gin.DefaultErrorWriter = writer
	})

	Router := gin.New()

	Router.Use(gin.LoggerWithConfig(gin.LoggerConfig{Output: writer, Formatter: writer.Formatter}), middleware.Recovery())
	Router.Use(middleware.Cors())
	APIGroup := Router.Group("")

	InitTaskRouter(APIGroup) // 注册任务路由

	return Router
}

func WrapperHandler(handler Handler) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		data, err := handler(ctx)
		if err != nil {
			logger.Errorf("handle %s error:%s", ctx.Request.URL, err)
			ctx.JSON(http.StatusInternalServerError, err.Error())
			return
		}
		// value, exist := ctx.Get(auth.CtxAccessToken)
		// if exist {
		// 	if val, ok := value.(string); ok {
		// 		ctx.JSON(http.StatusOK, gin.H{"data": data, "newAccessToken": val})
		// 		return
		// 	}
		// }
		ctx.JSON(http.StatusOK, data)
	}
}

// WrapperSLAPIHandler
func WrapperSLAPIHandler(handler SLHandler) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		data, err := handler(ctx)
		if err != nil {
			logger.Errorf("slapi handle %s error:%s", ctx.Request.URL, err)
			ctx.String(http.StatusInternalServerError, err.Error())
			return
		}
		ctx.String(http.StatusOK, data)
	}
}

type JSONError struct {
	Status int    `json:"status"`
	Err    string `json:"err"`
}

func (e *JSONError) Error() error {
	return fmt.Errorf("[%v]%s", e.Status, e.Err)
}

type JSONReply struct {
	Status int    `json:"status"`
	Data   string `json:"data"`
}

func JSONWrap(handler Handler) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		data, err := handler(ctx)
		if err != nil {
			logger.Errorf("handle %s error:%s", ctx.Request.URL, err)
			ctx.JSON(http.StatusInternalServerError, &JSONError{
				Status: http.StatusInternalServerError,
				Err:    err.Error(),
			})
			return
		}
		if data == nil {
			ctx.JSON(http.StatusOK, &JSONReply{
				Status: http.StatusOK,
				Data:   "ok",
			})
			return
		}
		ctx.JSON(http.StatusOK, data)
	}
}
