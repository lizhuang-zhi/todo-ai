package logger

type Driver interface {
	Open() error
	Close()
	Named(name string) Driver
	With(key string, val interface{}) Driver
	Log(level Level, args ...interface{})
	Logf(level Level, format string, args ...interface{})
	Sugar(skip int) Driver
	Level() string
	SetLevel(level string)

	// 额外字段
	Bool(key string, val bool) interface{}
	String(key string, val string) interface{}
	Int(key string, val int) interface{}
	Int32(key string, val int32) interface{}
	Int64(key string, val int64) interface{}
	Uint(key string, val uint) interface{}
	Uint32(key string, val uint32) interface{}
	Uint64(key string, val uint64) interface{}
	Float64(key string, val float64) interface{}
}
