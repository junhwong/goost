package security

import "fmt"

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
