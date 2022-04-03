package web

import "github.com/gin-gonic/gin"

// Context 请求上下文。
// 注意：禁止使用 `*gin.Context`，因为以后可能改变实现。
type Context = *gin.Context

type Result interface {
	Render(ctx Context) error
}

type Midleware gin.HandlerFunc

type Handler func(Context) Result
