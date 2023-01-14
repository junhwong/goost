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
}
type FormatLogger interface {
	Debugf(string, ...interface{})
	Infof(string, ...interface{})
	Warnf(string, ...interface{})
	Errorf(string, ...interface{})
	Fatalf(string, ...interface{})
}

// ==================== EntryInterface ====================
type logImpl struct {
	calldepth  int // 1
	fields     Fields
	dispatcher Dispatcher
	ctx        context.Context
}

func (l *logImpl) SetCalldepth(a int) { l.calldepth = a }
func (l *logImpl) WithFields(fs ...Field) Interface {
	cl := l.clone()
	cl.fields = append(cl.fields, fs...)
	return l
}
func (l *logImpl) clone() *logImpl {
	fieldsCopy := make([]field.Field, len(l.fields))
	copy(fieldsCopy, l.fields)

	return &logImpl{
		calldepth:  l.calldepth,
		dispatcher: l.dispatcher,
		fields:     fieldsCopy,
	}
}

func (l logImpl) Log(level Level, args []interface{}) {
	l.calldepth++
	l.LogFS(args, LevelField(int(level)))
}

func (l logImpl) LogFS(args []interface{}, fs ...Field) {
	entry := make(field.Fields, 5)
	entry.Set(l.fields...)
	entry.Set(fs...)
	// var err error
	a := []interface{}{}
	ctxs := []context.Context{}
	var serr *StacktraceError
	for _, f := range args {
		switch f := f.(type) {
		case Field:
			entry.Set(f)
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
	calldepth := l.calldepth
	if calldepth > -1 {
		info := Caller(calldepth + 1)
		caller = info.Caller()
		entry.Set(TracebackCaller(caller))
	}
	if serr != nil {
		if _, ok := entry[ErrorMethodKey]; !ok {
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
				entry.Set(ErrorMethod(strings.Join(arr, ",")))
			}
		}
		if _, ok := entry[ErrorStackTraceKey]; !ok {
			entry.Set(ErrorStackTrace("%s", serr.Stack))
		}
	}

	if _, ok := entry[TraceIDKey]; !ok {
		for _, ctx := range ctxs {
			tid, sid := GetTraceID(ctx)
			if len(tid) > 0 {
				entry.Set(TraceID(tid))
				if _, ok := entry[SpanIDKey]; !ok {
					if len(sid) > 0 {
						entry.Set(SpanID(sid))
					}
				}
				break
			}
		}
	}

	if len(a) > 0 {
		entry.Set(Message("", a...))
	}

	if _, ok := entry[TimeKey]; !ok {
		entry.Set(Time(time.Now()))
	}

	l.dispatcher.Dispatch(FieldsEntry(entry))
}

func (l *logImpl) Debug(a ...interface{}) { l.Log(LevelDebug, a) }
func (l *logImpl) Info(a ...interface{})  { l.Log(LevelInfo, a) }
func (l *logImpl) Warn(a ...interface{})  { l.Log(LevelWarn, a) }
func (l *logImpl) Error(a ...interface{}) { l.Log(LevelError, a) }
func (l *logImpl) Fatal(a ...interface{}) { l.Log(LevelFatal, a) }

// func (log *logImpl) Logf(level LogLevel, format string, a []interface{}) {
// 	log.logger.Logf(log.ctx, log.calldepth+1, level, format, a)
// }
// func (log *entryLog) Debugf(format string, a ...interface{}) { log.Logf(Debug, format, a) }
// func (log *entryLog) Infof(format string, a ...interface{})  { log.Logf(Info, format, a) }
// func (log *entryLog) Warnf(format string, a ...interface{})  { log.Logf(Warn, format, a) }
// func (log *entryLog) Errorf(format string, a ...interface{}) { log.Logf(Error, format, a) }
// func (log *entryLog) Fatalf(format string, a ...interface{}) { log.Logf(Fatal, format, a) }
