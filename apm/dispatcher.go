package apm

import "sync"

type Dispatcher interface {
	AddHandlers(handlers ...Handler)
	GetHandlers() []Handler
	Dispatch(e Entry)
	Flush() error
	// Close() error
}

type syncDispatcher struct {
	mu       sync.RWMutex
	handlers handlerSlice
}

func (logger *syncDispatcher) AddHandlers(handlers ...Handler) {
	logger.mu.Lock()
	old := logger.gethandlers()
	logger.mu.Unlock()

	for _, it := range handlers {
		if it == nil {
			continue
		}
		old = append(old, it)
	}
	old.Sort()
	logger.mu.Lock()
	logger.handlers = old
	logger.mu.Unlock()
}

func (logger *syncDispatcher) gethandlers() handlerSlice {
	handlers := make(handlerSlice, logger.handlers.Len())
	copy(handlers, logger.handlers)
	return handlers
}
func (logger *syncDispatcher) GetHandlers() []Handler {
	logger.mu.Lock()
	old := logger.gethandlers()
	logger.mu.Unlock()
	return old
}

func (logger *syncDispatcher) Dispatch(e Entry) {
	logger.handlers.handle(e)
}

func (logger *syncDispatcher) Flush() error {
	return nil
}
func (logger *syncDispatcher) Close() error {
	return nil
}

type asyncDispatcher struct {
	syncDispatcher
	queue chan Entry
	once  sync.Once
}

func (logger *asyncDispatcher) Dispatch(e Entry) {
	logger.mu.RLock()
	defer logger.mu.RUnlock()

	logger.queue <- e
}
func (logger *asyncDispatcher) doFlush() error {

	for {
		select {
		case e, ok := <-logger.queue:
			if !ok {
				return nil
			}
			logger.syncDispatcher.Dispatch(e)
		default:
			return nil
		}
	}

}
func (logger *asyncDispatcher) Flush() error {
	logger.mu.Lock()
	defer logger.mu.Unlock()

	return logger.doFlush()
}

func (logger *asyncDispatcher) Close() error {
	logger.mu.Lock()
	defer logger.mu.Unlock()

	close(logger.queue)

	return logger.doFlush()
}

func (logger *asyncDispatcher) loop() {
	logger.once.Do(func() {
		for e := range logger.queue {
			if e == nil {
				return
			}
			logger.syncDispatcher.Dispatch(e)
		}
	})
}
