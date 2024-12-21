package common

import (
	"fmt"
	"os"
	"path"
	"todo-ai/core/utils"

	"github.com/spf13/viper"
)

// 初始化配置信息
func InitConfig(configFile string) error {
	v := viper.New()
	v.SetConfigFile(configFile)
	err := v.ReadInConfig()
	if err != nil {
		return fmt.Errorf("fatal error config file: %s", err)
	}
	if err := v.Unmarshal(&Config); err != nil {
		return err
	}

	if Config.System.Server == "" {
		Config.System.Server = utils.GetLocalIP()
	}
	return nil
}

// 初始化临时缓存目录
func InitTmpDir() error {
	tmpDir := path.Clean(Config.System.TmpDir)
	_, err := os.Stat(tmpDir)
	if !(err == nil || os.IsExist(err)) {
		if err := os.MkdirAll(tmpDir, 0777); err != nil { //os.ModePerm
			return err
		}
	}
	return nil
}
