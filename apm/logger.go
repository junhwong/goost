package apm

import (
	"context"
	"sync"
	"time"

	"github.com/junhwong/goost/apm/level"
	"github.com/junhwong/goost/pkg/field"
	"github.com/junhwong/goost/runtime"
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

	Debugf(string, ...interface{})
	Infof(string, ...interface{})
	Warnf(string, ...interface{})
	Errorf(string, ...interface{})
	Fatalf(string, ...interface{})
}

type DefaultLogger struct {
	mu       sync.Mutex
	wg       sync.WaitGroup
	queue    chan Entry
	cancel   context.CancelFunc
	handlers handlerSlice
}

func (logger *DefaultLogger) handle(entry Entry) {
	std.mu.Lock()
	defer std.mu.Unlock()

	size := logger.handlers.Len()
	crt := 0
	var next func()
	next = func() {
		if crt >= size {
			return
		}
		h := logger.handlers[crt]
		crt++
		h.Handle(entry, next)
	}
	next()
}

func (logger *DefaultLogger) flush() {
	for {
		select {
		case entry := <-logger.queue:
			logger.handle(entry)
		default:
			return
		}
	}
}

func (logger *DefaultLogger) Run(stopCh <-chan struct{}) {
	for {
		select {
		case entry := <-logger.queue:
			logger.handle(entry)
		case <-stopCh:
			goto END
		}
	}

END:
	logger.flush()
}

func (logger *DefaultLogger) Close() error {
	logger.cancel()
	time.Sleep(time.Millisecond) // 给协程一点时间启动
	logger.wg.Wait()
	return nil
}

func (logger *DefaultLogger) NewSpan(ctx context.Context, options ...SpanOption) (context.Context, Span) {
	return newSpan(ctx, logger, options)
}

type _GetCallLastInfo interface {
	GetCallLastInfo() runtime.CallerInfo
}

func (entry *DefaultLogger) Logf(ctx context.Context, calldepth int, level level.Level, format string, args []interface{}) {
	fs := make(field.Fields, 5)

	var err error
	a := []interface{}{}
	for _, f := range args {
		if fd, ok := f.(field.Field); ok {
			fs.Set(fd)
		} else {
			a = append(a, f)
			if ex, ok := f.(error); ok {
				err = ex
			}
		}
	}

	if info, ok := runtime.GetCallLastFromError(err); ok {
		fs.Set(ErrorMethod(info.Method))
	}

	if calldepth > -1 {
		info := runtime.Caller(calldepth + 1)
		fs.Set(TracebackCaller(getSplitLast(info.Method, "/")))
		fs.Set(TracebackLineNo(info.Line))
		fs.Set(TracebackPath(info.File))
	}

	if _, ok := fs[TraceIDKey]; !ok && ctx != nil {
		fs.Set(TraceID(getTraceID(ctx)))
	}

	fs.Set(Message(format, a...))
	fs.Set(Level(level))
	if _, ok := fs[TimeKey]; !ok {
		fs.Set(Time(time.Now()))
	}

	entry.queue <- Entry(fs)
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

func (entry *DefaultLogger) Write(lvl level.Level, ts time.Time, msg string, caller runtime.CallerInfo) {
	ent := logEntry{}
	ent.Caller = caller
	ent.Level = lvl
	ent.Message = msg
	ent.Time = ts

	// io.Discard.Write([]byte(fmt.Sprintf("ent: %v", ent)))

	for _, v := range entry.handlers {
		v.Handle(nil, nil)
	}

	// ent := make(field.Fields, 5)
	// for _, it := range fs {
	// 	ent.Set(it)
	// }
	// if _, ok := ent[LevelKey]; !ok {
	// 	return
	// }
	// // fmt.Println("", ent)
	// entry.queue <- ent
}

// ==================== EntryInterface ====================
type entryLog struct {
	calldepth int // 1
	logger    LoggerInterface
	ctx       context.Context
}

func (log *entryLog) Debug(a ...interface{}) { log.logger.Log(log.ctx, log.calldepth, level.Debug, a) }
func (log *entryLog) Info(a ...interface{})  { log.logger.Log(log.ctx, log.calldepth, level.Info, a) }
func (log *entryLog) Warn(a ...interface{})  { log.logger.Log(log.ctx, log.calldepth, level.Warn, a) }
func (log *entryLog) Error(a ...interface{}) { log.logger.Log(log.ctx, log.calldepth, level.Error, a) }
func (log *entryLog) Fatal(a ...interface{}) { log.logger.Log(log.ctx, log.calldepth, level.Fatal, a) }

func (log *entryLog) Debugf(format string, a ...interface{}) {
	log.logger.Logf(log.ctx, log.calldepth, level.Debug, format, a)
}
func (log *entryLog) Infof(format string, a ...interface{}) {
	log.logger.Logf(log.ctx, log.calldepth, level.Info, format, a)
}
func (log *entryLog) Warnf(format string, a ...interface{}) {
	log.logger.Logf(log.ctx, log.calldepth, level.Warn, format, a)
}
func (log *entryLog) Errorf(format string, a ...interface{}) {
	log.logger.Logf(log.ctx, log.calldepth, level.Error, format, a)
}
func (log *entryLog) Fatalf(format string, a ...interface{}) {
	log.logger.Logf(log.ctx, log.calldepth, level.Fatal, format, a)
}
