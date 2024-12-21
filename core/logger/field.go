package logger

import (
	"runtime"
	"strconv"
	"strings"
)

const (
	FlagNoCaller uint8 = 1 << iota // 不需要Caller
)

// EmptyFields 空Fields
var EmptyFields = &Fields{}

// 日志字段
type Fields struct {
	flags  uint8 // 标记位
	skip   int   // Caller跳过
	values []interface{}
}

// 构造字段
func NewFields(values ...interface{}) *Fields {
	return &Fields{
		values: values,
	}
}

// 添加字段
func (f *Fields) AddField(field interface{}) {
	f.values = append(f.values, field)
}

// 添加多个字段
func (f *Fields) AddFields(fields ...interface{}) {
	f.values = append(f.values, fields...)
}

// 添加KV
func (f *Fields) Add(key string, val interface{}) {
	f.values = append(f.values, key, val)
}

// 重置
func (f *Fields) Reset() {
	f.values = f.values[:0]
	f.flags = 0
	f.skip = 0
}

// 获取数据
func (f *Fields) Values() []interface{} {
	return f.values
}

// 添加标记
func (f *Fields) AddFlag(flag uint8) {
	f.flags |= flag
}

// 是否包含标记
func (f *Fields) HasFlag(flag uint8) bool {
	return f.flags&flag > 0
}

// 设置Caller跳过
func (f *Fields) SetSkip(skip int) {
	f.skip = skip
}

// 获取Caller跳过
func (f *Fields) GetSkip() int {
	return f.skip
}

// Clone
func (f *Fields) Clone() *Fields {
	new := &Fields{values: make([]interface{}, len(f.values))}
	copy(new.values, f.values)
	return new
}

// parseFields 解析参数中的Fields
func ParseFields(args []interface{}) ([]interface{}, *Fields) {
	if len(args) > 0 {
		tail := len(args) - 1
		if fields, ok := args[tail].(*Fields); ok {
			return args[:tail], fields
		}
	}

	return args, EmptyFields
}

// 获取调用者(文件:行号)
func GetCaller(skip int) string {
	_, file, line, ok := runtime.Caller(skip)
	if !ok {
		return ""
	}

	idx := strings.LastIndexByte(file, '/')
	if idx == -1 {
		return file + ":" + strconv.Itoa(line)
	}

	idx = strings.LastIndexByte(file[:idx], '/')
	if idx == -1 {
		return file + ":" + strconv.Itoa(line)
	}

	return file[idx+1:] + ":" + strconv.Itoa(line)
}
