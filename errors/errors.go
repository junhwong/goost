package errors

import (
	"errors"
	"fmt"

	"github.com/junhwong/goost/security"
)

type Error struct {
	Message string // 错误信息
	Raise   error  // 由什么错误引起的
	File    string // 源码文件名
	Line    int    // 源码文件行号
	Stack   []byte // 原始栈
	Code    string // 错误码, 用于识别错误并处理
}

func (err *Error) Error() string {
	switch {
	case err.Message != "" && err.Code != "":
		return err.Code + ":" + err.Message
	case err.Message != "":
		return err.Message
	case err.Raise != nil:
		return err.Raise.Error()
	case err.Code != "":
		return err.Code
	}
	return ""
}

// MarshalJSON impl json.Marshaller to starand-logs error fields
func (err *Error) MarshalJSON() ([]byte, error) {
	return nil, nil
}

type AccessDeniedError struct {
	Err    error
	Any    bool
	Denied []security.Permission
}

func (err *AccessDeniedError) Error() string {
	temp := "Access Denied: mismatch all of %v"
	if err.Any {
		temp = "Access Denied: mismatch any of %v"
	}
	return fmt.Sprintf(temp, err.Denied)
}

func (err *AccessDeniedError) Unwrap() error {
	return err.Err
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

type UnauthorizedError struct {
	Err                        error
	WwwAuthenticateHeaderValue string // https://tools.ietf.org/html/rfc2617#section-3.2.1
}

func (err *UnauthorizedError) Unwrap() error {
	return err.Err
}
func (err *UnauthorizedError) Error() string {
	return "Unauthorized"
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

// type UnauthorizedError struct {
// 	description string
// }

// func NewUnauthorizedError(description string) error {
// 	return &UnauthorizedError{description: description}
// }

// func (err *UnauthorizedError) Error() string {
// 	return "Unauthorized error: " + err.description
// }

//
func Is(err error, target error) bool {
	return errors.Is(err, target)
}
func As(err error, target interface{}) bool {
	return errors.As(err, target)
}
func AsArgumentError(err error) (target *ArgumentError) {
	if !errors.As(err, &target) {
		return nil
	}
	return
}

type ErrorResponseEntry struct {
	Code    string `json:"code,omitempty"`
	Message string `json:"message,omitempty"`
	Field   string `json:"field,omitempty"`
}
type ErrorResponseBody struct {
	ErrorResponseEntry `json:",inline"`
	Errors             []ErrorResponseEntry `json:"errors,omitempty"`
}
