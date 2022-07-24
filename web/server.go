package web

import (
	"context"
	"fmt"
	"net/http"

	"github.com/junhwong/goost/errors"
	"github.com/junhwong/goost/pkg/field"
	"github.com/junhwong/goost/runtime"
)

// // Server
// type Server struct {
// 	controllers map[string][]Controller
// 	entryPoints []AuthenticationEntryPoint
// 	errorHooks  []func(ctx Context, err error) error
// 	mu          sync.Mutex
// 	settings    serverSettings
// }

// // NewServer 返回一个新的 Server
// func NewServer(options ...ServerOption) *Server {
// 	settings := serverSettings{
// 		Addr:                    ":8086",
// 		GracefulShutdownTimeout: time.Second * 60,
// 	}
// 	for _, it := range options {
// 		if it != nil {
// 			it.apply(&settings)
// 		}
// 	}
// 	return &Server{
// 		settings:    settings,
// 		controllers: make(map[string][]Controller),
// 		entryPoints: []AuthenticationEntryPoint{},
// 		errorHooks:  make([]func(ctx Context, err error) error, 0),
// 	}
// }

// func (srv *Server) Run(stopCh <-chan struct{}) error {
// 	srv.mu.Lock()
// 	defer srv.mu.Unlock()

// 	engine := gin.Default()
// 	// gin.Mode()
// 	engine.Use(srv.enterHandler)
// 	engine.Use(srv.authHandler)
// 	engine.Use(cors())

// 	for relativePath, controllers := range srv.controllers {
// 		g := engine.Group(relativePath)
// 		for _, ctl := range controllers {
// 			impl := &routesImpl{
// 				server:   srv,
// 				routes:   g,
// 				basePath: relativePath,
// 			}
// 			if h, ok := ctl.(ErrorHandler); ok {
// 				impl.errorHandler = h
// 			}
// 			if err := ctl.Init(srv, impl); err != nil {
// 				return err
// 			}
// 		}
// 	}
// 	// fmt.Println("===", srv.settings.Addr)
// 	lis := &http.Server{
// 		Addr:    srv.settings.Addr,
// 		Handler: engine,
// 		//ErrorLog
// 	}

// 	go func() {
// 		if err := lis.ListenAndServe(); err != nil && err != http.ErrServerClosed {
// 			log.Fatalf("listen: %s\n", err)
// 		}
// 	}()
// 	<-stopCh

// 	// Graceful shutdown wait timeout
// 	ctx, cancel := context.WithTimeout(context.Background(), srv.settings.GracefulShutdownTimeout)
// 	defer cancel()
// 	if err := lis.Shutdown(ctx); err != nil {
// 		log.Fatal("Server forced to shutdown:", err)
// 	}
// 	log.Println("done")
// 	return nil
// }

// func (srv *Server) OnError(hook func(ctx Context, err error) error) {
// 	srv.errorHooks = append(srv.errorHooks, hook)
// }

// func (srv *Server) handleError(c Context, err error) {
// 	if err == nil {
// 		return
// 	}
// 	for _, hook := range srv.errorHooks {
// 		if err == nil {
// 			break
// 		}
// 		err = hook(c, err)
// 	}
// 	if err != nil {
// 		// TODO gin.Recovery()
// 		panic(err)
// 	}

// 	c.Abort() // gin middleware 实现有bug. https://github.com/gin-gonic/gin/issues/2221
// }

const (
	_traceNameInContextKey = "$web.TraceName"
)

// func (srv *Server) authHandler(ctx Context) {
// 	// TODO 是否每次都需要调用
// 	for _, entry := range srv.entryPoints {
// 		if ctx.IsAborted() {
// 			return
// 		}
// 		auth, err := entry.Commence(ctx)
// 		if err != nil {
// 			srv.handleError(ctx, err)
// 			return
// 		}
// 		if auth != nil && auth.IsAuthenticated() {
// 			ctx.Set(security.AuthenticationContextKey, auth)
// 			break
// 		}
// 	}
// 	if !ctx.IsAborted() {
// 		ctx.Next()
// 	}
// }

var (
	ClientIPKey, clientIP             = field.String("client.ip")
	HTTPResponseStatusKey, httpStatus = field.Int("http.response.status_code")
	HTTPRequestMethodKey, httpMethod  = field.String("http.request.method")
)

// // 入口
// func (srv *Server) enterHandler(ctx Context) {
// 	path := ctx.Request.URL.Path
// 	method := ctx.Request.Method

// 	//TODO 服务器版本、客户端版本、容器id、主机名称、IP、等等

// 	//c.Writer.Size()
// 	// "os.arch":"x86_64","os.platform":"linux","os.release":"ubuntu","os.version":"8.0"
// 	span := apm.Start(ctx,
// 		apm.WithName(method+" "+path),
// 		apm.WithTrimFieldPrefix("__web."),
// 		apm.WithFields(
// 			clientIP(ctx.ClientIP()),
// 			httpMethod(method),
// 		),
// 	)
// 	defer func() {
// 		span.End(
// 			// 用于替换自定义名称
// 			apm.WithReplaceSpanName(func() (s string) {
// 				s, _ = ctx.Value(_traceNameInContextKey).(string)
// 				return
// 			}),
// 			apm.WithFields(
// 				httpStatus(ctx.Writer.Status()),
// 			),
// 		)
// 	}()

