package apm

import (
	"fmt"
	"io"
	"runtime/debug"
)

// 如果 err 不为 nil 则包装错误并panic
func PanicIf(err error) {
	if err == nil {
		return
	}
	panic(err)
}

type closeLogger interface {
	Error(...interface{})
	Debug(...interface{})
}

// 安全关闭资源。如果有返回错误则记录, 如果fn返回错误将记录并返回 true, 否则返回false
// todo: 记录调用方行号
func Close(closer interface{}, log ...closeLogger) (b bool) {
	if closer == nil {
		return false
	}

	defer func() {
		if o := recover(); o != nil {
			b = logErr(fmt.Errorf("%v", o), log, false)
			logErr(fmt.Errorf("recovered: %s", debug.Stack()), log, false)
		}
	}()

	if c, _ := closer.(io.Closer); c != nil {
		return logErr(c.Close(), log, false)

	} else if fn, _ := closer.(func()); fn != nil {
		fn()
		return false
	} else if fn, _ := closer.(func() error); fn != nil {
		return logErr(fn(), log, false)
	}

	return logErr(fmt.Errorf("暂不支持 closer 的类型: %T", closer), log, false)
}

func logErr(err error, log []closeLogger, debug bool) bool {
	if err == nil {
		return false
	}
	var l closeLogger = Default()
	if len(log) > 0 {
		l = log[len(log)-1]
	}
	if debug {
		l.Debug(err)
		return true
	}
	l.Error(err)
	return true
}
