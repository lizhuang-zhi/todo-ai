package logger

// 通用日志字段，都会设置索引
const (
	LogFieldID         = "id"          // 请求ID
	LogFieldCommand    = "command"     // 请求命令
	LogFieldIP         = "ip"          // 请求IP
	LogFieldRemoteAddr = "remove_addr" // 请求IP端口
	LogFieldCost       = "cost"        // 请求耗时
	LogFieldError      = "error"       // 请求错误
	LogFieldTrace      = "trace"       // 请求追踪ID
	LogFieldCaller     = "caller"      // 调用来源
	LogFieldSize       = "size"        // 请求大小
	LogFieldSource     = "source"      // 来源
	LogFieldRes        = "res"         // 返回信息
)

// Bool
func Bool(key string, val bool) interface{} {
	return defaultLoggerDriver.Bool(key, val)
}

// String
func String(key string, val string) interface{} {
	return defaultLoggerDriver.String(key, val)
}

// Int
func Int(key string, val int) interface{} {
	return defaultLoggerDriver.Int(key, val)
}

// Int32
func Int32(key string, val int32) interface{} {
	return defaultLoggerDriver.Int32(key, val)
}

// Int64
func Int64(key string, val int64) interface{} {
	return defaultLoggerDriver.Int64(key, val)
}

// Uint
func Uint(key string, val uint) interface{} {
	return defaultLoggerDriver.Uint(key, val)
}

// Uint32
func Uint32(key string, val uint32) interface{} {
	return defaultLoggerDriver.Uint32(key, val)
}

// Uint64
func Uint64(key string, val uint64) interface{} {
	return defaultLoggerDriver.Uint64(key, val)
}

// Float64
func Float64(key string, val float64) interface{} {
	return defaultLoggerDriver.Float64(key, val)
}
