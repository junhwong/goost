package apm

import (
	"sync"
	"sync/atomic"
)

type Dispatcher interface {
	AddHandlers(handlers ...Handler)
	RemoveHandlers(handlers ...Handler)
	GetHandlers() []Handler
	Dispatch(e Entry)
	Flush() error
	// Close() error
}

type syncDispatcher struct {
	mu       sync.RWMutex
	handlers atomic.Value
}

func (d *syncDispatcher) AddHandlers(handlers ...Handler) {
	if len(handlers) == 0 {
		panic("apm: handlers cannot empty")
	}

	d.mu.Lock()
	defer d.mu.Unlock()

	dst := d.getHandlers()
	for _, it := range handlers {
		if it == nil {
			continue
		}
		dst = append(dst, it)
	}
	dst.Sort()
	d.handlers.Store(dst)
}
func (d *syncDispatcher) RemoveHandlers(handlers ...Handler) {
	if len(handlers) == 0 {
		return
	}

	d.mu.Lock()
	defer d.mu.Unlock()

	var dst handlerSlice
	for _, v := range d.getHandlers() {
		found := false
		for _, it := range handlers {
			if it == v {
				found = true
				break
			}
		}
		if !found {
			dst = append(dst, v)
		}
	}
	dst.Sort()
	d.handlers.Store(dst)
}

func (d *syncDispatcher) getHandlers() handlerSlice {
	obj := d.handlers.Load()
	if obj == nil {
		return handlerSlice{}
	}
	return obj.(handlerSlice)
}

func (d *syncDispatcher) GetHandlers() []Handler {
	d.mu.Lock()
	old := d.getHandlers()
	d.mu.Unlock()
	return old
}

func (d *syncDispatcher) Dispatch(entry Entry) {
	if entry == nil {
		return
	}
	handlers := d.getHandlers()
	if len(handlers) == 0 {
		return
	}

	size := int32(handlers.Len())
	var crt atomic.Int32
	var once sync.Once
	var release = func() {
		once.Do(func() {
			// todo 将entry释放
		})
	}
	// defer release()?

	var next func()
	next = func() {
		i := crt.Add(1) - 1
		if i >= size {
			release()
			return
		}
		h := handlers[i]
		h.Handle(entry, next, release)
	}
	next()
}

func (d *syncDispatcher) Flush() error {
	return nil
}
func (d *syncDispatcher) Close() error {
	return nil
}

type asyncDispatcher struct {
	syncDispatcher
	queue chan Entry
	once  sync.Once
}

func (d *asyncDispatcher) Dispatch(e Entry) {
	if e == nil {
		return
	}

	// d.mu.RLock()
	// defer d.mu.RUnlock()

	d.queue <- e
}

func (d *asyncDispatcher) doFlush() error {

	for {
		select {
		case e, ok := <-d.queue:
			if !ok {
				return nil
			}
			d.syncDispatcher.Dispatch(e)
		default:
			return nil
		}
	}

}
func (d *asyncDispatcher) Flush() error {
	d.mu.Lock()
	defer d.mu.Unlock()

	return d.doFlush()
}

func (d *asyncDispatcher) Close() error {
	d.mu.Lock()
	defer d.mu.Unlock()

	close(d.queue)

	return d.doFlush()
}

func (d *asyncDispatcher) loop() {
	d.once.Do(func() {
		for e := range d.queue {
			if e == nil {
				return
			}
			d.syncDispatcher.Dispatch(e)
		}
	})
}
