package runtime

import (
	"container/list"
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"reflect"
	"runtime"
	"strconv"
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
	hook  func(ctx context.Context)
	delay time.Duration
}

// func (h *Hook) doRun(wg *sync.WaitGroup, next func(), stop func(s string), callerNames *sync.Map) {
// 	h.once.Do(func() {
// 		fn := FuncName(h.servingHook)
// 		callerNames.Store(fn, true)
// 		defer callerNames.Delete(fn)
// 		defer stop(fn)
// 		defer wg.Done()
// 		defer h.cancel()
// 		h.servingHook(h.ctx, next)
// 	})
// }

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
	life      *lifecycle
	err       error
}

func (app *appImpl) Provide(constructor interface{}, opts ...ProvideOption) {
	app.life.mu.Lock()
	b := app.life.aborted
	app.life.mu.Unlock()
	if b {
		return
	}

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
	app.life.mu.Lock()
	b := app.life.aborted
	app.life.mu.Unlock()
	if b {
		return
	}

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
	app.life.Wait()
	return nil
}

func New() Application {
	//dig.DeferAcyclicVerification()
	app := &appImpl{
		container: dig.New(),
	}
	app.life, _ = NewLifecycle(context.Background())
	_ = app.container.Provide(func() Lifecycle {
		return app.life
	})
	_ = app.container.Provide(func() context.Context {
		return app.life
	})
	go func() {
		ch := make(chan os.Signal)
		signal.Notify(ch, os.Interrupt, syscall.SIGQUIT, syscall.SIGHUP)

		select {
		case b := <-ch:
			app.life.stop(fmt.Sprintf("%v signal", b))
		case <-app.life.Done():
			app.life.stop(fmt.Sprintf(app.life.Err().Error()))
		}

		select {
		case <-ch:
			Debug("\nforce quit")
			os.Exit(1)
		case <-time.After(time.Minute * 5):
			Debug("terminating timeout(5m), force quit")
			if s := app.life.LastCallerName(); len(s) > 0 {
				Debugf("blocking terminated hook: %v", s)
			}
			os.Exit(1)
		}

	}()
	return app
}

type Lifecycle interface {
	context.Context
	Append(func(ctx context.Context), ...HookAppendOption)
	WaitTerminateWithTimeout(time.Duration, func())
}

func NewLifecycle(ctx context.Context) (*lifecycle, func()) {
	ctx, startCancel := context.WithCancel(ctx)
	l := &lifecycle{Context: ctx, cc: list.New(), done: make(chan struct{})}
	stop := func(s string) {
		l.mu.Lock()
		if l.aborted {
			l.mu.Unlock()
			return
		}
		l.aborted = true
		n := len(l.cancels)
		l.mu.Unlock()

		go func() {
			defer startCancel()
			if len(s) != 0 {
				Debugf("runtime: terminating caused by %s", s)
			}

			for n > 0 {
				l.mu.Lock()
				n--
				ctx := l.cancels[n]
				l.mu.Unlock()
				ctx.cancel()
				select {
				case <-ctx.done:
				case <-time.After(time.Minute * 1):
					Debugf("runtime: waiting for %s to done timeout", ctx.name)
				}
			}
			l.mu.Lock()
			close(l.done)
			l.mu.Unlock()
		}()
	}

	l.stop = stop
	return l, func() { stop("") } // stopCtx
}

type hookCtx struct {
	context.Context
	cancel context.CancelFunc
	index  int
	name   string
	done   chan struct{}
	next   func()
}

func (c *hookCtx) Value(k any) any {
	if k == nextHookKey {
		return c.next
	}
	return c.Context.Value(k)
}

type lifecycle struct {
	context.Context

	hooks   []*Hook
	wgStop  sync.WaitGroup // 等待结束
	wg      sync.WaitGroup // hook 组
	stop    func(s string)
	mu      sync.Mutex
	started bool
	timeout time.Duration
	cancels []*hookCtx
	cc      *list.List
	done    chan struct{}
	aborted bool
}

type HookAppendOption func(*Hook)

var (
	nextHookKey = struct{ bool }{}
)

