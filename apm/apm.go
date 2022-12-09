package apm

import (
	"context"
)

var (
	std *loggerImpl
	// defi Interface
	// asyncD Dispatcher
)

func init() {
	// ctx, cancel := context.WithCancel(context.Background())
	// f := NewTextFormatter() // NewJsonFormatter() //
	// std = &DefaultLogger{
	// 	queue:    make(chan Entry, 1024),
	// 	inqueue:  make(chan Entry, 1024),
	// 	handlers: []Handler{&ConsoleHandler{Formatter: f}},
	// 	cancel:   cancel,
	// }
	// go std.Run(ctx.Done())
	// defi = New(context.Background())

	std = &loggerImpl{logImpl: logImpl{ctx: context.TODO(), dispatcher: &syncDispatcher{}, calldepth: 1}}
	// asyncD = &asyncDispatcher{}
	// defi = New(context.Background())
}

func Done() {
	// std.Close()
}
func Flush() {
	std.dispatcher.Flush()
}

// 适配接口
type Adapter interface {
	Dispatch(Entry)
}

// 同一接口
type Interface interface {
	Logger
	SpanFactory
}

func GetAdapter() Adapter {
	return std.dispatcher
}
func SetDispatcher(a Dispatcher) {
	old := std.dispatcher
	defer old.Close()

	handlers := std.dispatcher.GetHandlers()
	a.AddHandlers(handlers...)
	std.dispatcher = a
}

func UseAsyncDispatcher() {
	d := &asyncDispatcher{queue: make(chan Entry, 1024)}
	SetDispatcher(d)
	go d.loop()
}

func SetDefault(writer LoggerInterface) Adapter {
	// std = writer
	// defi = New(context.Background())
	// return defi
	return nil
}

func Default() Interface {
	return std
}

func AddHandlers(handlers ...Handler) {
	std.dispatcher.AddHandlers(handlers...)
}

type Option interface {
	applyInterface(*loggerImpl)
}

// func WithFields(fs ...field.Field) appendFields {
// 	return func() []field.Field {
// 		return fs
// 	}
// }

// func New(ctx context.Context, options ...Option) Interface {
// 	r := &loggerImpl{logImpl: logImpl{ctx: ctx, dispatcher: std, calldepth: 1}}
// 	return r
// }

// type Provider interface {
// 	Out(Entry)
// 	NewLogger() Logger
// 	NewSpan() Span
// }

// type Outer interface {
// 	Out(Entry)
// }

// type syncOuter struct {
// }

// func (syncOuter) Out(Entry) {

// }
