package runtime

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"reflect"
	"runtime"
	"sync"
	"syscall"
	"time"

	"go.uber.org/dig"
)

var (
	Debug = func(a ...any) { fmt.Println(a...) }
)

func Debugf(format string, a ...any) { Debug(fmt.Sprintf(format+"\n", a...)) }

type Hook struct {
	*hookCtx
	once    sync.Once
	serving bool

	hook        func(ctx context.Context)
	servingHook func(ctx context.Context, onStarted func())
}
type hookCtx struct {
	ctx    context.Context
	cancel context.CancelFunc
}

func (h *Hook) doRun(wg *sync.WaitGroup, next func(), stop func(s string), m *sync.Map) {
	h.once.Do(func() {
		defer wg.Done()
		defer h.cancel()
		if h.servingHook != nil {
			fn := FuncName(h.servingHook)
			m.Store(fn, true)
			defer stop(fn)
			defer func() { m.Delete(fn) }()
			h.servingHook(h.ctx, next)
			return
		}
		h.hook(h.ctx)
		next()
	})
}

type Lifecycle interface {
	Append(func(ctx context.Context))
	AppendServing(func(ctx context.Context, onStarted func()))
	WaitTerminateWithTimeout(time.Duration, func())
}

// Application 定义的 DI 容器
type Application interface {
	Provide(constructor interface{}, opts ...ProvideOption) // 注册一个依赖构造器
	Run(constructor interface{}, opts ...InvokeOption)      // 注册一个任务
	Wait() error                                            // 阻塞到所有任务完成.
}

// 别名, 不要直接调用 dig
type (
	ProvideOption = dig.ProvideOption
	InvokeOption  = dig.InvokeOption
	In            = dig.In
	Out           = dig.Out
)

// 别名, 不要直接调用 dig
var (
	Name            = dig.Name
	Group           = dig.Group
	As              = dig.As
	Export          = dig.Export
	FillProvideInfo = dig.FillProvideInfo
	LocationForPC   = dig.LocationForPC
)

type appImpl struct {
	container *dig.Container
	mu        sync.Mutex
	provides  []provideOption
	invokes   []invokeOption
}

func (app *appImpl) doInvokes() error {
	for _, it := range app.provides {
		if err := app.container.Provide(it.constructor, it.opts...); err != nil {
			return err
		}
	}
	for _, it := range app.invokes {
		if err := app.container.Invoke(it.constructor, it.opts...); err != nil {
			return err
		}
	}
	return nil
}

type provideOption struct {
	constructor interface{}
	opts        []ProvideOption
}
type invokeOption struct {
	constructor interface{}
	opts        []InvokeOption
}

func (app *appImpl) Provide(constructor interface{}, opts ...ProvideOption) {
	app.mu.Lock()
	defer app.mu.Unlock()

	app.provides = append(app.provides, provideOption{
		constructor: constructor,
		opts:        opts,
	})
}
func (app *appImpl) Run(constructor interface{}, opts ...InvokeOption) {
	app.mu.Lock()
	defer app.mu.Unlock()

	app.invokes = append(app.invokes, invokeOption{
		constructor: constructor,
		opts:        opts,
	})
}

func RootCause(err error) error {
	return dig.RootCause(err)
}
func (app *appImpl) Wait() error {
	startCtx, startCancel := context.WithCancel(context.Background())
	defer startCancel()
	stopCtx, stopCancel := context.WithCancel(context.Background())
	defer stopCancel()
	var once sync.Once
	stop := func(s string) {
		once.Do(func() {
			startCancel()
			Debugf("runtime: terminating caused by %s", s)
		})
	}

	builder := lifecycle{ctx: stopCtx}
	_ = app.container.Provide(func() Lifecycle {
		return &builder
	})
	_ = app.container.Provide(func() context.Context {
		return startCtx
	})

	err := app.doInvokes()
	if err != nil {
		return err
	}

	var wg sync.WaitGroup
	var m sync.Map
	go watchInterrupt(startCtx, stop, &m)
	builder.build(startCtx, &wg, stop, &m)()
	wg.Wait()
	stopCancel()
	builder.wgStop.Wait()

	return nil
}

func New() Application {
	//dig.DeferAcyclicVerification()
	app := &appImpl{
		container: dig.New(),
		provides:  []provideOption{},
		invokes:   []invokeOption{},
	}
	// app.Run(WatchInterrupt()) // todo options
	return app
}

type lifecycle struct {
	hooks  []*Hook
	ctx    context.Context
	wgStop sync.WaitGroup
}

func (l *lifecycle) Append(fn func(context.Context)) {
	l.hooks = append(l.hooks, &Hook{
		serving: false,
		hook:    fn,
	})
}
func (l *lifecycle) AppendServing(fn func(ctx context.Context, onStarted func())) {
	l.hooks = append(l.hooks, &Hook{
		serving:     true,
		servingHook: fn,
	})
}
func (l *lifecycle) WaitTerminateWithTimeout(t time.Duration, task func()) {
	l.wgStop.Add(1)
	ctx, cancel := context.WithTimeout(l.ctx, t)
	go func() {
		defer cancel()
		defer l.wgStop.Done()

		<-ctx.Done()
		task()
	}()
}

func (l *lifecycle) build(ctx context.Context, wg *sync.WaitGroup, stop func(s string), m *sync.Map) (next func()) {
	builder := l.hooks

	contexts := []*hookCtx{}
	for range builder {
		nextCtx, cancel := context.WithCancel(ctx)
		contexts = append(contexts, &hookCtx{
			ctx:    nextCtx,
			cancel: cancel,
		})
	}
	n := len(contexts) - 1
	// fmt.Printf("n: %v\n", n)
	for i, h := range builder {
		h.hookCtx = contexts[n-i]
		// if h.serving {
		// 	// wg.Add(1)
		// 	// h.hookCtx.ctx = context.WithValue(h.hookCtx.ctx, "hookName", funcName(h.servingHook))
		// }
	}

	var i = 0
	var mu sync.Mutex

	next = func() {
		mu.Lock()
		// fmt.Printf("i: %v\n", i)
		done := i > n
		select {
		case <-ctx.Done():
			done = true
		default:
		}
		if ctx.Err() != nil {
			done = true
		}
		if done {
			mu.Unlock()
			return
		}

		h := builder[i]
		//fn:=funcName(h.servingHook)
		i++
		wg.Add(1)
		mu.Unlock()

		go h.doRun(wg, next, stop, m)

	}

	return
}

func watchInterrupt(ctx context.Context, cancel func(s string), m *sync.Map, sig ...os.Signal) {
	if len(sig) == 0 {
		sig = []os.Signal{os.Interrupt, syscall.SIGHUP}
	}
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, sig...)

	var b os.Signal
	select {
	case b = <-ch:
	case <-ctx.Done():
		b = syscall.SIGPIPE // TODO 自定义退出
	}
	cancel(fmt.Sprintf("signal: %v", b))

	select {
	case <-ch:
		Debug("\nforce quit")
	case <-time.After(time.Minute * 1):
		Debug("terminating timeout 5m force quit")
		m.Range(func(key, value any) bool {
			Debugf("blocking terminated hook: %v", key)
			return false
		})
	}
	os.Exit(1)
}

func FuncName(f any) string {
	rv := reflect.ValueOf(f)
	if !rv.IsValid() || rv.Kind() != reflect.Func {
		return ""
	}
	return runtime.FuncForPC(rv.Pointer()).Name()
}
