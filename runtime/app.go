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
	serving bool // 如果true,OnStart 结束, 则退出 整个app
	// OnStart     func(ctx context.Context) // app开始时触发
	// OnStop      func(ctx context.Context) // app结束时触发
	// Run         func(ctx context.Context, running func())
	// OnRuning    func()
	hook        func(ctx context.Context)
	servingHook ServingHookFunc
	// ctx         context.Context
}

type ServingHookFunc func(ctx context.Context, onStarted func())

//
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

// func (app *appImpl) watchSignal() {
// 	sigch := make(chan os.Signal, 1)
// 	signal.Notify(sigch,
// 		syscall.SIGHUP,  // 用户终端连接(正常或非正常)结束时发出, 通常是在终端的控制进程结束时, 通知同一session内的各个作业, 这时它们与控制终端不再关联。
// 		syscall.SIGINT,  // 程序终止(interrupt)信号, 在用户键入INTR字符(通常是Ctrl-C)时发出，用于通知前台进程组终止进程
// 		syscall.SIGTERM, // 程序结束(terminate)信号, 与SIGKILL不同的是该信号可以被阻塞和处理。通常用来要求程序自己正常退出，shell命令kill缺省产生这个信号
// 		syscall.SIGQUIT, // 和SIGINT类似, 但由QUIT字符(通常是Ctrl-\)来控制. 进程在因收到SIGQUIT退出时会产生core文件, 在这个意义上类似于一个程序错误信号
// 	)
// 	interruptCount := 0
// FOR:
// 	for sig := range sigch {
// 		switch sig {
// 		case syscall.SIGINT:
// 			interruptCount++
// 			if interruptCount > 1 {
// 				fmt.Println("强制退出")
// 				os.Exit(1)
// 			}
// 		default:
// 			break FOR
// 		}
// 	}
// }

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

	var wg sync.WaitGroup
	root, cancel := context.WithCancel(context.Background())
	app.cancel = cancel

	startCtx, startCancel := context.WithCancel(root)
	defer startCancel()

	var startCancelOnce sync.Once

	// go app.doExit(startCtx, &wg, builder)

	hookrs := []*hookr{}

	// var prevCtx context.Context
	// var prevCancel context.CancelFunc
	// var nextCtx context.Context
	// var nextCancel context.CancelFunc
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

// func (hooks *hookBuilder) Append(hook Hook) {
// 	fmt.Println("Append")
// 	arr := []*Hook(*hooks)
// 	target := hookBuilder(append(arr, &hook))
// 	*hooks = target
// }
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

// 跟踪中断信号。
func WatchInterrupt(sig ...os.Signal) func(Lifecycle) {
	if len(sig) == 0 {
		// systemd restart terminated , syscall.SIGTERM 子进程?
		sig = []os.Signal{os.Interrupt, syscall.SIGHUP}
	}
	ch := make(chan os.Signal, 1)
	return func(life Lifecycle) {
		// life.Append(Hook{
		// 	Serving: true,
		// 	OnStart: func(ctx context.Context) {
		// 		signal.Notify(ch, sig...)
		// 		b := <-ch
		// 		// TODO 临时检测具体信号
		// 		fmt.Println("runtime/WatchInterrupt: 收到信号", b)
		// 		go func() {
		// 			for sig := range ch {
		// 				if sig == os.Interrupt {
		// 					break
		// 				}
		// 			}
		// 			fmt.Println("force quit")
		// 			os.Exit(1)
		// 		}()
		// 	},
		// 	OnStop: func(ctx context.Context) {
		// 		ch <- syscall.SIGPIPE // 自定义退出
		// 	},
		// })
		life.Append(func(ctx context.Context) {
			// go func() {
			// 	running()
			// 	<-ctx.Done()
			// 	ch <- syscall.SIGPIPE // 自定义退出
			// }()
			signal.Notify(ch, sig...)
			// running()
			var b os.Signal
			select {
			case b = <-ch:
			case <-ctx.Done():
				b = syscall.SIGPIPE // TODO
			}

			// TODO 临时检测具体信号
			fmt.Println("runtime/WatchInterrupt: 收到信号", b)
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
