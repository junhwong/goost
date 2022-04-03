package errors

import (
	"errors"
	"fmt"
	"strings"
)

// Exception 表示一个异常, 通常由`panic`引发并捕获。
type Exception struct {
	Err   error
	Raise interface{} // recover 捕获的值
	File  string      // 源码文件名
	Line  int         // 源码文件行号
	Stack []byte      // 原始栈
}

func (ex *Exception) Unwrap() error {
	return ex.Err
}
func (ex *Exception) Error() string {
	return fmt.Sprintf("%#v", ex)
}

// Error 通用错误
type Error struct {
	Message string // 错误信息。通常该值是可公开的，不包含敏感信息。
	Err     error  // 由什么错误引起的
	// File    string // 源码文件名
	// Line    int    // 源码文件行号
	// Stack   []byte // 原始栈
	Code   string // 错误码, 用于识别错误并处理
	Status int    // 状态码
	Field  string // 字段名称
}

func (err *Error) Unwrap() error {
	return err.Err
}
func (err *Error) Error() string {
	switch {
	case err.Message != "" && err.Code != "":
		return err.Code + ":" + err.Message
	case err.Message != "":
		return err.Message
	case err.Err != nil:
		return err.Err.Error()
	case err.Code != "":
		return err.Code
	}

	return "Error"
}

type Cause struct {
	Err error
}

func (err *Cause) Unwrap() error {
	return err.Err
}

// MarshalJSON impl json.Marshaller to starand-logs error fields
func (err *Error) MarshalJSON() ([]byte, error) {
	return nil, nil
}

type InvalidParameterError struct {
	Err     error
	Field   string
	Message string
}

func (err *InvalidParameterError) Unwrap() error {
	return err.Err
}

func (err *InvalidParameterError) Error() string {

	return fmt.Sprintf("Invalid parameter: %q, %s", err.Field, err.Message)
}

type IllegalArgumentError struct {
	Err     error
	Code    string // missing missing_field invalid already_exists unprocessable custom
	Field   string
	Message string
}

func (err *IllegalArgumentError) Unwrap() error {
	return err.Err
}
func (err *IllegalArgumentError) Error() string {
	return fmt.Sprintf("Illegal argument: %v", err.Err)
}

type ArgumentError struct {
	Err     error
	Code    string // missing missing_field invalid already_exists unprocessable custom
	Field   string
	Message string
}

func (err *ArgumentError) Unwrap() error {
	return err.Err
}
func (err *ArgumentError) Error() string {
	s := "ArgumentError: "
	if err.Code != "" {
		s += err.Code
		if err.Message != "" {
			s += ": " + err.Message
		}
	} else if err.Message != "" {
		s += err.Message
	}

	if err.Err != nil {
		s += ", " + err.Err.Error()
	}
	if err.Field != "" {
		s += ". field: " + err.Field
	}
	return s
}

func NewInvalidArgumentError(field, message string, err error) *ArgumentError {
	return &ArgumentError{Err: err,
		Code:    "invalid_argument",
		Field:   field,
		Message: message,
	}
}
func WrapInvalidArgumentError(err error, field string) *ArgumentError {
	return &ArgumentError{Err: err,
		Code:  "invalid_argument",
		Field: field,
	}
}

// type UnauthorizedError struct {
// 	description string
// }

// func NewUnauthorizedError(description string) error {
// 	return &UnauthorizedError{description: description}
// }

// func (err *UnauthorizedError) Error() string {
// 	return "Unauthorized error: " + err.description
// }

var (
	Is = errors.Is
	As = errors.As
)

// func Is(err error, target error) bool {
// 	return errors.Is(err, target)
// }
// func As(err error, target interface{}) bool {
// 	return errors.As(err, target)
// }
func AsArgumentError(err error) (target *ArgumentError) {
	if !errors.As(err, &target) {
		return nil
	}
	return
}

// ErrorCodeGetter 错误码
type ErrorCodeGetter interface {
	ErrorCode() string
}

// ErrorCodeGetter 错误码
type ErrorMessageGetter interface {
	ErrorMessage() string
}

// StatusCodeGetter 状态码
type StatusCodeGetter interface {
	StatusCode() int
}

////

type CodeError struct {
	Err  error
	Code string
	// Message string
}

type codeError string

func (err codeError) Error() string {
	return string(err)
}
func (err codeError) ErrorCode() string {
	return strings.SplitN(string(err), " ", 2)[0]
}
func (err codeError) ErrorMessage() string {
	if arr := strings.SplitN(string(err), " ", 2); len(arr) == 2 {
		return arr[1]
	}
	return ""
}

func New(code string, message ...string) error {
	code = strings.TrimSpace(code)
	if len(code) == 0 {
		panic("errors: code cannot empty")
	}
	s := strings.Join(message, "")
	if s == "" {
		return codeError(code)
	}
	return codeError(code + " " + s)
}

func GetErrorCode(err error) string {
	if err == nil {
		return ""
	}
	if f, ok := err.(ErrorCodeGetter); ok {
		if s := f.ErrorCode(); s != "" {
			return s
		}
	}
	return GetErrorCode(errors.Unwrap(err))
}
func GetErrorMessage(err error) string {
	if err == nil {
		return ""
	}
	if f, ok := err.(ErrorMessageGetter); ok {
		if s := f.ErrorMessage(); s != "" {
			return s
		}
	}

	return GetErrorCode(errors.Unwrap(err))
}

////
func WithErrorCode(err error, code string) error {
	if len(code) == 0 {
		panic("errors: code cannot empty")
	}
	base, _ := err.(*Error)
	if base != nil && base.Code == "" {
		base.Code = code
		return base
	}

	e := &Error{
		Err:  err,
		Code: code,
	}
	return e
}
func WithStatusCode(err error, v int) error {
	if v == 0 {
		panic("errors: v cannot be zero")
	}
	base, _ := err.(*Error)
	if base != nil && base.Code == "" {
		base.Status = v
		return base
	}

	e := &Error{
		Err:    err,
		Status: v,
	}
	return e
}
