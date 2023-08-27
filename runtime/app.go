package runtime

import (
	"context"
	"errors"
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
	servingHook func(ctx context.Context, next func())
}
type hookCtx struct {
	ctx    context.Context
	cancel context.CancelFunc
}

func (h *Hook) doRun(wg *sync.WaitGroup, next func(), stop func(s string), callerNames *sync.Map) {
	h.once.Do(func() {
		fn := FuncName(h.servingHook)
		callerNames.Store(fn, true)
		defer callerNames.Delete(fn)
		if h.serving {
			defer stop(fn)
		}
		defer wg.Done()
		defer h.cancel()
		h.servingHook(h.ctx, next)
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
	regs      []reg
	life      *lifecycle
	stop      func(string)
	err       error
}

func (app *appImpl) doInvokes() error {
	for _, v := range app.regs {
		if err := v.run(app.container); err != nil {
			return err
		}
	}
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

type reg interface {
	run(container *dig.Container) error
}
type provideOption struct {
	constructor interface{}
	opts        []ProvideOption
}

func (o *provideOption) run(container *dig.Container) error {
	return container.Provide(o.constructor, o.opts...)
}

type invokeOption struct {
	constructor interface{}
	opts        []InvokeOption
}

func (o *invokeOption) run(container *dig.Container) error {
	return container.Invoke(o.constructor, o.opts...)
}

func (app *appImpl) Provide(constructor interface{}, opts ...ProvideOption) {
	app.mu.Lock()
	if app.err != nil {
		app.mu.Unlock()
		return
	}
	app.mu.Unlock()
	err := app.container.Provide(constructor, opts...)
	if err != nil {
		app.mu.Lock()
		if app.err != nil {
			app.err = errors.Join(app.err, err)
		} else {
			app.err = err
		}
		app.mu.Unlock()
	}
}
func (app *appImpl) Run(constructor interface{}, opts ...InvokeOption) {
	app.mu.Lock()
	if app.err != nil {
		app.mu.Unlock()
		return
	}
	app.mu.Unlock()
	err := app.container.Invoke(constructor, opts...)
	if err != nil {
		app.mu.Lock()
		if app.err != nil {
			app.err = errors.Join(app.err, err)
		} else {
			app.err = err
		}
		app.mu.Unlock()
	}
}

func RootCause(err error) error {
	return dig.RootCause(err)
}
func (app *appImpl) Wait() error {
	app.mu.Lock()
	if app.err != nil {
		app.mu.Unlock()
		return app.err
	}
	app.mu.Unlock()
	life, stop := app.life, app.stop
	defer stop("")

	life.Wait()
	return nil
}

func New() Application {
	//dig.DeferAcyclicVerification()
	app := &appImpl{
		container: dig.New(),
		provides:  []provideOption{},
		invokes:   []invokeOption{},
	}
	app.life, app.stop = NewLifecycle(context.Background())
	_ = app.container.Provide(func() Lifecycle {
		return app.life
	})
	_ = app.container.Provide(func() context.Context {
		return app.life.Context()
	})
	go watchInterrupt(app.life.Context(), app.stop, app.life.CallerNames())
	// app.Run(WatchInterrupt()) // todo options
	return app
}

func NewLifecycle(ctx context.Context) (*lifecycle, func(string)) {
	_, startCancel := context.WithCancel(ctx)
	defer startCancel()

	stopCtx, stopCancel := context.WithCancel(context.Background())
	// defer stopCancel()
	var once sync.Once
	stop := func(s string) {
		once.Do(func() {
			// startCancel()
			stopCancel()
			if len(s) == 0 {
				return
			}
			Debugf("runtime: terminating caused by %s", s)
		})
	}

	return &lifecycle{ctx: stopCtx, startCancel: startCancel, stop: stop}, stop // stopCtx
}

type lifecycle struct {
	hooks       []*Hook
	ctx         context.Context
	wgStop      sync.WaitGroup // 等待结束
	wg          sync.WaitGroup // hook 组
	startCancel context.CancelFunc
	stop        func(s string)
	callerNames sync.Map
	mu          sync.Mutex
	started     bool
}

func (l *lifecycle) Append(fn func(context.Context)) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.hooks = append(l.hooks, &Hook{
		serving: false,
		servingHook: func(ctx context.Context, next func()) {
			fn(ctx)
			next()
		},
	})
}
func (l *lifecycle) AppendServing(fn func(context.Context, func())) {
	l.mu.Lock()
	defer l.mu.Unlock()
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
func (l *lifecycle) Context() context.Context {
	return l.ctx
}
func (l *lifecycle) CallerNames() *sync.Map {
	return &l.callerNames
}
func (l *lifecycle) Start() {
	l.mu.Lock()
	if l.started {
		l.mu.Unlock()
		return
	}
	l.started = true
	builder := l.hooks
	ctx := l.ctx
	wg := &l.wg
	stop := l.stop
	contexts := []*hookCtx{}
	//ctx context.Context, wg *sync.WaitGroup, stop func(s string), callerNames *sync.Map
	callerNames := l.CallerNames()

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
	}

	var i = 0
	// var mu sync.Mutex
	var next func()
	next = func() {
		l.mu.Lock()
		if i > n {
			l.mu.Unlock()
			return
		}
		select {
		case <-ctx.Done():
			l.mu.Unlock()
			return
		default:
		}

		h := builder[i]
		//fn:=funcName(h.servingHook)
		i++
		wg.Add(1)
		l.mu.Unlock()

		go h.doRun(wg, next, stop, callerNames)

	}
	l.mu.Unlock()
	next()

	return
}

func (l *lifecycle) Wait() {
	l.Start()
	l.wg.Wait()
	l.stop("")
	l.wgStop.Wait()
}

func watchInterrupt(ctx context.Context, stop func(s string), callerNames *sync.Map, sig ...os.Signal) {
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
	stop(fmt.Sprintf("signal: %v", b))

	select {
	case <-ch:
		Debug("\nforce quit")
	case <-time.After(time.Minute * 1):
		Debug("terminating timeout 5m force quit")
		callerNames.Range(func(key, value any) bool {
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
