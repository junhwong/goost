package runtime

import (
	"errors"
	"fmt"
	"net/http"
	"runtime"
)

// from k8s

var (
	// ReallyCrash controls the behavior of HandleCrash and now defaults
	// true. It's still exposed so components can optionally set to false
	// to restore prior behavior.
	ReallyCrash = true
)

// PanicHandlers is a list of functions which will be invoked when a panic happens.
var PanicHandlers = []func(error){logPanic}

// HandleCrash simply catches a crash and logs an error. Meant to be called via
// defer.  Additional context-specific handlers can be provided, and will be
// called in case of panic.  HandleCrash actually crashes, after calling the
// handlers and logging the panic message.
//
// E.g., you can provide one or more additional handlers for something like shutting down go routines gracefully.
func HandleCrash(handlers ...func(error)) {
	r := recover()
	if r == nil {
		return
	}
	err, _ := r.(error)
	if err == nil {
		err = &recoveredError{recovered: r}
		// TODO: 获取堆栈的首行
	}
	err = wrapTracebackLastError(err, 1, false)
	err = wrapTracebackStackError(err, 1, false)

	ok := false
	for _, fn := range handlers {
		if fn == nil {
			continue
		}
		ok = true
		fn(err)
	}
	if ok {
		return
	}

	// 默认处理
	for _, fn := range PanicHandlers {
		if fn == nil {
			continue
		}
		ok = true
		fn(err)
	}
	if !ok {
		logPanic(err)
	}
	if ReallyCrash {
		// Actually proceed to panic.
		// TODO: 有可能死循环
		panic(err)
	}

}

type recoveredError struct {
	recovered interface{}
}

func (err *recoveredError) RecoverValue() interface{} {
	return err.recovered
}

func (err *recoveredError) Error() string {
	return fmt.Sprintf("PanicError: %#v", err.recovered)
}

// logPanic logs the caller tree when a panic occurs (except in the special case of http.ErrAbortHandler).
func logPanic(r error) {
	if r == http.ErrAbortHandler {
		// honor the http.ErrAbortHandler sentinel panic value:
		//   ErrAbortHandler is a sentinel panic value to abort a handler.
		//   While any panic from ServeHTTP aborts the response to the client,
		//   panicking with ErrAbortHandler also suppresses logging of a stack trace to the server's error log.
		return
	}

	// Same as stdlib http server code. Manually allocate stack trace buffer size
	// to prevent excessively large logs
	const size = 64 << 10
	stacktrace := make([]byte, size)
	stacktrace = stacktrace[:runtime.Stack(stacktrace, false)]
	fmt.Printf("Observed a panic: %#v (%v)\n%s", r, r, stacktrace)
}

// PanicIf panics on non-nil errors. Useful to handling programmer level errors.
func PanicIf(err error) {
	if err == nil {
		return
	}
	err = wrapTracebackLastError(err, 1, false)
	err = wrapTracebackStackError(err, 1, false)
	panic(err)
}

func Caller(depth int) (funcName, file string, line int) {
	var pc uintptr
	var ok bool
	pc, file, line, ok = runtime.Caller(depth + 1)

	if ok {
		funcName = runtime.FuncForPC(pc).Name()
	}
	return
}
func wrapTracebackLastError(err error, depth int, f bool) (ex *TracebackLastError) {
	if err == nil {
		panic("err cannot nil")
	}
	if errors.As(err, &ex) && !f {
		return
	}
	ex = &TracebackLastError{
		Err:   err,
		depth: depth + 1,
	}
	ex.Method, ex.File, ex.Line = Caller(ex.depth)

	return
}

type TracebackLast struct {
	File   string
	Line   int
	Method string
}

type TracebackLastError struct {
	TracebackLast
	Err error

	depth int
}

func (err *TracebackLastError) Unwrap() error {
	return err.Err
}
func (err *TracebackLastError) Error() string {
	return err.Err.Error()
}

func wrapTracebackStackError(err error, depth int, f bool) (ex *TracebackStackError) {
	if err == nil {
		panic("err cannot nil")
	}
	if errors.As(err, &ex) && !f {
		return
	}
	ex = &TracebackStackError{
		Err:   err,
		depth: depth + 1,
	}
	const size = 64 << 10
	stacktrace := make([]byte, size)
	ex.Stack = stacktrace[:runtime.Stack(stacktrace, false)]
	return
}

type TracebackStackError struct {
	Err   error
	Stack []byte

	depth int
}

func (err *TracebackStackError) Unwrap() error {
	return err.Err
}
func (err *TracebackStackError) Error() string {
	return err.Err.Error()
}

// type StopCh = <-chan struct{}