// 	defer func() {
// 		if r := recover(); r != nil {
// 			var err error
// 			if e, _ := r.(error); e != nil {
// 				err = e
// 			} else {
// 				// TODO 封装错误
// 				err = fmt.Errorf("recoverd error: %v", r)
// 			}
// 			span.Fail()
// 			srv.handleError(ctx, err)
// 		}
// 	}()

// 	ctx.Next()
// }

// // 注册认证入口
// func (srv *Server) EntryPoint(v ...AuthenticationEntryPoint) {
// 	srv.mu.Lock()
// 	defer srv.mu.Unlock()

// 	for _, p := range v {
// 		if p == nil {
// 			continue
// 		}
// 		srv.entryPoints = append(srv.entryPoints, p)
// 	}
// }

// func (srv *Server) RegisterAuthenticationFilters(filters ...AuthenticationFilter) {

// }

// // Deprecated: 请使用 RegisterController 替换
// func (srv *Server) Route(relativePath string, controllers ...Controller) {
// 	srv.mu.Lock()
// 	defer srv.mu.Unlock()

// 	if relativePath == "" {
// 		relativePath = "/"
// 	}
// 	for _, ctl := range controllers {
// 		if ctl == nil {
// 			continue
// 		}
// 		srv.controllers[relativePath] = append(srv.controllers[relativePath], ctl)
// 	}
// }

// // 注册Controller
// func (srv *Server) RegisterController(relativePath string, controllers ...Controller) {
// 	srv.mu.Lock()
// 	defer srv.mu.Unlock()

// 	if relativePath == "" {
// 		relativePath = "/"
// 	}
// 	for _, ctl := range controllers {
// 		if ctl == nil {
// 			continue
// 		}
// 		srv.controllers[relativePath] = append(srv.controllers[relativePath], ctl)
// 	}
// }

// type ServerContext interface {
// 	// security.AuthenticationManager
// }

// func cors() gin.HandlerFunc {
// 	return func(c *gin.Context) {

// 		origin := "*" // c.Request.Header.Get("Origin")
// 		// if origin == "" {
// 		// 	origin = "*"
// 		// }
// 		exposes := []string{
// 			"Set-Access-Token",
// 			"X-CSRF-Token",
// 			"X-State",
// 		}
// 		allowHeaders := []string{
// 			// "Access-Control-Request-Headers",
// 			// "Access-Control-Request-Method",
// 			"Authorization",
// 			"Content-Type",
// 			"Content-Length",
// 			"Cookie",
// 			"Accept",
// 			"Accept-Encoding",
// 			"Accept-Language",
// 			"User-Agent",
// 			"Host",
// 			"Referer",
// 			"Cache-Control",
// 			"Connection",
// 			"DNT",
// 			"Origin",
// 			"Pragma",
// 			"TE",
// 			"Host",

// 			"X-Requested-With",
// 			"X-State",
// 			"X-App-Version",
// 		}
// 		c.Header("Access-Control-Allow-Origin", origin)
// 		c.Header("Access-Control-Allow-Credentials", "true")
// 		c.Header("Access-Control-Expose-Headers", strings.Join(exposes, ","))
// 		c.Header("Access-Control-Allow-Headers", strings.Join(allowHeaders, ","))
// 		c.Header("Access-Control-Allow-Methods", "POST,GET,PUT,DELETE,CONNECT,TRACE,PATCH,HEAD")
// 		c.Header("Access-Control-Max-Age", "3600")

// 		if strings.ToUpper(c.Request.Method) == "OPTIONS" {
// 			c.AbortWithStatus(http.StatusOK)
// 			return
// 		}
// 		c.Next()
// 	}
// }

type ServerOption interface {
	applyServerOption(*http.Server) error
}

// ServerFactory 返回 http.Server 构造函数
func ServerFactory(opts ...ServerOption) func(runtime.Lifecycle, Router) error {
	return func(lifecycle runtime.Lifecycle, router Router) error {
		lis := &http.Server{
			Addr:    ":8086",
			Handler: router,
			//ErrorLog
		}

		stoped := false
		lifecycle.Append(func(ctx context.Context) {
			fmt.Println("ListenAndServe:", lis.Addr)
			go func() {
				<-ctx.Done()
				lis.Shutdown(context.TODO())
			}()
			err := lis.ListenAndServe()
			if err == nil || (stoped && errors.Is(err, http.ErrServerClosed)) {
				return
			}
			if err != nil {
				fmt.Println("ListenAndServe Faile:", err)
			}

		})
		// lifecycle.Append(runtime.Hook{
		// 	OnStart: func(c context.Context) {
		// 		fmt.Println("ListenAndServe:", lis.Addr)
		// 		err := lis.ListenAndServe()
		// 		if err == nil || (stoped && errors.Is(err, http.ErrServerClosed)) {
		// 			return
		// 		}
		// 		if err != nil {
		// 			fmt.Println("ListenAndServe Faile:", err)
		// 		}

		// 	},
		// 	OnStop: func(c context.Context) {
		// 		if stoped {
		// 			return
		// 		}
		// 		stoped = true
		// 		err := lis.Shutdown(c)
		// 		if err == nil || errors.Is(err, http.ErrServerClosed) {
		// 			return
		// 		}
		// 		if err != nil {
		// 			fmt.Println("Shutdown Faile:", err)
		// 		}
		// 	},
		// })
		return nil
	}
}
