package logger

import (
	"fmt"
	"os"
	"path"

	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var levelMap = map[string]zapcore.Level{
	"debug": zap.DebugLevel,
	"info":  zap.InfoLevel,
	"warn":  zap.WarnLevel,
	"error": zap.ErrorLevel,
	"fatal": zap.FatalLevel,
	"panic": zap.PanicLevel,
}

func newEncoderConfig() zapcore.EncoderConfig {
	return zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "name",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalLevelEncoder, // 大写编码器
		EncodeTime:     zapcore.ISO8601TimeEncoder,  // ISO8601 UTC 时间格式
		EncodeDuration: zapcore.SecondsDurationEncoder,
	}
}

type zapLogger struct {
	config      *Config
	logger      *zap.SugaredLogger
	atomicLevel zap.AtomicLevel
	skip        int // 额外的Skip
}

// 构造日志
func NewLogger(config *Config) Driver {
	return &zapLogger{
		config:      config,
		atomicLevel: zap.NewAtomicLevel(),
	}
}

// ZapLogger
func (zl *zapLogger) ZapLogger() *zap.Logger {
	return zl.logger.Desugar()
}

// 打开日志
func (zl *zapLogger) Open() error {
	var cores []zapcore.Core

	zl.atomicLevel = zap.NewAtomicLevelAt(zl.zapLevel())

	// Console logger
	{
		var encoder zapcore.Encoder

		// 日志编码
		if zl.config.Encoding == "json" {
			encoder = zapcore.NewJSONEncoder(newEncoderConfig())
		} else {
			// 输出颜色
			config := newEncoderConfig()
			if zl.config.Color {
				config.EncodeLevel = zapcore.CapitalColorLevelEncoder
			}
			encoder = zapcore.NewConsoleEncoder(config)
		}

		// 采用缓写模式
		if zl.config.Buffer {
			cores = append(cores, zapcore.NewCore(encoder, zapcore.AddSync(&zapcore.BufferedWriteSyncer{WS: zapcore.AddSync(os.Stdout)}), zl.atomicLevel))
		} else {
			cores = append(cores, zapcore.NewCore(encoder, zapcore.AddSync(os.Stdout), zl.atomicLevel))
		}
	}

	// File logger
	if zl.config.Path != "" {
		pathName, err := CheckLogPath(zl.config.Path)
		if err != nil {
			return err
		}

		// 文件默认采用JSON编码器
		encoder := zapcore.NewJSONEncoder(newEncoderConfig())

		// 日志写入器
		lumberJackLogger := &lumberjack.Logger{
			Filename:   path.Join(pathName, fmt.Sprintf("%v.log", zl.config.Name)), // 日志文件的位置
			MaxSize:    zl.config.MaxSize,                                          // 在进行切割之前，日志文件的最大大小（以MB为单位）
			MaxBackups: zl.config.MaxBackups,                                       // 保留旧文件的最大个数
			MaxAge:     zl.config.MaxAge,                                           // 保留旧文件的最大天数
			Compress:   false,                                                      // 是否压缩/归档旧文件
		}

		// 采用缓写模式
		if zl.config.Buffer {
			cores = append(cores, zapcore.NewCore(encoder, zapcore.AddSync(&zapcore.BufferedWriteSyncer{WS: zapcore.AddSync(lumberJackLogger)}), zl.atomicLevel))
		} else {
			cores = append(cores, zapcore.NewCore(encoder, zapcore.AddSync(lumberJackLogger), zl.atomicLevel))
		}
	}

	zl.logger = zap.New(zapcore.NewTee(cores...)).Named(zl.config.Name).Sugar()
	return nil
}

// 获取日志等级
func (zl *zapLogger) Level() string {
	return zl.atomicLevel.String()
}

// 设置日志等级
func (zl *zapLogger) SetLevel(val string) {
	level, ok := levelMap[val]
	if !ok {
		return
	}

	zl.atomicLevel.SetLevel(level)
}

// 关闭
func (zl *zapLogger) Close() {
	_ = zl.logger.Sync()
}

func (zl *zapLogger) zapLevel() zapcore.Level {
	level, ok := levelMap[zl.config.Level]
	if ok {
		return level
	}
	return zap.DebugLevel
}

// 构造日志
func (zl *zapLogger) Named(name string) Driver {
	return &zapLogger{
		config:      zl.config,
		logger:      zl.logger.Named(name),
		atomicLevel: zl.atomicLevel,
	}
}

// 构造日志
func (zl *zapLogger) With(key string, val interface{}) Driver {
	return &zapLogger{
		config:      zl.config,
		logger:      zl.logger.With(key, val),
		atomicLevel: zl.atomicLevel,
	}
}

// Sugar 根据指定的skip
func (zl *zapLogger) Sugar(skip int) Driver {
	l := &zapLogger{
		config:      zl.config,
		logger:      zl.logger,
		atomicLevel: zl.atomicLevel,
		skip:        skip + zl.skip,
	}

	return l
}

// 写日志
func (zl *zapLogger) Log(level Level, args ...interface{}) {
	// 对于Debug日志，提前做Level判断
	if level == DebugLevel && !zl.atomicLevel.Enabled(zap.DebugLevel) {
		return
	}

	args, fields := ParseFields(args)
	zl.logwf(level, fields, "", args...)
}

// 写日志
func (zl *zapLogger) Logf(level Level, format string, args ...interface{}) {
	// 对于Debug日志，提前做Level判断
	if level == DebugLevel && !zl.atomicLevel.Enabled(zap.DebugLevel) {
		return
	}

	args, fields := ParseFields(args)
	zl.logwf(level, fields, format, args...)
}

// 无模板的日志
func (zl *zapLogger) logwf(level Level, fields *Fields, format string, args ...interface{}) {
	var msg string
	if format == "" {
		msg = fmt.Sprint(args...)
	} else {
		msg = fmt.Sprintf(format, args...)
	}

	values := fields.Values()
	// 添加Caller
	if !fields.HasFlag(FlagNoCaller) {
		values = append(values, zl.String("caller", GetCaller(4+zl.skip+fields.GetSkip())))
	}

	switch level {
	case DebugLevel:
		zl.logger.Debugw(msg, values...)
	case InfoLevel:
		zl.logger.Infow(msg, values...)
	case WarnLevel:
		zl.logger.Warnw(msg, values...)
	case ErrorLevel:
		zl.logger.Errorw(msg, values...)
	case FatalLevel:
		zl.logger.Fatalw(msg, values...)
	case PanicLevel:
		zl.logger.Panicw(msg, values...)
	}
}

// Bool
func (zl *zapLogger) Bool(key string, val bool) interface{} {
	return zap.Bool(key, val)
}

// String
func (zl *zapLogger) String(key string, val string) interface{} {
	return zap.String(key, val)
}

// Int
func (zl *zapLogger) Int(key string, val int) interface{} {
	return zap.Int(key, val)
}

// Int32
func (zl *zapLogger) Int32(key string, val int32) interface{} {
	return zap.Int32(key, val)
}

// Int64
func (zl *zapLogger) Int64(key string, val int64) interface{} {
	return zap.Int64(key, val)
}

// Uint
func (zl *zapLogger) Uint(key string, val uint) interface{} {
	return zap.Uint(key, val)
}

// Uint32
func (zl *zapLogger) Uint32(key string, val uint32) interface{} {
	return zap.Uint32(key, val)
}

// Uint64
func (zl *zapLogger) Uint64(key string, val uint64) interface{} {
	return zap.Uint64(key, val)
}

// Float64
func (zl *zapLogger) Float64(key string, val float64) interface{} {
	return zap.Float64(key, val)
}
