package apm

import (
	"runtime"
	"sync"

	"github.com/junhwong/goost/apm/field"
)

var (
	defaultEntry *factoryEntry
	initOnce     sync.Once
	gmu          sync.Mutex
	queue        chan *field.Field
	queuewg      sync.WaitGroup
)

func init() {
	initOnce.Do(func() {
		defaultEntry = &factoryEntry{calldepth: 1, Field: field.MakeRoot()}

		handler, _ := NewConsole()
		handler.HandlerPriority -= 999
		AddHandlers(handler)

		queue = make(chan *field.Field, runtime.NumCPU()*2)

		go loop()
	})
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
