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

// 安全关闭资源。如果有返回错误则记录, 如果fn返回错误将记录并返回 true, 否则返回false
// todo: 记录调用方行号
func Close(closer interface{}, log ...ErrorLogger) (b bool) {
	if closer == nil {
		return false
	}

	defer func() {
		if o := recover(); o != nil {
			b = logErr(fmt.Errorf("%v", o), log)
		}
	}()
	if c, _ := closer.(io.Closer); c != nil {
		b = logErr(c.Close(), log)
	} else if fn, _ := closer.(func() error); fn != nil {
		b = logErr(fn(), log)
	} else if fn, _ := closer.(func()); fn != nil {
		fn()
	} else {
		panic(fmt.Sprintf("暂不支持 closer 的类型: %T", closer))
	}
	return
}

func logErr(err error, log []ErrorLogger) bool {
	if err == nil {
		return false
	}
	var l ErrorLogger = Default()
	if len(log) > 0 {
		l = log[len(log)-1]
	}
	l.Error(err)
	return true
}
