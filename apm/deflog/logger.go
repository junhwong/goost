package deflog

import (
	"context"
	"sync"
	"time"

	"github.com/junhwong/goost/apm"
)

func NewLogger() *DefaultLogger {
	return &DefaultLogger{
		queue:    make(chan apm.Entry, 1024),
		inqueue:  make(chan apm.Entry, 1024),
		handlers: []apm.Handler{},
	}
}

type DefaultLogger struct {
	mu       sync.RWMutex
	wg       sync.WaitGroup
	queue    chan apm.Entry
	inqueue  chan apm.Entry
	cancel   context.CancelFunc
	handlers handlerSlice
}

func (logger *DefaultLogger) AddHandlers(handlers ...apm.Handler) {
	logger.mu.Lock()
	old := logger.handlers
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
func (logger *DefaultLogger) gethandlers() handlerSlice {
	handlers := make(handlerSlice, logger.handlers.Len())
	copy(handlers, logger.handlers)
	return handlers
}

func (logger *DefaultLogger) handle(entry apm.Entry) {
	// std.mu.Lock()
	handlers := logger.handlers
	// std.mu.Unlock()

	handlers.handle(entry)
}

func (logger *DefaultLogger) Flush() error {
	// std.mu.Lock()
	// defer std.mu.Unlock()

	// handlers := logger.handlers
	// for entry := range logger.inqueue {
	// 	handlers.handle(entry)
	// 	logger.wg.Done()
	// }
	return nil
}

func (logger *DefaultLogger) Run(stopCh <-chan struct{}) {
	for {
		select {
		case entry, ok := <-logger.queue:
			if !ok {
				continue
			}

			logger.wg.Add(1)
			logger.inqueue <- entry
		case entry, ok := <-logger.inqueue:
			if !ok {
				return
			}

			logger.handle(entry)
			logger.wg.Done()
		case <-stopCh:
			goto END
		}
	}

END:
	logger.Flush()
}

func (logger *DefaultLogger) Close() error {
	time.Sleep(time.Millisecond) // 给协程一点时间启动
	logger.cancel()
	go logger.Flush()
	close(logger.queue)
	logger.wg.Wait()
	close(logger.inqueue)
	return nil
}

// func (logger *DefaultLogger) NewSpan(ctx context.Context, options ...apm.SpanOption) (context.Context, apm.Span) {
// 	calldepth := 0
// 	return newSpan(ctx, logger, calldepth+1, options)
// }

// func (entry *DefaultLogger) Logf(ctx context.Context, calldepth int, level apm.LogLevel, format string, args []interface{}) {
// 	fs := make(field.Fields, 5)

// 	var err error
// 	a := []interface{}{}
// 	ctxs := []context.Context{ctx}
// 	for _, f := range args {
// 		switch f := f.(type) {
// 		case field.Field:
// 			fs.Set(f)
// 		case context.Context:
// 			ctxs = append(ctxs, f)
// 		case error:
// 			a = append(a, f)
// 			err = f
// 		default:
// 			a = append(a, f)
// 		}
// 		// if fd, ok := f.(field.Field); ok {
// 		// 	fs.Set(fd)
// 		// } else {
// 		// 	a = append(a, f)
// 		// 	if ex, ok := f.(error); ok {
// 		// 		err = ex
// 		// 	}
// 		// }
// 	}

// 	if info, ok := runtime.GetCallLastFromError(err); ok {
// 		fs.Set(apm.ErrorMethod(info.Method))
// 	}

// 	if stack, ok := runtime.GetCallStackFromError(err); ok {
// 		arr := []string{}
// 		for _, it := range stack {
// 			caller := it.Path
// 			if i := strings.LastIndex(caller, "/"); i > 0 {
// 				caller = caller[i+1:]
// 			}
// 			arr = append(arr, caller+"/"+it.File+":"+cast.ToString(it.Line))
// 		}
// 		fs.Set(apm.ErrorMethod(strings.Join(arr, ",")))
// 	}

// 	if calldepth > -1 {
// 		info := runtime.Caller(calldepth + 1)
// 		fs.Set(apm.TracebackCaller(getSplitLast(info.Method, "/")))
// 		fs.Set(apm.TracebackLineNo(info.Line))
// 		// fmt.Printf("info.Path: %v\n", info.Path)
// 		p := info.Path
// 		if i := strings.LastIndex(p, "/"); i > 0 {
// 			p = p[i+1:]
// 		}
// 		if len(p) != 0 {
// 			p += "/"
// 		}

// 		fs.Set(apm.TracebackPath(p + info.File))
// 	}

// 	if _, ok := fs[apm.TraceIDKey]; !ok {
// 		for _, ctx := range ctxs {
// 			tid, sid := apm.GetTraceID(ctx)
// 			if len(tid) > 0 {
// 				fs.Set(apm.TraceID(tid))
// 				if _, ok := fs[apm.SpanIDKey]; !ok {
// 					if len(sid) > 0 {
// 						fs.Set(apm.SpanID(sid))
// 					}
// 				}
// 				break
// 			}
// 		}
// 	}

// 	fs.Set(apm.Message(format, a...))
// 	fs.Set(apm.Level(level))
// 	if _, ok := fs[apm.TimeKey]; !ok {
// 		fs.Set(apm.Time(time.Now()))
// 	}

// 	// TODO 异步
// 	// entry.queue <- Entry(fs)
// 	entry.handlers.handle(apm.Entry(fs))
// }

// func (entry *DefaultLogger) Log(ctx context.Context, calldepth int, level apm.LogLevel, args []interface{}) {
// 	entry.Logf(ctx, calldepth+1, level, "", args)
// }

//	func getSplitLast(s string, substr string) string {
//		i := strings.LastIndex(s, substr)
//		if i > 0 {
//			s = s[i+1:]
//		}
//		return s
//	}
func (logger *DefaultLogger) Log(e apm.Entry) {
	logger.handlers.handle(e)
}
