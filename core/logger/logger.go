package logger

import (
	"fmt"
	"os"
)

var defaultLoggerDriver Driver

// Open logger
func Open(config *Config) error {
	logDriver := NewLogger(config)

	if err := logDriver.Open(); err != nil {
		return err
	}

	defaultLoggerDriver = logDriver
	return nil
}

// 关闭
func Close() {
	defaultLoggerDriver.Close()
}

// debug日志
func Debug(args ...interface{}) {
	defaultLoggerDriver.Log(DebugLevel, args...)
}

// info日志
func Info(args ...interface{}) {
	defaultLoggerDriver.Log(InfoLevel, args...)
}

// warn日志
func Warn(args ...interface{}) {
	defaultLoggerDriver.Log(WarnLevel, args...)
}

// 错误日志
func Error(args ...interface{}) {
	defaultLoggerDriver.Log(ErrorLevel, args...)
}

// fatal日志
func Fatal(args ...interface{}) {
	defaultLoggerDriver.Log(FatalLevel, args...)
	os.Exit(1)
}

// panic日志
func Panic(args ...interface{}) {
	defaultLoggerDriver.Log(PanicLevel, args...)
	panic(fmt.Sprint(args...))
}

// debug日志
func Debugf(format string, args ...interface{}) {
	defaultLoggerDriver.Logf(DebugLevel, format, args...)
}

// info日志
func Infof(format string, args ...interface{}) {
	defaultLoggerDriver.Logf(InfoLevel, format, args...)
}

// warn日志
func Warnf(format string, args ...interface{}) {
	defaultLoggerDriver.Logf(WarnLevel, format, args...)
}

// error日志
func Errorf(format string, args ...interface{}) {
	defaultLoggerDriver.Logf(ErrorLevel, format, args...)
}

// fatal日志
func Fatalf(format string, args ...interface{}) {
	defaultLoggerDriver.Logf(FatalLevel, format, args...)
	os.Exit(1)
}

// panic日志
func Panicf(format string, args ...interface{}) {
	defaultLoggerDriver.Logf(PanicLevel, format, args...)
	panic(fmt.Sprintf(format, args...))
}
