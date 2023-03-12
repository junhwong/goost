package web

import (
	"fmt"
	"net/http"
	"path"
	"reflect"
	"strings"

	"log"

	"github.com/gin-gonic/gin"
	"github.com/junhwong/goost/apm"
	"github.com/junhwong/goost/errors"
	"github.com/junhwong/goost/runtime"
	"github.com/junhwong/goost/security"
)

type Action = func(Context) Result
type ActionOption interface {
	applyOption(*routeOptions) error
}
type BindControllerOption interface {
	ActionOption
	applyBindControllerOption(*routeOptions) error
}
type bindControllerOptionFunc func(*routeOptions) error

func (f bindControllerOptionFunc) applyBindControllerOption(opt *routeOptions) error { return f(opt) }
func (f bindControllerOptionFunc) applyOption(opt *routeOptions) error               { return f(opt) }

func WithAuthed() bindControllerOptionFunc {
	return func(ro *routeOptions) error {
		return nil
	}
}

type Mapper interface {
	Controller(relativePath string, controller interface{}, options ...ActionOption) error
	//Action(method, relativePath string, action Action, options ...ActionOption) error
	//GET(relativePath string, action Action, options ...ActionOption) error
	//POST(relativePath string, action Action, options ...ActionOption) error
	//Dir(relativePath string, dir Dir, options ...ActionOption) error
	Static(relativePath, root string) error
}

type Router interface {
	http.Handler
	Mapper

	// Authenticated 请求要求身份验证中间件
	Authenticated() ActionOption
}

type ginRouter struct {
	*gin.Engine
	errorHooks  []func(ctx Context, err error) error
	authFilters []AuthenticationFilter
}

