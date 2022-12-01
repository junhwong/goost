package runtime

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"go.uber.org/dig"
)

type Hook struct {
	*hookCtx
	once    sync.Once
	serving bool

	hook        func(ctx context.Context)
	servingHook ServingHookFunc
}
type hookCtx struct {
	ctx    context.Context
	cancel context.CancelFunc
}

func (h *Hook) doRun(wg *sync.WaitGroup, next func(), stop func()) {
	h.once.Do(func() {
		defer h.cancel()
		if h.servingHook != nil {
			defer wg.Done()
			defer stop()
			h.servingHook(h.ctx, next)
			return
		}
		h.hook(h.ctx)
		next()
	})
}

type ServingHookFunc func(ctx context.Context, onStarted func())

type Lifecycle interface {
	Append(func(ctx context.Context))
	AppendServing(ServingHookFunc)
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
	Name  = dig.Name
	Group = dig.Group
	As    = dig.As
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

func (app *appImpl) Wait() error {
	builder := hookBuilder{}
	_ = app.container.Provide(func() Lifecycle {
		return &builder
	})

	err := app.doInvokes()
	if err != nil {
		return err
		return dig.RootCause(err)
	}

	var wg sync.WaitGroup

	startCtx, startCancel := context.WithCancel(context.Background())
	defer startCancel()
	next := builder.build(startCtx, &wg, stop(startCancel))
	next()

	wg.Wait()
	return nil
}

func stop(startCancel context.CancelFunc) func() {
	var once sync.Once
	return func() {
		once.Do(func() {
			defer startCancel()
			go func() {
				timer := time.NewTimer(time.Minute * 5)
				defer timer.Stop()

				<-timer.C
				fmt.Println("terminating timeout was forced to quit")
				os.Exit(1)
			}()
		})
	}
}

func New() Application {
	//dig.DeferAcyclicVerification()
	app := &appImpl{
		container: dig.New(),
		provides:  []provideOption{},
		invokes:   []invokeOption{},
	}
	app.Run(WatchInterrupt()) // todo options
	return app
}

type hookBuilder []*Hook

func (hooks *hookBuilder) Append(fn func(context.Context)) {
	arr := []*Hook(*hooks)
	target := hookBuilder(append(arr, &Hook{
		serving: false,
		hook:    fn,
	}))
	*hooks = target
}
func (hooks *hookBuilder) AppendServing(fn ServingHookFunc) {
	arr := []*Hook(*hooks)
	target := hookBuilder(append(arr, &Hook{
		serving:     true,
		servingHook: fn,
	}))
	*hooks = target
}

func (hooks hookBuilder) build(ctx context.Context, wg *sync.WaitGroup, stop func()) (next func()) {
	builder := hooks

	contexts := []*hookCtx{}
	for range builder {
		next, cancel := context.WithCancel(ctx)
		contexts = append(contexts, &hookCtx{
			ctx:    next,
			cancel: cancel,
		})
	}
	n := len(contexts) - 1
	// fmt.Printf("n: %v\n", n)
	for i, h := range builder {
		h.hookCtx = contexts[n-i]
		if h.serving {
			wg.Add(1)
		}
	}

	var i = 0
	var mu sync.Mutex

	next = func() {
		mu.Lock()
		// fmt.Printf("i: %v\n", i)
		if i > n {
			mu.Unlock()
			return
		}

		h := builder[i]
		i++
		mu.Unlock()

		go h.doRun(wg, next, stop)

	}
	return
}

// 跟踪中断信号。
func WatchInterrupt(sig ...os.Signal) func(Lifecycle) {
	if len(sig) == 0 {
		sig = []os.Signal{os.Interrupt, syscall.SIGHUP}
	}
	ch := make(chan os.Signal, 1)
	return func(life Lifecycle) {
		signal.Notify(ch, sig...)
		life.AppendServing(func(ctx context.Context, onStarted func()) {
			onStarted()
			var b os.Signal
			select {
			case b = <-ch:
			case <-ctx.Done():
				b = syscall.SIGPIPE // TODO 自定义退出
			}

			// TODO 临时检测具体信号
			fmt.Println("runtime/WatchInterrupt: caught signal: ", b)
			fmt.Println("runtime/WatchInterrupt: shutting down")
			go func() {
				for sig := range ch {
					if sig == os.Interrupt {
						break
					}
				}
				fmt.Println("force quit")
				os.Exit(1)
			}()
		})
	}
}

func rootCancel(ctx context.Context, cancel context.CancelFunc) {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)
	select {
	case <-ctx.Done():
		close(ch)
		return
	case <-ch:
		cancel()
		return
	}
}
