package runtime

import (
	"errors"
	"runtime"
	"strings"
)

// 对标准库 `runtime.Caller` 的封装
func Caller(depth int) (info CallerInfo) {
	info.depth = depth + 1
	info.pc, info.File, info.Line, info.ok = runtime.Caller(info.depth)

	if info.ok {
		info.Method, info.Package = split(runtime.FuncForPC(info.pc).Name())
	}
	info.File, info.Path = split(info.File)
	return
}
func split(s string) (string, string) {
	i := strings.LastIndex(s, "/")
	if i > 0 {
		return s[i+1:], s[:i]
	}
	return s, ""
}

// 函数调用的名称等简单信息
type CallerInfo struct {
	Path    string
	File    string
	Package string
	Method  string
	Line    int

	depth int
	pc    uintptr
	ok    bool
}

type wrappedCallLastError struct {
	CallerInfo
	Err error
}

func (err *wrappedCallLastError) Unwrap() error {
	return err.Err
}
func (err *wrappedCallLastError) Error() string {
	return err.Err.Error()
}

// 接口: 获取调用栈的最后
func (err *wrappedCallLastError) GetCallLastInfo() CallerInfo {
	return err.CallerInfo
}

func WrapCallLast(err error, depth int, forceWrap ...bool) error {
	if err == nil {
		panic("err cannot nil")
	}
	b := false
	if i := len(forceWrap); i > 0 {
		b = forceWrap[i-1]
	}
	var ex *wrappedCallLastError
	if !b && errors.As(err, &ex) {
		return err
	}
	ex = &wrappedCallLastError{
		Err:        err,
		CallerInfo: Caller(depth + 1),
	}

	return ex
}

func GetCallLastFromError(err error) (info CallerInfo, ok bool) {
	var ex *wrappedCallLastError
	if errors.As(err, &ex) {
		info = ex.CallerInfo
		ok = true
	}
	return
}

type wrappedCallStackError struct {
	Stack []CallerInfo
	Err   error
}

func (err *wrappedCallStackError) Unwrap() error {
	return err.Err
}
func (err *wrappedCallStackError) Error() string {
	return err.Err.Error()
}

// 接口: 获取调用栈的最后
func (err *wrappedCallStackError) GetCallStack() []CallerInfo {
	return err.Stack
}

func WrapCallStacktrace(err error, depth int) error {
	if err == nil {
		return nil
	}

	var ex *wrappedCallStackError
	if !errors.As(err, &ex) {
		ex = &wrappedCallStackError{
			Err:   err,
			Stack: []CallerInfo{},
		}
	}
	ex.Stack = append(ex.Stack, Caller(depth+1))
	return ex
}

func GetCallStackFromError(err error) (info []CallerInfo, ok bool) {
	var ex *wrappedCallStackError
	if errors.As(err, &ex) {
		info = ex.Stack
		ok = true
	}
	return
}
