package apm

import (
	"context"
	"strings"
	"sync"
	"time"

	"github.com/junhwong/goost/apm/level"
	"github.com/junhwong/goost/pkg/field"
	"github.com/junhwong/goost/runtime"
	"github.com/spf13/cast"
)

// LoggerInterface 日志记录接口
type LoggerInterface interface {
	Log(ctx context.Context, calldepth int, level level.Level, args []interface{})
	Logf(ctx context.Context, calldepth int, level level.Level, format string, args []interface{})
}

// Logger 日志操作接口
type Logger interface {
	Debug(...interface{})
	Info(...interface{})
	Warn(...interface{})
	Error(...interface{})
	Fatal(...interface{})

	// Deprecated
	Debugf(string, ...interface{})
	// Deprecated
	Infof(string, ...interface{})
	// Deprecated
	Warnf(string, ...interface{})
	// Deprecated: use Error
	Errorf(string, ...interface{})
	// Deprecated:
	Fatalf(string, ...interface{})
}

type DefaultLogger struct {
	mu       sync.RWMutex
	wg       sync.WaitGroup
	queue    chan Entry
	inqueue  chan Entry
	cancel   context.CancelFunc
	handlers handlerSlice
}

func (logger *DefaultLogger) gethandlers() handlerSlice {
	handlers := make(handlerSlice, logger.handlers.Len())
	copy(handlers, logger.handlers)
	return handlers
}

func (logger *DefaultLogger) handle(entry Entry) {
	// std.mu.Lock()
	handlers := logger.handlers
	// std.mu.Unlock()

	handlers.handle(entry)
}

func (logger *DefaultLogger) Flush() {
	// std.mu.Lock()
	// defer std.mu.Unlock()

	// handlers := logger.handlers
	// for entry := range logger.inqueue {
	// 	handlers.handle(entry)
	// 	logger.wg.Done()
	// }
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

func (logger *DefaultLogger) NewSpan(ctx context.Context, calldepth int, options ...SpanOption) (context.Context, Span) {
	return newSpan(ctx, logger, calldepth+1, options)
}

func (entry *DefaultLogger) Logf(ctx context.Context, calldepth int, level level.Level, format string, args []interface{}) {
	fs := make(field.Fields, 5)

	var err error
	a := []interface{}{}
	ctxs := []context.Context{ctx}
	for _, f := range args {
		switch f := f.(type) {
		case field.Field:
			fs.Set(f)
		case context.Context:
			ctxs = append(ctxs, f)
		case error:
			a = append(a, f)
			err = f
		default:
			a = append(a, f)
		}
		// if fd, ok := f.(field.Field); ok {
		// 	fs.Set(fd)
		// } else {
		// 	a = append(a, f)
		// 	if ex, ok := f.(error); ok {
		// 		err = ex
		// 	}
		// }
	}

	if info, ok := runtime.GetCallLastFromError(err); ok {
		fs.Set(ErrorMethod(info.Method))
	}

	if stack, ok := runtime.GetCallStackFromError(err); ok {
		arr := []string{}
		for _, it := range stack {
			caller := it.Path
			if i := strings.LastIndex(caller, "/"); i > 0 {
				caller = caller[i+1:]
			}
			arr = append(arr, caller+"/"+it.File+":"+cast.ToString(it.Line))
		}
		fs.Set(ErrorMethod(strings.Join(arr, ",")))
	}

	if calldepth > -1 {
		info := runtime.Caller(calldepth + 1)
		fs.Set(TracebackCaller(getSplitLast(info.Method, "/")))
		fs.Set(TracebackLineNo(info.Line))
		// fmt.Printf("info.Path: %v\n", info.Path)
		p := info.Path
		if i := strings.LastIndex(p, "/"); i > 0 {
			p = p[i+1:]
		}
		if len(p) != 0 {
			p += "/"
		}

		fs.Set(TracebackPath(p + info.File))
	}

	if _, ok := fs[TraceIDKey]; !ok {
		for _, ctx := range ctxs {
			tid, sid := getTraceID(ctx)
			if len(tid) > 0 {
				fs.Set(TraceID(tid))
				if _, ok := fs[SpanIDKey]; !ok {
					if len(sid) > 0 {
						fs.Set(SpanID(sid))
					}
				}
				break
			}
		}
	}

	fs.Set(Message(format, a...))
	fs.Set(Level(level))
	if _, ok := fs[TimeKey]; !ok {
		fs.Set(Time(time.Now()))
	}

	// TODO 异步
	// entry.queue <- Entry(fs)
	entry.handlers.handle(Entry(fs))
}

func (entry *DefaultLogger) Log(ctx context.Context, calldepth int, level level.Level, args []interface{}) {
	entry.Logf(ctx, calldepth+1, level, "", args)
}

type logEntry struct {
	Level   level.Level
	Time    time.Time
	Message string
	Caller  runtime.CallerInfo
	Fields  field.Fields
}

// func (entry *DefaultLogger) Write(lvl level.Level, ts time.Time, msg string, caller runtime.CallerInfo) {
// 	ent := logEntry{}
// 	ent.Caller = caller
// 	ent.Level = lvl
// 	ent.Message = msg
// 	ent.Time = ts

// 	for _, v := range entry.handlers {
// 		v.Handle(nil, nil)
// 	}
// }

// ==================== EntryInterface ====================
type entryLog struct {
	calldepth int // 1
	logger    LoggerInterface
	ctx       context.Context
}

func (log *entryLog) Log(level int, a []interface{}) {
	log.logger.Log(log.ctx, log.calldepth+1, level, a)
}
func (log *entryLog) Logf(level int, format string, a []interface{}) {
	log.logger.Logf(log.ctx, log.calldepth+1, level, format, a)
}

func (log *entryLog) Debug(a ...interface{}) { log.Log(level.Debug, a) }
func (log *entryLog) Info(a ...interface{})  { log.Log(level.Info, a) }
func (log *entryLog) Warn(a ...interface{})  { log.Log(level.Warn, a) }
func (log *entryLog) Error(a ...interface{}) { log.Log(level.Error, a) }
func (log *entryLog) Fatal(a ...interface{}) { log.Log(level.Fatal, a) }

func (log *entryLog) Debugf(format string, a ...interface{}) { log.Logf(level.Debug, format, a) }
func (log *entryLog) Infof(format string, a ...interface{})  { log.Logf(level.Info, format, a) }
func (log *entryLog) Warnf(format string, a ...interface{})  { log.Logf(level.Warn, format, a) }
func (log *entryLog) Errorf(format string, a ...interface{}) { log.Logf(level.Error, format, a) }
func (log *entryLog) Fatalf(format string, a ...interface{}) { log.Logf(level.Fatal, format, a) }
