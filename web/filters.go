package web

import "github.com/junhwong/goost/security"

type AuthenticationFilter interface {
	// 处理认证
	Filter(Context) (security.Authentication, error)
}
