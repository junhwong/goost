package web

import (
	"github.com/junhwong/goost/errors"
	"github.com/junhwong/goost/security"

	"github.com/gin-gonic/gin"
)

type MappingOption interface {
	apply(ops *MappingOptions)
}
type MappingOptions struct {
	TraceName     string
	NoTrace       bool
	BeforHandlers []func(Context)
	AfterHandlers []func(Context)
	do            func(*MappingOptions)
}

func (m *MappingOptions) apply(ops *MappingOptions) {
	m.do(ops)
}

type Routes interface {
	GET(path string, handler func(Context), options ...MappingOption) Routes
	POST(path string, handler func(Context), options ...MappingOption) Routes
	Handle(path, method string, handler func(Context), options ...MappingOption) Routes
}
type GroupRoutes interface {
	Routes
	Use(handlers ...func(Context)) Routes
}
type ControllerRoutes interface {
	GroupRoutes
	UseAuthenticate() GroupRoutes
}
type routesImpl struct {
	server   *Server
	basePath string
	routes   gin.IRoutes
}

func (r *routesImpl) Handle(path, method string, handler func(Context), options ...MappingOption) Routes {
	ops := &MappingOptions{
		TraceName:     r.basePath + path,
		BeforHandlers: make([]func(Context), 0),
		AfterHandlers: make([]func(Context), 0),
	}
	for _, option := range options {
		option.apply(ops)
	}
	handlers := []gin.HandlerFunc{}

	handlers = append(handlers, func(ctx Context) {
		ctx.Set(_traceNameInContextKey, ops.TraceName)
	})

	for _, it := range ops.BeforHandlers {
		handlers = append(handlers, gin.HandlerFunc(it))
	}

	handlers = append(handlers, gin.HandlerFunc(handler))

	for _, it := range ops.AfterHandlers {
		handlers = append(handlers, gin.HandlerFunc(it))
	}
	r.routes.Handle(method, path, handlers...)
	return r
}
func (r *routesImpl) GET(path string, handler func(Context), options ...MappingOption) Routes {
	return r.Handle(path, "GET", handler, options...)
}
func (r *routesImpl) POST(path string, handler func(Context), options ...MappingOption) Routes {
	return r.Handle(path, "POST", handler, options...)
}

func (r *routesImpl) Use(handlers ...func(Context)) Routes {
	m := []gin.HandlerFunc{}
	for _, it := range handlers {
		m = append(m, gin.HandlerFunc(it))
	}
	r.routes = r.routes.Use(m...)
	return r
}
func (r *routesImpl) UseAuthenticate() GroupRoutes {
	r.routes = r.routes.Use(func(ctx Context) {
		auth := security.AuthenticationFromContext(ctx)
		if !auth.IsAuthenticated() {
			panic(&errors.UnauthorizedError{})
		}
	})
	return r
}
