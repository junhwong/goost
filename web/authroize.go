package web

import "github.com/junhwong/goost/security"

// AuthenticationEntryPoint 定义认证入口点
type AuthenticationEntryPoint interface {
	Commence(Context) (security.Authentication, error)
}
