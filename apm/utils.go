package apm

import (
	"fmt"
	"io"
)

// 如果 err 不为 nil 则包装错误并panic
func PanicIf(err error) {
	if err == nil {
		return
	}
	panic(err)
}

type ErrorLogger interface {
	Error(...interface{})
}

// 安全关闭资源。如果有返回错误则记录, 如果fn返回错误将记录并返回 false, 否则返回true
func Close(closer interface{}, log ...ErrorLogger) bool {
	if closer == nil {
		return false
	}
	var err error
	if c, _ := closer.(io.Closer); c != nil {
		err = c.Close()
	} else if fn, _ := closer.(func() error); fn != nil {
		err = fn()
	} else {
		panic(fmt.Sprintf("暂不支持 closer 的类型: %T", closer))
	}

	if err != nil {
		var l ErrorLogger = Default()
		if len(log) > 0 {
			l = log[len(log)-1]
		}
		l.Error(err)
		return false
	}

	return true

}
