package web

import (
	"net/http"

	"github.com/junhwong/goost/security"
)

// AuthenticationEntryPoint 定义认证入口点
type AuthenticationEntryPoint interface {
	// 处理认证
	Commence(Context) (security.Authentication, error)
}

// Authority
func Authority(perms ...security.Permission) bindControllerOptionFunc {
	if len(perms) == 0 {
		panic("perms cannot be nil")
	}
	return func(ro *routeOptions) error {
		handler := func(ctx Context) {
			auth := security.AuthenticationFromContext(ctx)
			for _, p := range perms {
				if !auth.IsGranted(p) {
					_ = ctx.AbortWithError(http.StatusForbidden, &security.AccessDeniedError{Any: false, Denied: []security.Permission{p}})
					return
				}
			}
			ctx.Next()
		}
		ro.BeforHandlers = append(ro.BeforHandlers, handler)
		return nil
	}
}

// AuthorityAny
func AuthorityAny(perms ...security.Permission) bindControllerOptionFunc {
	if len(perms) == 0 {
		panic("perms cannot be nil")
	}
	return func(ro *routeOptions) error {
		handler := func(ctx Context) {
			auth := security.AuthenticationFromContext(ctx)
			for _, p := range perms {
				if auth.IsGranted(p) {
					ctx.Next()
					return
				}
			}
			// ctx.AbortWithStatusJSON(http.StatusForbidden, nil)
			_ = ctx.AbortWithError(http.StatusForbidden, &security.AccessDeniedError{Any: true, Denied: perms})
		}
		ro.BeforHandlers = append(ro.BeforHandlers, handler)
		return nil
	}
}

func (r *ginRouter) Authenticated() ActionOption {
	return bindControllerOptionFunc(func(opts *routeOptions) error {
		handler := func(ctx Context) {
			auth := security.AuthenticationFromContext(ctx)
			if !auth.IsAuthenticated() {
				var err error
				auth, err = r.doAuth(ctx)
				if err != nil {
					r.handleError(ctx, err)
					if ctx.IsAborted() {
						return
					}
					ctx.JSON(http.StatusUnauthorized, map[string]string{"code": "unauthorized"})
					ctx.Abort()
					return
				}
				ctx.Set(security.AuthenticationContextKey, auth)
			}

			ctx.Next()
		}
		opts.BeforHandlers = append(opts.BeforHandlers, handler)
		return nil
	})
}
