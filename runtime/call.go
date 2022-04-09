package runtime

import (
	"errors"
	"runtime"
)

func Caller(depth int) (info CallSourceInfo) {
	var pc uintptr
	var ok bool
	pc, info.File, info.Line, ok = runtime.Caller(depth + 1)

	if ok {
		info.Method = runtime.FuncForPC(pc).Name()
	}
	return
}

// 调用源简单信息
type CallSourceInfo struct {
	File   string
	Line   int
	Method string
}

type wrappedCallLastError struct {
	CallSourceInfo
	Err error

	depth int
}

func (err *wrappedCallLastError) Unwrap() error {
	return err.Err
}
func (err *wrappedCallLastError) Error() string {
	return err.Err.Error()
}

// 接口: 获取调用栈的最后
func (err *wrappedCallLastError) GetCallLastInfo() CallSourceInfo {
	return err.CallSourceInfo
}

func WrapCallLast(err error, depth int, forceWrap ...bool) (ex *wrappedCallLastError) {
	if err == nil {
		panic("err cannot nil")
	}
	b := false
	if i := len(forceWrap); i > 0 {
		b = forceWrap[i-1]
	}
	if errors.As(err, &ex) && !b {
		return
	}
	ex = &wrappedCallLastError{
		Err:            err,
		depth:          depth + 1,
		CallSourceInfo: Caller(depth + 1),
	}

	return
}
