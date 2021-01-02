package web

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/junhwong/goost/apm"
	"github.com/junhwong/goost/security"

	"github.com/gin-gonic/gin"
)

type Server struct {
	controllers map[string][]Controller
	entryPoints []AuthenticationEntryPoint
	errHooks    []func(ctx Context, err error) error
	mu          sync.Mutex
}

func NewServer() *Server {
	return &Server{
		controllers: make(map[string][]Controller),
		entryPoints: []AuthenticationEntryPoint{},
		errHooks:    make([]func(ctx *gin.Context, err error) error, 0),
	}
}

func (srv *Server) Run(stopCh <-chan struct{}) error {
	srv.mu.Lock()
	defer srv.mu.Unlock()

	engine := gin.Default()
	engine.Use(srv.errorHandler)
	engine.Use(cors())
	points := make([]AuthenticationEntryPoint, len(srv.entryPoints))
	copy(points, srv.entryPoints)
	engine.Use(func(c *gin.Context) {
		for _, entry := range points {
			if c.IsAborted() {
				return
			}
			auth, err := entry.Commence(c)
			if err != nil {
				srv.handleError(c, err)
				return
			}
			if auth != nil && auth.IsAuthenticated() {
				c.Set(security.AuthenticationContextKey, auth)
				break
			}
		}
		if !c.IsAborted() {
			c.Next()
		}
	})

	for relativePath, controllers := range srv.controllers {
		g := engine.Group(relativePath)
		for _, ctl := range controllers {
			if err := ctl.Init(srv, &routesImpl{server: srv, routes: g, basePath: relativePath}); err != nil {
				return err
			}
		}
	}

	lis := &http.Server{
		Addr:    ":8080",
		Handler: engine,
	}
	go func() {
		if err := lis.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()
	<-stopCh
	//Graceful shutdown wait timeout
	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Second)
	defer cancel()
	if err := lis.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}
	log.Println("done")
	return nil
}
func (srv *Server) OnError(hook func(ctx Context, err error) error) {
	srv.errHooks = append(srv.errHooks, hook)
}
func (srv *Server) handleError(c *gin.Context, err error) {
	if err == nil {
		return
	}
	for _, hook := range srv.errHooks {
		if err == nil {
			break
		}
		err = hook(c, err)
	}
	if err != nil {
		// gin.Recovery()
		panic(err)
	}

	c.Abort() // gin middleware 实现有bug. https://github.com/gin-gonic/gin/issues/2221
}

const (
	_traceNameInContextKey = "$web.TraceName"
)

func (srv *Server) errorHandler(ctx *gin.Context) {
	//TODO 服务器版本、客户端版本、容器id、主机名称、IP、等等
	span := apm.Start(ctx, apm.WithName(ctx.Request.Method+" "+ctx.Request.URL.Path))
	defer span.Finish(apm.WithReplaceSpanName(func() (s string) {
		s, _ = ctx.Value(_traceNameInContextKey).(string)
		return
	}))

	defer func() {
		if r := recover(); r != nil {
			var err error
			if e, _ := r.(error); e != nil {
				err = e
			} else {
				// TODO 封装错误
				err = fmt.Errorf("recoverd error: %v", r)
			}
			span.Fail()
			srv.handleError(ctx, err)
		}
	}()

	ctx.Next()
}

// 注册认证入口
func (srv *Server) EntryPoint(v ...AuthenticationEntryPoint) {
	srv.mu.Lock()
	defer srv.mu.Unlock()

	for _, p := range v {
		if p == nil {
			continue
		}
		srv.entryPoints = append(srv.entryPoints, p)
	}
}

func (srv *Server) RegisterAuthenticationFilters(filters ...AuthenticationFilter) {

}

func (srv *Server) Route(relativePath string, controllers ...Controller) {
	srv.mu.Lock()
	defer srv.mu.Unlock()

	if relativePath == "" {
		relativePath = "/"
	}
	for _, ctl := range controllers {
		if ctl == nil {
			continue
		}
		srv.controllers[relativePath] = append(srv.controllers[relativePath], ctl)
	}
}

type ServerContext interface {
	// security.AuthenticationManager
}

func cors() gin.HandlerFunc {
	return func(c *gin.Context) {

		origin := "*" // c.Request.Header.Get("Origin")
		// if origin == "" {
		// 	origin = "*"
		// }
		exposes := []string{
			"Set-Access-Token",
			"X-CSRF-Token",
			"X-State",
		}
		allowHeaders := []string{
			// "Access-Control-Request-Headers",
			// "Access-Control-Request-Method",
			"Authorization",
			"Content-Type",
			"Content-Length",
			"Cookie",
			"Accept",
			"Accept-Encoding",
			"Accept-Language",
			"User-Agent",
			"Host",
			"Referer",
			"Cache-Control",
			"Connection",
			"DNT",
			"Origin",
			"Pragma",
			"TE",
			"Host",

			"X-Requested-With",
			"X-State",
			"X-App-Version",
		}
		c.Header("Access-Control-Allow-Origin", origin)
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Expose-Headers", strings.Join(exposes, ","))
		c.Header("Access-Control-Allow-Headers", strings.Join(allowHeaders, ","))
		c.Header("Access-Control-Allow-Methods", "POST,GET,PUT,DELETE,CONNECT,TRACE,PATCH,HEAD")
		c.Header("Access-Control-Max-Age", "3600")

		if strings.ToUpper(c.Request.Method) == "OPTIONS" {
			c.AbortWithStatus(http.StatusOK)
			return
		}
		c.Next()
	}
}