func (gr *ginRouter) route(method, relativePath string, action Action, options *routeOptions) error {
	ch := []gin.HandlerFunc{}

	for _, it := range options.BeforHandlers {
		ch = append(ch, gin.HandlerFunc(it))
	}

	ch = append(ch, func(c *gin.Context) {
		result := action(c)
		if result == nil {
			return
		}
		gr.handleError(c, result.Render(c))
	})

	for _, it := range options.AfterHandlers {
		ch = append(ch, gin.HandlerFunc(it))
	}

	gr.Handle(method, relativePath, ch...)

	return nil
}
func (r *ginRouter) Controller(relativePath string, controller interface{}, opts ...ActionOption) error {

	actions, err := getActions(controller)
	if err != nil {
		return err
	}
	options := &routeOptions{}
	for _, o := range opts {
		if err := o.applyOption(options); err != nil {
			return err
		}
	}
	for _, it := range actions {
		err = r.route(it.Method, path.Join(relativePath, it.Path), it.Action, options)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *ginRouter) Static(relativePath, root string) error {
	r.Engine.Static(relativePath, root)
	return nil
}

// 入口
func (r *ginRouter) enterHandler(ctx Context) {
	path := ctx.Request.URL.Path
	method := ctx.Request.Method

	// TODO 服务器版本、客户端版本、容器id、主机名称、IP、等等
	// TODO request-id  https://www.kancloud.cn/linimbus/envoyproxy/498765
	//c.Writer.Size()
	// "os.arch":"x86_64","os.platform":"linux","os.release":"ubuntu","os.version":"8.0"

	panic("todo")
	_, span := apm.Start(ctx,
		apm.WithName(method+" "+path),
		// apm.WithTrimFieldPrefix("__web."),
		// apm.WithFields(
		// 	clientIP(ctx.ClientIP()),
		// 	httpMethod(method),
		// ),
	)
	defer func() {
		span.End(
			// 用于替换自定义名称
			apm.WithReplaceSpanName(func() (s string) {
				s, _ = ctx.Value(_traceNameInContextKey).(string)
				return
			}),
			// apm.WithFields(
			// 	httpStatus(ctx.Writer.Status()),
			// ),
		)
	}()

	defer func() {
		if re := recover(); re != nil {
			var err = &errors.Exception{} //TODO stack
			if ex, _ := re.(error); ex != nil {
				err.Err = ex
			} else {
				err.Raise = re
			}
			span.Fail(err)
			r.handleError(ctx, err)
		}
	}()

	//ctx.Error()

	ctx.Next()

	if len(ctx.Errors) == 0 {
		return
	}
	for _, err := range ctx.Errors {
		// TODO LOG
		if err == nil || ctx.IsAborted() {
			continue
		}
		r.handleError(ctx, err)
	}
}
func (r *ginRouter) doAuth(ctx Context) (security.Authentication, error) {
	for _, filter := range r.authFilters {
		if ctx.IsAborted() {
			continue
		}
		auth, err := filter.Filter(ctx)
		if err != nil {
			return nil, err
		}
		if auth != nil && auth.IsAuthenticated() {
			return auth, nil
		}
	}
	return nil, &security.UnauthorizedError{Cause: errors.Cause{Err: fmt.Errorf("not impl")}}
}
func (r *ginRouter) handleError(c Context, err error) {
	if err == nil {
		return
	}
	for _, hook := range r.errorHooks {
		err = hook(c, err)
		if err == nil {
			return // 只成功处理一次
		}
	}
	// TODO: log
	log.Printf("web.handleError: %+v\n", err)
	if !c.IsAborted() {
		if res, ok := err.(errors.ErrorResponse); ok {
			c.JSON(res.ResponseStatusCode(), res.ResponseData())
			c.Abort()
			return
		}
		if res := new(security.UnauthorizedError); errors.As(err, &res) {
			if res.WwwAuthenticateHeaderValue != "" {
				c.Header("WWW-Authenticate", res.WwwAuthenticateHeaderValue)
			}

			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
	}

	c.Abort() // gin middleware 实现有bug. https://github.com/gin-gonic/gin/issues/2221
	if err != nil {
		// TODO gin.Recovery()
		panic(err)
	}
}

type actionMap struct {
	Method string
	Path   string
	Action Action
}

func getActions(controller interface{}) ([]actionMap, error) {
	refVal := reflect.ValueOf(controller)
	if refVal.Kind() != reflect.Ptr {
		return nil, fmt.Errorf("无效")
	}
	ct := refVal.Type()
	actions := []actionMap{}
	for i := 0; i < ct.NumMethod(); i++ {
		method := ct.Method(i)
		var action Action
		if a, ok := refVal.Method(i).Interface().(Action); ok {
			action = a
		}
		if action == nil {
			continue
		}
		switch {
		case strings.HasPrefix(method.Name, "Get"):
			actions = append(actions, actionMap{
				Method: "GET",
				Path:   strings.ToLower(method.Name[3:]),
				Action: action,
			})
		case strings.HasPrefix(method.Name, "Post"):
			actions = append(actions, actionMap{
				Method: "POST",
				Path:   strings.ToLower(method.Name[4:]),
				Action: action,
			})
		default:
			continue
		}

	}
	return actions, nil
}

type RouterOption interface {
	apply(*ginRouter) error
}
type routerOptionFunc func(*ginRouter) error

func (f routerOptionFunc) apply(opts *ginRouter) error { return f(opts) }

type RouterFactoryInit struct {
	runtime.In

	AuthenticationFilters []AuthenticationFilter `group:"authenticationFilters"`
}

// RouterFactory 构造路由
func RouterFactory(opts ...RouterOption) func(RouterFactoryInit) (Router, error) {
	return func(in RouterFactoryInit) (Router, error) {
		engine := gin.Default()
		r := &ginRouter{
			Engine: engine,
		}
		engine.NoRoute(func(c *gin.Context) { c.AbortWithStatus(404) })
		engine.NoMethod(func(c *gin.Context) { c.AbortWithStatus(405) })
		engine.Use(r.enterHandler)

		for _, it := range in.AuthenticationFilters {
			if it == nil {
				continue
			}
			r.authFilters = append(r.authFilters, it)
		}
		for _, opt := range opts {
			if opt == nil {
				continue
			}
			if err := opt.apply(r); err != nil {
				return nil, err
			}
		}
		return r, nil
	}
}

func WithErrorHook(f func(ctx Context, err error) error) routerOptionFunc {
	return func(gr *ginRouter) error {
		if f == nil {
			return fmt.Errorf("f cannot be nil")
		}
		gr.errorHooks = append(gr.errorHooks, f)
		return nil
	}
}
func WithAuthenticationFilter(v ...AuthenticationFilter) routerOptionFunc {
	return func(gr *ginRouter) error {
		for _, it := range v {
			if it == nil {
				continue
			}
			gr.authFilters = append(gr.authFilters, it)
		}
		return nil
	}
}
