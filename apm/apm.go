package apm

import (
	"context"
)

var (
	std *logImpl
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

	provider := &syncDispatcher{}
	provider.AddHandlers(Console())
	std = &logImpl{
		ctx:        context.TODO(),
		dispatcher: provider,
		calldepth:  1,
	}
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
	WithFields(fs ...Field) Interface
	SpanFactory
}

func GetAdapter() Adapter {
	return std.dispatcher
}
func SetDispatcher(a Dispatcher) {
	old := std.dispatcher
	defer old.Flush()

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

func AddHandlers(handlers ...Handler) {
	std.dispatcher.AddHandlers(handlers...)
}

type Option interface {
	applyInterface(*logImpl)
}
type FieldsAppender interface {
	AppendFields(fs Fields)
}

func WithFields(fs ...Field) any {
	return func(appender FieldsAppender) {
		appender.AppendFields(fs)
	}
}

func Default() Interface {
	return std
}
