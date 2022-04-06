package apm

// 如果 err 不为 nil 则包装错误并panic
func PanicIf(err error) {
	if err == nil {
		return
	}
	panic(err)
}
