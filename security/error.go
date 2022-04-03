package security

import (
	"fmt"
	"net/http"

	"github.com/junhwong/goost/errors"
)

// 无效的token错误。参考： https://tools.ietf.org/html/rfc2617#section-3.2.1
type InvaildTokenError struct {
	Err       error
	TokenType string
	Realm     string
	Scope     string
	Message   string
}

func (err *InvaildTokenError) Unwrap() error {
	return err.Err
}
func (err *InvaildTokenError) Error() string {
	return fmt.Sprintf("Invaild Token Error: %s", err.Err)
}
func (err *InvaildTokenError) ResponseStatusCode() int {
	return http.StatusUnauthorized
}
func (err *InvaildTokenError) ResponseData() interface{} {
	code := "invaild_token"
	if code == "" {
		code = "invaild_token"
	}
	details := []errors.ErrorDetail{}
	// for _, e := range err.Details {
	// 	if rd, ok := e.(errors.ErrorResponseDetail); ok {
	// 		details = append(details, rd.ResponseDetailData())
	// 	}
	// }
	if len(details) == 0 {
		details = nil
	}
	return errors.ErrorResult{
		Code:    code,
		Message: err.Message,
		Details: details,
	}
}

type UnauthorizedError struct {
	errors.Cause
	WwwAuthenticateHeaderValue string // https://tools.ietf.org/html/rfc2617#section-3.2.1
}

func (err *UnauthorizedError) Error() string {
	return "Unauthorized"
}

// func (err *UnauthorizedError) Result() string {
// 	return "Unauthorized"
// }

type AccessDeniedError struct {
	Err    error
	Any    bool
	Denied []Permission
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

// 范围不足
type InsufficientScopeError struct {
	Cause   error
	Scope   string
	Message string
}

func (err *InsufficientScopeError) Error() string {
	return "InsufficientScopeError"
}
func (err *InsufficientScopeError) Unwrap() error {
	return err.Cause
}

// 账户比匹配错误。
type AccountMismatchedError struct {
	Err error
}
