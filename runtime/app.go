package runtime

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"sync/atomic"
	"syscall"
	"time"

	"go.uber.org/dig"
)

type Hook struct {
	Serving bool                      // 如果true,OnStart 结束, 则退出 整个app
	OnStart func(ctx context.Context) // app开始时触发
	OnStop  func(ctx context.Context) // app结束时触发
}

type Lifecycle interface {
	Append(Hook)
}

type Application interface {
	Provide(constructor interface{}, opts ...ProvideOption)
	Run(constructor interface{}, opts ...InvokeOption)
	AwaitTermination(ctx context.Context) error
	Wait() error // Wait 阻塞到所有任务完成.
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

	mu       sync.Mutex
	running  bool
	hooks    []*Hook
	provides []provideOption
	invokes  []invokeOption
	r        atomic.Value
	cancel   context.CancelFunc
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

func (app *appImpl) AwaitTermination(ctx context.Context) error {
	app.mu.Lock()
	// if app.running {
	// 	app.mu.Unlock()
	// 	return nil
	// }
	if app.r.CompareAndSwap(false, true) {
		app.mu.Unlock()
		return nil
	}
	app.running = true
	// app.mu.Unlock()

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	builder := hookBuilder{}
	_ = app.container.Provide(func() Lifecycle {
		return &builder
	})

	app.hooks = builder // []*Hook(builder)

	var wg sync.WaitGroup
	// for _, it := range app.hooks {
	// 	if it == nil {
	// 		continue
	// 	}
	// 	if it.Serving {
	// 		it.wg = &wg
	// 		wg.Add(1)
	// 	}
	// 	go it.doStart(ctx)
	// }
	go app.watchSignal()

	wg.Wait()
	app.stop(ctx) // nolint
	return nil
}
func (app *appImpl) watchSignal() {
	sigch := make(chan os.Signal, 1)
	signal.Notify(sigch,
		syscall.SIGHUP,  // 用户终端连接(正常或非正常)结束时发出, 通常是在终端的控制进程结束时, 通知同一session内的各个作业, 这时它们与控制终端不再关联。
		syscall.SIGINT,  // 程序终止(interrupt)信号, 在用户键入INTR字符(通常是Ctrl-C)时发出，用于通知前台进程组终止进程
		syscall.SIGTERM, // 程序结束(terminate)信号, 与SIGKILL不同的是该信号可以被阻塞和处理。通常用来要求程序自己正常退出，shell命令kill缺省产生这个信号
		syscall.SIGQUIT, // 和SIGINT类似, 但由QUIT字符(通常是Ctrl-\)来控制. 进程在因收到SIGQUIT退出时会产生core文件, 在这个意义上类似于一个程序错误信号
	)
	interruptCount := 0
FOR:
	for sig := range sigch {
		switch sig {
		case syscall.SIGINT:
			interruptCount++
			if interruptCount > 1 {
				fmt.Println("强制退出")
				os.Exit(1)
			}
			go app.stop(context.TODO()) // nolint
		default:
			break FOR
		}
	}

}
func (app *appImpl) stop(ctx context.Context) {
	// app.mu.Lock()
	// defer app.mu.Unlock()
	// if !app.running {
	// 	return
	// }
	if app.r.CompareAndSwap(true, false) {
		return
	}

	app.running = false
	ctx, cancel := context.WithTimeout(ctx, time.Minute)
	defer cancel()
	for _, it := range app.hooks {
		// go it.doStop(ctx)
		it.OnStart(ctx)
	}
	// fmt.Println("exited")
	os.Exit(0)
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
	}
	var wg sync.WaitGroup
	root, cancel := context.WithCancel(context.Background())
	app.cancel = cancel

	startCtx, startCancel := context.WithCancel(root)
	defer startCancel()
	for _, hook := range builder {
		if hook.OnStart == nil {
			continue
		}
		wg.Add(1)
		go func(hook *Hook) {
			defer wg.Done()
			defer func() {
				if hook.Serving {
					startCancel()
				}
			}()
			defer HandleCrash()
			hook.OnStart(startCtx)
		}(hook)
	}
	wg.Wait()
	startCancel()

	stopCtx, stopCancel := context.WithTimeout(context.Background(), time.Minute)
	defer stopCancel()
	for _, hook := range builder {
		if hook.OnStop == nil {
			continue
		}
		wg.Add(1)
		go func(hook *Hook) {
			defer wg.Done()
			hook.OnStop(stopCtx)
		}(hook)
	}
	wg.Wait()

	return nil
	// return app.AwaitTermination(context.TODO())
}

func New() Application {
	//dig.DeferAcyclicVerification()
	app := &appImpl{container: dig.New()}

	return app
}

type hookBuilder []*Hook

func (hooks *hookBuilder) Append(hook Hook) {
	arr := []*Hook(*hooks)
	target := hookBuilder(append(arr, &hook))
	*hooks = target
}
