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
	serving bool // 如果true,OnStart 结束, 则退出 整个app
	// OnStart     func(ctx context.Context) // app开始时触发
	// OnStop      func(ctx context.Context) // app结束时触发
	// Run         func(ctx context.Context, running func())
	// OnRuning    func()
	hook        func(ctx context.Context)
	servingHook ServingHookFunc
	// ctx         context.Context
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
	// AwaitTermination(ctx context.Context) error
	Wait() error // 阻塞到所有任务完成.
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

	mu sync.Mutex
	// running  bool
	// hooks    []*Hook
	provides []provideOption
	invokes  []invokeOption
	// r        atomic.Value
	cancel context.CancelFunc
}

func (app *appImpl) doInvokes() error {
	// fmt.Println("runtime: all invokes", len(app.provides)+len(app.invokes))
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

// func (app *appImpl) doExit(ctx context.Context, wg *sync.WaitGroup, builder hookBuilder) {
// 	<-ctx.Done()
// 	stopCtx, stopCancel := context.WithTimeout(context.Background(), time.Second*120)
// 	defer stopCancel()
// 	go func() {
// 		<-stopCtx.Done()
// 		if errors.Is(stopCtx.Err(), context.DeadlineExceeded) {
// 			fmt.Println("terminating timeout was forced to quit")
// 			os.Exit(1)
// 		}
// 	}()

// 	for _, hook := range builder {
// 		if hook.OnStop == nil {
// 			continue
// 		}
// 		wg.Add(1)
// 		go func(hook *Hook) {
// 			defer wg.Done()
// 			hook.OnStop(stopCtx)
// 		}(hook)
// 	}
// 	wg.Wait()
// }

type hookr struct {
	// run    func(int, *hookr)
	ctx    context.Context
	cancel context.CancelFunc
	// next   *hookr
	index int
	hook  *Hook
}

func (app *appImpl) Wait() error {
	builder := hookBuilder{}
	_ = app.container.Provide(func() Lifecycle {
		return &builder
	})

	err := app.doInvokes()
	if err != nil {
		// TODO return err
		return dig.RootCause(err)
	}
	// fmt.Println("runtime: all hooks", len(builder))

	var wg sync.WaitGroup
	// root, cancel := context.WithCancel(context.Background())
	// app.cancel = cancel

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

func (app *appImpl) WaitOld() error {
	builder := hookBuilder{}
	_ = app.container.Provide(func() Lifecycle {
		return &builder
	})

	err := app.doInvokes()
	if err != nil {
		// TODO return err
		return dig.RootCause(err)
	}
	// fmt.Println("runtime: all hooks", len(builder))

	var wg sync.WaitGroup
	root, cancel := context.WithCancel(context.Background())
	app.cancel = cancel

	startCtx, startCancel := context.WithCancel(root)
	defer startCancel()

	var startCancelOnce sync.Once

	// go app.doExit(startCtx, &wg, builder)

	hookrs := []*hookr{}

	var hookCtx context.Context
	var hookCancel context.CancelFunc
	for i, hook := range builder {
		if hook.serving {
			wg.Add(1)
		}
		if i == 0 {
			hookCtx, hookCancel = context.WithCancel(startCtx)
		} else {
			hookCtx, hookCancel = context.WithCancel(hookCtx)
		}
		hr := &hookr{
			ctx:    hookCtx,
			cancel: hookCancel,
			index:  i,
			hook:   hook,
		}
		hookrs = append(hookrs, hr)
	}

	var run func(i int, crt *hookr)
	run = func(i int, crt *hookr) {
		time.Sleep(time.Microsecond) // 稍稍延迟, 以尽量保持执行顺序
		var next *hookr
		if crt.index+1 < len(hookrs) {
			next = hookrs[crt.index+1]
		}
		ctx := crt.ctx
		cancel := crt.cancel
		// if next != nil {
		// 	ctx = next.ctx
		// 	cancel = next.cancel
		// }
		serving := crt.hook.serving
		defer cancel()
		defer func() {
			if serving {
				wg.Done()
				startCancel()
				startCancelOnce.Do(func() {
					go func() {
						timer := time.NewTimer(time.Minute * 5)
						defer timer.Stop()

						<-timer.C
						fmt.Println("terminating timeout was forced to quit")
						os.Exit(1)
					}()
				})
			}
		}()

		defer hookCancel()
		defer HandleCrash()

		if serving {
			crt.hook.servingHook(ctx, func() {
				if next == nil {
					return
				}
				go run(i+1, next)
			})
		} else {
			if next != nil {
				go run(i+1, next)
			}
			crt.hook.hook(ctx)

		}
	}
	hookrsCopy := make([]*hookr, 0)
	for n := len(hookrs) - 1; n >= 0; n-- {
		hookrsCopy = append(hookrsCopy, hookrs[n])
	}
	for i := range hookrs {
		hookrs[i].ctx = hookrsCopy[i].ctx
		hookrs[i].cancel = hookrsCopy[i].cancel
	}
	// fmt.Printf("hookrs: %v\n", hookrs)
	if len(hookrs) > 0 {
		go run(0, hookrs[0])
	}

	// for _, hook := range builder {
	// 	if hook.OnStart == nil {
	// 		continue
	// 	}
	// 	wg.Add(1)
	// 	go func(hook *Hook) {
	// 		defer wg.Done()
	// 		defer func() {
	// 			if hook.Serving {
	// 				startCancel()
	// 			}
	// 		}()
	// 		defer HandleCrash()
	// 		hook.OnStart(startCtx)
	// 	}(hook)
	// }
	wg.Wait()
	return nil
}

func New() Application {
	//dig.DeferAcyclicVerification()
	app := &appImpl{container: dig.New()}

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

func (hooks *hookBuilder) build(ctx context.Context, wg *sync.WaitGroup, stop func()) (next func()) {
	contexts := []*hookCtx{}
	builder := *hooks
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
