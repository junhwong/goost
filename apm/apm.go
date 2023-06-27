package apm

import (
	"sync"
)

var (
	std *FieldsEntry
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
		handler.HandlerPriority -= 999
		dispatcher.AddHandlers(handler)
		std = &FieldsEntry{calldepth: 1} // 0 Default() ok
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
	SpanFactory
	// With(options ...WithOption) Interface
}

type WithOption interface {
	applyWithOption(*FieldsEntry)
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

// type Option interface {
// 	applyInterface(*FieldsEntry)
// }

// type attributesSetter interface {
// 	SetAttributes(a ...*field.Field)
// }
// type funcSetAttrsOption func(attributesSetter)

// func (f funcSetAttrsOption) applySpanOption(target *spanImpl) {
// 	f(target)
// }

func Default(options ...WithOption) Interface {
	if len(options) == 0 {
		return std
	}
	cl := std.new()
	// cl.calldepth++
	for _, o := range options {
		if o != nil {
			o.applyWithOption(cl)
		}
	}
	return cl
}
