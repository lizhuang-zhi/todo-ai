package logger

import (
	"errors"
	"fmt"
	"os"
	"path"
	"strings"
)

// 配置
type Config struct {
	Name       string // 日志名称
	Level      string // 日志等级
	Path       string // 日志路径
	Encoding   string // 编码
	Color      bool   // 是否显示颜色
	Buffer     bool   // 是否采用缓写模式
	MaxSize    int    // 在进行切割之前，日志文件的最大大小（以MB为单位）
	MaxBackups int    // 保留旧文件的最大个数
	MaxAge     int    // 保留旧文件的最大天数
}

type Level uint32

const (
	PanicLevel Level = iota // 崩溃等级，会停止进程打印堆栈
	FatalLevel              // 致命等级，会停止进程
	ErrorLevel              // 错误等级
	WarnLevel               // 警告等级
	InfoLevel               // 通知等级
	DebugLevel              // 调试等级
)

// ToLevel
func ToLevel(level string) Level {
	switch strings.ToLower(level) {
	case "panic":
		return PanicLevel
	case "fatal":
		return FatalLevel
	case "error":
		return ErrorLevel
	case "warn":
		return WarnLevel
	case "info":
		return InfoLevel
	case "debug":
		return DebugLevel
	}
	return DebugLevel
}

// 日志等级
func (level Level) String() string {
	switch level {
	case PanicLevel:
		return "panic"
	case FatalLevel:
		return "fatal"
	case ErrorLevel:
		return "error"
	case WarnLevel:
		return "warn"
	case InfoLevel:
		return "info"
	case DebugLevel:
		return "debug"
	}
	return "Unknown"
}

// 检查日志路径
func CheckLogPath(logPath string) (string, error) {
	logPath = path.Clean(logPath)

	if logPath == "" {
		return "", errors.New("logPath is empty")
	}

	if !FileExists(logPath) {
		if err := os.MkdirAll(logPath, os.ModePerm); err != nil {
			return "", fmt.Errorf("checkLogPath with err: %w", err)
		}
	}

	return logPath, nil
}

// 文件是否存在
func FileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil || os.IsExist(err)
}
