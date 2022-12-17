package apm

import (
	"bytes"
	"errors"
	"fmt"
	"runtime/debug"
	"strings"
	"sync"

	"github.com/junhwong/goost/apm/field"
	"github.com/spf13/cast"
)

// 标准错误
type CodeError struct {
	code   string // 固定的全局系统唯一的错误码, 如: NOTFOUND
	msg    string // 固定安全的错误消息, 用于对外部系统暂时. 如: 密码或账号不正确
	status int    // 用于如 HTTP/RPC 等接口返回的状态码, -1 表示未定义
}

func (err *CodeError) Code() string    { return err.code }
func (err *CodeError) Status() int     { return err.status }
func (err *CodeError) Message() string { return err.msg }

func (err *CodeError) Error() string {
	return fmt.Sprintf("%s(%d)", err.code, err.status)
}

var codeErrors = sync.Map{}

func loadOrStoreCodeErr(code string, status int, msg []string) (err *CodeError) {
	var s string
	if len(msg) > 0 {
		s = msg[len(msg)-1]
	}
	obj, loaded := codeErrors.LoadOrStore(code, &CodeError{
		code:   code,
		msg:    s,
		status: status,
	})
	err = obj.(*CodeError) // TODO 强制转换, 可能出现bug, 持续跟踪
	// if err == nil {
	// 	err = &CodeError{
	// 		Code:    code,
	// 		Message: msg,
	// 		Status:  status,
	// 	}
	// }
	if !loaded && (err.status != status || err.msg != s) {
		panic(fmt.Sprintf("apm: 错误码冲突: code=%q, status=%q, msg=%q", code, status, msg))
	}
	return
}

// func GetCodeErr(code string) {

// }

func NewError(code string, status int, desc ...string) (error, func(...interface{}) error) {
	err := loadOrStoreCodeErr(code, status, desc)
	return err, func(a ...interface{}) error {
		if len(a) == 0 {
			return err
		}
		return fmt.Errorf("%w: %s", err, fmt.Sprint(a...))
	}
}
func NewErrorf(code string, status int, desc ...string) (error, func(string, ...interface{}) error) {
	err := loadOrStoreCodeErr(code, status, desc)
	return err, func(f string, a ...interface{}) error {
		switch {
		case f != "" && len(a) != 0:
			f = fmt.Sprintf(f, a...)
		case f != "":
		case len(a) != 0:
			f = fmt.Sprint(a...)
		default:
			return err
		}

		return fmt.Errorf("%w: %s", err, f)
	}
}

type fieldsError struct {
	Err    error
	Fields []field.Field
}

func (err *fieldsError) Unwrap() error {
	return err.Err
}

func (err *fieldsError) Error() string {
	return err.Err.Error()
}

func (err *fieldsError) GetFields() []field.Field {
	return err.Fields
}

func GetFieldsFromError(err error) []field.Field {
	var ex *fieldsError
	if !errors.As(err, &ex) || ex == nil {
		return nil
	}
	return ex.Fields
}

func WrapFields(err error, fs ...field.Field) error {
	if err == nil {
		return nil
	}
	return &fieldsError{Err: err, Fields: fs}
}

func WrapCallStack(err error) error {
	if err == nil {
		return nil
	}
	var ex *StacktraceError
	if errors.Is(err, ex) {
		return err
	}

	return &StacktraceError{
		Err:   err,
		Stack: debug.Stack(),
	}
}

type StacktraceError struct {
	Stack []byte
	Err   error
}

func (err *StacktraceError) Unwrap() error {
	return err.Err
}

func (err *StacktraceError) Error() string {
	return err.Err.Error()
}

func StackToCallerInfo(stack []byte) []CallerInfo {
	lines := bytes.Split(stack, []byte{'\n'})
	// fmt.Printf("stack: %s\n", stack)

	// runtime: goroutine stack exceeds 1000000000-byte limit
	// runtime: sp=0xc020d50348 stack=[0xc020d50000, 0xc040d50000]
	// fatal error: stack overflow
	//
	// runtime stack:
	// ...
	// goroutine 218 [running]:
	// goroutine 4 [GC scavenge wait]:
	// goroutine 5 [finalizer wait]:
	// ...
	// goroutine 34 [GC worker (idle), 1 minutes]:
	start := false
	var dst []CallerInfo
	for i := 0; i < len(lines); i++ {
		l := lines[i]
		if bytes.Contains(l, []byte("goroutine ")) {
			i++
			continue
		}
		if !start {
			if bytes.Equal(l, []byte("runtime/debug.Stack()")) {
				start = true
				i++
			}
			continue
		}
		if bytes.Contains(l, []byte("apm.WrapCallStack({")) {
			i++
			continue
		}
		if bytes.Contains(l, []byte("testing.tRunner(")) {
			i++
			continue
		}
		if bytes.HasPrefix(l, []byte("created by ")) {
			break
		}

		i++
		if i >= len(lines) {
			fmt.Printf("例外的行: %v\n", i)
			fmt.Printf("%s\n", stack)
			break
		}
		arr := bytes.Split(bytes.Trim(lines[i], "\t"), []byte{' '})
		ci := CallerInfo{
			Method: string(bytes.SplitN(l, []byte{'('}, 2)[0]),
			File:   string(arr[0]),
		}
		{
			i := strings.LastIndex(ci.Method, "/")
			if i > 1 {
				ci.Package = ci.Method[:i]
				ci.Method = ci.Method[i+1:]
			}
		}
		{
			i := strings.LastIndex(ci.File, ":")
			if i > 1 {
				ci.Line = cast.ToInt(ci.File[i+1:])
				ci.File = ci.File[:i]
			}

			// i = strings.LastIndex(ci.File, "/")
			// if i > 1 {
			// 	ci.Path = ci.File[:i]
			// 	ci.File = ci.File[i+1:]
			// } else {
			// 	i = strings.LastIndex(ci.File, "\\")
			// 	if i > 1 {
			// 		ci.Path = ci.File[:i]
			// 		ci.File = ci.File[i+1:]
			// 	}
			// }
		}
		dst = append(dst, ci)
		// fmt.Printf("ci: %+v\n", ci)

	}
	return dst
}
func getSplitLast(s string, substr string) string {
	i := strings.LastIndex(s, substr)
	if i > 0 {
		s = s[i+1:]
	}
	return s
}
