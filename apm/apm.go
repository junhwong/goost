package apm

import (
	"sync"
	"sync/atomic"

	"github.com/junhwong/goost/apm/field"
)

var (
	defaultEntry *FieldsEntry
	// defi Interface
	// asyncD Dispatcher
	dispatcher atomic.Value //Dispatcher = &syncDispatcher{}
	initOnce   sync.Once
	gmu        sync.Mutex
)

func init() {
	initOnce.Do(func() {
		defaultEntry = &FieldsEntry{calldepth: 1, Field: *field.NewRoot()}

		handler, _ := NewConsole()
		handler.HandlerPriority -= 999
		d := &syncDispatcher{}
		d.AddHandlers(handler)

		dispatcher.Store(d)

	})
	var a atomic.Value
	a.Store(dispatcher)
	a.Load()
}

func GetDispatcher() Dispatcher {
	obj := dispatcher.Load()
	if obj == nil {
		return nil
	}
	return obj.(Dispatcher)
}

func Flush() {
	if d := GetDispatcher(); d != nil {
		d.Flush()
	}
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
	return GetDispatcher()
}
func SetDispatcher(d Dispatcher) {
	gmu.Lock()
	defer gmu.Unlock()

	old := GetDispatcher()
	if old != nil {
		defer old.Flush()

		handlers := old.GetHandlers()
		d.AddHandlers(handlers...)
	}
	dispatcher.Store(d)
}
func AddHandlers(handlers ...Handler) {
	if d := GetDispatcher(); d != nil {
		d.AddHandlers(handlers...)
	}
}
func RemoveHandlers(handlers ...Handler) {
	if d := GetDispatcher(); d != nil {
		d.RemoveHandlers(handlers...)
	}
}

func SetHandlers(handlers ...Handler) {
	if d := GetDispatcher(); d != nil {
		d.RemoveHandlers(d.GetHandlers()...)
		d.AddHandlers(handlers...)
	}
}

func UseAsyncDispatcher() {
	gmu.Lock()
	defer gmu.Unlock()

	d := &asyncDispatcher{queue: make(chan Entry, 1024)}
	SetDispatcher(d)
	go d.loop()
}

func Default(options ...WithOption) Interface {
	if len(options) == 0 {
		return defaultEntry
	}
	cl := defaultEntry.new()
	// cl.calldepth++
	for _, o := range options {
		if o != nil {
			o.applyWithOption(cl)
		}
	}
	return cl
}
