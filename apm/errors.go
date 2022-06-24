package apm

import (
	"errors"
	"fmt"
	"sync"

	"github.com/junhwong/goost/pkg/field"
	"github.com/junhwong/goost/runtime"
)

// 标准错误
type CodeError struct {
	code   string // 固定的全局系统唯一的错误码, 如: NOTFOUND
	msg    string // 固定安全的错误消息, 用于对外部系统暂时. 如: 密码或账号不正确
	status int    // 用于如 HTTP/GRPC 等接口返回的状态码, -1 表示未定义
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

// Deprecated: use WrapCallStack
func WrapCallLast(err error, forceWrap ...bool) error {
	return runtime.WrapCallLast(err, 1, forceWrap...)
}

//
func WrapCallStack(err error) error {
	return runtime.WrapCallStacktrace(err, 1)
}
