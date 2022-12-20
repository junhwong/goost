package apm

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/junhwong/goost/apm/field"
)

// LoggerInterface 日志记录接口
type LoggerInterface interface {
	Close() error
	Flush() error

	Log(Entry)
}

// Logger 日志操作接口
type Logger interface {
	Debug(...interface{})
	Info(...interface{})
	Warn(...interface{})
	Error(...interface{})
	Fatal(...interface{})

	// // Deprecated
	// Debugf(string, ...interface{})
	// // Deprecated
	// Infof(string, ...interface{})
	// // Deprecated
	// Warnf(string, ...interface{})
	// // Deprecated: use Error
	// Errorf(string, ...interface{})
	// // Deprecated:
	// Fatalf(string, ...interface{})
}

type loggerImpl struct {
	logImpl
	fields []field.Field
}

// ==================== EntryInterface ====================
type logImpl struct {
	calldepth  int // 1
	dispatcher Dispatcher
	ctx        context.Context
}

func (log *logImpl) Log(level LogLevel, args []interface{}) {
	calldepth := log.calldepth
	fs := make(field.Fields, 5)

	// var err error
	a := []interface{}{}
	ctxs := []context.Context{}
	var serr *StacktraceError
	for _, f := range args {
		switch f := f.(type) {
		case Field:
			fs.Set(f)
		case context.Context:
			ctxs = append(ctxs, f)
		case error:
			a = append(a, f)
			if serr == nil {
				errors.As(f, &serr)
			}
		default:
			a = append(a, f)
		}
	}
	caller := ""
	if calldepth > -1 {
		info := Caller(calldepth + 1)
		caller = info.Caller()
		fs.Set(TracebackCaller(caller))
	}
	if serr != nil {
		if _, ok := fs[ErrorMethodKey]; !ok {
			stack := StackToCallerInfo(serr.Stack)
			arr := []string{}
			for _, it := range stack {
				arr = append(arr, it.Caller())
			}
			if len(arr) > 0 && arr[0] == caller {
				arr = arr[1:]
			}
			if n := len(arr); n > 0 && arr[n-1] == caller {
				arr = arr[:n-1]
			}
			// fmt.Printf("stack: %s\n", serr.Stack)
			// fmt.Printf("arr: %v\n", arr)
			if len(arr) > 0 {
				fs.Set(ErrorMethod(strings.Join(arr, ",")))
			}
		}
		if _, ok := fs[ErrorStackTraceKey]; !ok {
			fs.Set(ErrorStackTrace("%s", serr.Stack))
		}
	}

	if _, ok := fs[TraceIDKey]; !ok {
		for _, ctx := range ctxs {
			tid, sid := GetTraceID(ctx)
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

	fs.Set(Message("", a...))
	fs.Set(Level(int(level)))
	if _, ok := fs[TimeKey]; !ok {
		fs.Set(Time(time.Now()))
	}

	log.dispatcher.Dispatch(FieldsEntry(fs))
}

// func (log *logImpl) Logf(level LogLevel, format string, a []interface{}) {
// 	log.logger.Logf(log.ctx, log.calldepth+1, level, format, a)
// }

func (log *logImpl) Debug(a ...interface{}) { log.Log(Debug, a) }
func (log *logImpl) Info(a ...interface{})  { log.Log(Info, a) }
func (log *logImpl) Warn(a ...interface{})  { log.Log(Warn, a) }
func (log *logImpl) Error(a ...interface{}) { log.Log(Error, a) }
func (log *logImpl) Fatal(a ...interface{}) { log.Log(Fatal, a) }

// func (log *entryLog) Debugf(format string, a ...interface{}) { log.Logf(Debug, format, a) }
// func (log *entryLog) Infof(format string, a ...interface{})  { log.Logf(Info, format, a) }
// func (log *entryLog) Warnf(format string, a ...interface{})  { log.Logf(Warn, format, a) }
// func (log *entryLog) Errorf(format string, a ...interface{}) { log.Logf(Error, format, a) }
// func (log *entryLog) Fatalf(format string, a ...interface{}) { log.Logf(Fatal, format, a) }
