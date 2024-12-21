package shttp

import (
	"strings"
	"todo-ai/core/logger"

	"github.com/gin-gonic/gin"
)

// 日志
type LoggerWriter struct {
	Formatter gin.LogFormatter
}

// 构造日志
func NewLogger(tag string) *LoggerWriter {
	return &LoggerWriter{Formatter: newFormatter(tag)}
}

// 写入日志
func (l LoggerWriter) Write(msg []byte) (n int, err error) {
	if len(msg) == 0 {
		return 0, nil
	}

	lv := logger.InfoLevel // 默认INFO日志

	s := string(msg)

	// 根据内容确定日志等级
	if strings.Contains(s, "panic") || strings.Contains(s, "PANIC") {
		lv = logger.ErrorLevel
	} else if strings.Contains(s, "error") || strings.Contains(s, "ERROR") {
		lv = logger.ErrorLevel
	} else if strings.Contains(s, "warn") || strings.Contains(s, "WARN") {
		lv = logger.WarnLevel
	} else if strings.Contains(s, "debug") || strings.Contains(s, "DEBUG") {
		lv = logger.DebugLevel
	}

	switch lv {
	case logger.DebugLevel:
		logger.Debug(s)
	case logger.InfoLevel:
		logger.Info(s)
	case logger.WarnLevel:
		logger.Warn(s)
	case logger.ErrorLevel, logger.FatalLevel, logger.PanicLevel:
		logger.Error(s)
	}
	return len(s), nil
}

func newFormatter(tag string) gin.LogFormatter {
	return func(param gin.LogFormatterParams) string {
		fields := logger.NewFields()
		fields.AddFields(logger.String(logger.LogFieldIP, param.ClientIP),
			logger.Int64(logger.LogFieldCost, param.Latency.Milliseconds()))

		fields.AddFlag(logger.FlagNoCaller)

		var err string

		if param.ErrorMessage != "" {
			err = param.ErrorMessage
		}

		if err != "" {
			logger.Warnf("[%v] %3d %-7s %v, failed: %v", tag, param.StatusCode, param.Method, param.Path, err, fields)
		} else {
			logger.Infof("[%v] %3d %-7s %v", tag, param.StatusCode, param.Method, param.Path, fields)
		}

		return ""
	}
}
