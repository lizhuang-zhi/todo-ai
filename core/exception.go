package core

// ExceptionOutPut 输出函数
var ExceptionOutPut func(err interface{})

// 恢复
func Recovery() {
	if err := recover(); err != nil {
		ExceptionOutPut(err)
	}
}

func GO(f func()) {
	go func() {
		defer Recovery()
		f()
	}()
}