func CallNextHook(ctx context.Context) {
	next, _ := ctx.Value(nextHookKey).(func())
	if next != nil {
		next()
	}
}

func WithDelayCallNext(delay time.Duration) HookAppendOption {
	if delay <= 0 {
		panic("runtime: delay must be >0")
	}
	return func(h *Hook) {
		h.delay = delay
	}
}

func WithRunAfterCallNext() HookAppendOption {
	return func(h *Hook) {
		h.delay = 0
	}
}
func WithManualCallNext() HookAppendOption {
	return func(h *Hook) {
		h.delay = -1
	}
}

func (l *lifecycle) Append(run func(context.Context), opts ...HookAppendOption) {
	l.mu.Lock()
	defer l.mu.Unlock()

	h := &Hook{delay: time.Millisecond, hook: run}

	for _, o := range opts {
		if o != nil {
			o(h)
		}
	}

	l.hooks = append(l.hooks, h)
}
func (l *lifecycle) WaitTerminateWithTimeout(t time.Duration, task func()) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.wgStop.Add(1)

	go func() {
		<-l.done
		ctx, cancel := context.WithTimeout(context.TODO(), t)
		defer cancel()
		defer l.wgStop.Done()
		go func() {
			defer cancel()
			task()
		}()
		<-ctx.Done()
		if errors.Is(ctx.Err(), context.DeadlineExceeded) {
			Debugf("runtime: waiting for %s to done timeout", FuncName(task, true))
		}

	}()
}

func (l *lifecycle) LastCallerName() string {
	l.mu.Lock()
	last := l.cc.Back()
	l.mu.Unlock()
	if last != nil {
		s, _ := last.Value.(string)
		return s
	}
	return ""
}

func (l *lifecycle) Start() {
	l.mu.Lock()
	if l.started {
		l.mu.Unlock()
		return
	}
	l.started = true

	ctx := l
	wg := &l.wg
	stop := l.stop

	var i = 0
	var next func()
	next = func() {
		l.mu.Lock()
		if l.aborted || i >= len(l.hooks) {
			l.mu.Unlock()
			return
		}
		select {
		case <-ctx.Done():
			l.mu.Unlock()
			return
		default:
		}

		h := l.hooks[i]
		hook := h.hook
		var once sync.Once
		ctx, cancel := context.WithCancel(context.TODO())
		hctx := &hookCtx{
			index:   i,
			Context: ctx,
			cancel:  cancel,
			name:    FuncName(hook, true),
			done:    make(chan struct{}),
			next:    func() { once.Do(next) },
		}

		i++
		wg.Add(1)
		l.mu.Unlock()

		go func(ctx *hookCtx) {
			defer wg.Done()

			l.mu.Lock()
			l.cancels = append(l.cancels, hctx)
			el := l.cc.PushBack(hctx.name)
			l.mu.Unlock()

			defer func() {
				l.mu.Lock()
				l.cc.Remove(el)
				close(ctx.done)
				ctx.cancel()
				l.mu.Unlock()
				Debugf("todo: end %s: %d", ctx.name, ctx.index)
			}()

			Debugf("todo: run %s: %d", ctx.name, ctx.index)
			defer stop(ctx.name)

			delayCancel := func() {}
			if h.delay > 0 { // 延时调用 next
				waitCtx, cancel := context.WithTimeout(ctx, h.delay)
				delayCancel = cancel
				go func() {
					<-waitCtx.Done()
					ctx.next()
				}()
			}

			hook(ctx)

			delayCancel()

			if h.delay > 0 {
				ctx.next()
			}
		}(hctx)
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

func FuncName(f any, includeLine ...bool) string {
	rv := reflect.ValueOf(f)
	if !rv.IsValid() || rv.Kind() != reflect.Func {
		return ""
	}
	pc := rv.Pointer()
	rf := runtime.FuncForPC(pc)
	s := rf.Name()
	if len(includeLine) > 0 && includeLine[len(includeLine)-1] {
		_, l := rf.FileLine(pc)
		s += ":" + strconv.Itoa(l)
	}

	return s
}
