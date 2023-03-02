package apm

import (
	"sync"
)

var (
	std *logImpl
	// defi Interface
	// asyncD Dispatcher
	dispatcher Dispatcher = &syncDispatcher{}
	initOnce   sync.Once
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

	initOnce.Do(func() {
		handler, _ := Console()
		dispatcher.AddHandlers(handler)
		std = &logImpl{calldepth: 1}
	})

	// asyncD = &asyncDispatcher{}
	// defi = New(context.Background())
}

func Flush() {
	dispatcher.Flush()
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
	return dispatcher
}
func SetDispatcher(a Dispatcher) {
	old := dispatcher
	defer old.Flush()

	handlers := dispatcher.GetHandlers()
	a.AddHandlers(handlers...)
	dispatcher = a
}

func UseAsyncDispatcher() {
	d := &asyncDispatcher{queue: make(chan Entry, 1024)}
	SetDispatcher(d)
	go d.loop()
}

func AddHandlers(handlers ...Handler) {
	dispatcher.AddHandlers(handlers...)
}

type Option interface {
	applyInterface(*logImpl)
}

// type funcSpanOption func(SpanOptionSetter)

func WithFields(fs ...Field) funcSpanOption {
	return funcSpanOption(func(appender SpanOptionSetter) {
		appender.SetAttributes(fs...)
	})
}

func Default() Interface {
	return std
}
