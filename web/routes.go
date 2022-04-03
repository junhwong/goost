package web

// type RouteOption interface {
// 	apply(ops *routeOptions)
// }

// type Routes interface {
// 	GET(path string, handler Handler, options ...RouteOption) Routes
// 	POST(path string, handler Handler, options ...RouteOption) Routes
// 	Static(path, dir string) Routes
// 	Handle(path, method string, handler Handler, options ...RouteOption) Routes
// }

// type GroupRoutes interface {
// 	Routes
// 	Use(handlers ...Midleware) Routes
// }

// type ControllerRoutes interface {
// 	GroupRoutes
// 	UseAuthenticate() GroupRoutes
// }

type routeOptions struct {
	TraceName     string
	NoTrace       bool
	BeforHandlers []Midleware
	AfterHandlers []Midleware
	do            func(*routeOptions)
}

func (m *routeOptions) apply(ops *routeOptions) {
	m.do(ops)
}

// type routesImpl struct {
// 	server       *Server
// 	basePath     string
// 	routes       gin.IRoutes
// 	errorHandler ErrorHandler
// }

// func (r *routesImpl) Handle(path, method string, handler Handler, options ...RouteOption) Routes {
// 	ops := &routeOptions{
// 		TraceName:     r.basePath + path,
// 		BeforHandlers: make([]Midleware, 0),
// 		AfterHandlers: make([]Midleware, 0),
// 	}
// 	for _, option := range options {
// 		option.apply(ops)
// 	}
// 	handlers := []gin.HandlerFunc{}

// 	handlers = append(handlers, func(ctx Context) {
// 		ctx.Set(_traceNameInContextKey, ops.TraceName)
// 	})

// 	for _, it := range ops.BeforHandlers {
// 		handlers = append(handlers, gin.HandlerFunc(it))
// 	}

// 	handlers = append(handlers, func(ctx Context) {
// 		result := handler(ctx)
// 		if result == nil {
// 			return
// 		}
// 		// TODO: 传入server对象供自定义处理
// 		if err := result.Render(ctx); err != nil {
// 			if r.errorHandler != nil {
// 				err = r.errorHandler.HandleError(ctx, err)
// 			}
// 			if err != nil {
// 				r.server.handleError(ctx, err)
// 			}
// 		}
// 	})

// 	for _, it := range ops.AfterHandlers {
// 		handlers = append(handlers, gin.HandlerFunc(it))
// 	}
// 	r.routes.Handle(method, path, handlers...)
// 	return r
// }
// func (r *routesImpl) GET(path string, handler Handler, options ...RouteOption) Routes {
// 	return r.Handle(path, "GET", handler, options...)
// }
// func (r *routesImpl) POST(path string, handler Handler, options ...RouteOption) Routes {
// 	return r.Handle(path, "POST", handler, options...)
// }

// func (r *routesImpl) Static(path, dir string) Routes {
// 	r.routes = r.routes.Static(path, dir)
// 	return r
// }

// func (r *routesImpl) Use(handlers ...Midleware) Routes {
// 	m := []gin.HandlerFunc{}
// 	for _, it := range handlers {
// 		m = append(m, gin.HandlerFunc(it))
// 	}

// 	r.routes = r.routes.Use(m...)
// 	return r
// }

// func (r *routesImpl) UseAuthenticate() GroupRoutes {
// 	r.routes = r.routes.Use(func(ctx Context) {
// 		auth := security.AuthenticationFromContext(ctx)
// 		if !auth.IsAuthenticated() {
// 			panic("&errors.UnauthorizedError{}")
// 		}
// 	})
// 	return r
// }
