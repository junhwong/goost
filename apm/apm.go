package apm

import (
	"sync"

	"github.com/junhwong/goost/apm/field"
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
		std = &FieldsEntry{calldepth: 1}
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
	WithFields(fs ...*field.Field) Interface
	CalldepthInc()
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

func Default() Interface {
	return std
}
