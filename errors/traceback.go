package errors

import (
	"fmt"
	"runtime"
)

func Caller(calldepth int) (funcName, file string, line int) {
	var pc uintptr
	var ok bool
	pc, file, line, ok = runtime.Caller(calldepth + 1)

	if ok {
		funcName = runtime.FuncForPC(pc).Name()
	}
	return
}

type TracebackError struct {
	Err    error
	Method string
	File   string
	Line   int
}

func (ex *TracebackError) Unwrap() error {
	return ex.Err
}
func (ex *TracebackError) Error() string {
	if ex.Err == nil {
		return fmt.Sprintf("TracebackError: %s:%d", ex.File, ex.Line)
	}
	return ex.Err.Error()
}

func WithTraceback(err error, forceWrap ...bool) error {
	if err == nil {
		return nil
	}
	force := false
	for _, it := range forceWrap {
		force = it
	}
	if !force {
		v := AsTraceback(err)
		if v != nil {
			return v
		}
	}
	method, file, line := Caller(1)

	return &TracebackError{
		Method: method,
		File:   file,
		Line:   line,
		Err:    err,
	}
}

func AsTraceback(err error) *TracebackError {
	if err == nil {
		return nil
	}
	var v *TracebackError
	As(err, &v)
	return v
}
