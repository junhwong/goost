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
	calldepth int // 1
	fields    Fields
	ctx       context.Context
}

func (l *logImpl) SetCalldepth(a int) { l.calldepth = a }
func (l *logImpl) WithFields(fs ...Field) Interface {
	cl := l.clone()
	cl.fields = append(cl.fields, fs...)
	return cl
}
func (l *logImpl) clone() *logImpl {
	fieldsCopy := make([]field.Field, len(l.fields))
	copy(fieldsCopy, l.fields)

	return &logImpl{
		calldepth: l.calldepth,
		fields:    fieldsCopy,
	}
}

var callerContextKey = CallerInfo{}

func WithCaller(ctx context.Context, depth ...int) context.Context {
	d := 2
	if len(depth) > 0 {
		d = depth[len(depth)-1]
	}
	info := Caller(d)
	return context.WithValue(ctx, callerContextKey, info)
}
func CallerFromContext(ctx context.Context) (CallerInfo, bool) {
	obj := ctx.Value(callerContextKey)
	info, ok := obj.(CallerInfo)
	return info, ok
}
func (l logImpl) Log(level Level, args []interface{}) {
	l.calldepth++
	entry := &FieldsEntry{
		Level: level,
	}
	entry.Labels = append(entry.Labels, l.fields...)
	l.LogFS(entry, args)
}
func (l logImpl) LogFS(entry *FieldsEntry, args []interface{}, fs ...Field) {
	// var err error
	a := []interface{}{}
	ctxs := []context.Context{}
	var serr *StacktraceError
	for _, f := range args {
		switch f := f.(type) {
		case Field:
			entry.Labels = append(entry.Labels, f)
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

	if !entry.CallerInfo.Ok {
		for _, ctx := range ctxs {
			info, ok := CallerFromContext(ctx)
			if !ok {
				continue
			}
			entry.CallerInfo = info
			break
		}
	}

	if !entry.CallerInfo.Ok && l.calldepth > -1 {
		doCaller(l.calldepth+1, &entry.CallerInfo)
	}

	if serr != nil { // todo 额外处理
		caller := entry.CallerInfo.Caller()
		if ls := entry.Lookup(ErrorMethodKey.Name()); len(ls) == 0 {
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
				entry.Labels = append(entry.Labels, ErrorMethod(strings.Join(arr, ",")))
			}
		}

		if ls := entry.Lookup(ErrorStackTraceKey.Name()); len(ls) == 0 {
			entry.Labels = append(entry.Labels, ErrorStackTrace("%s", serr.Stack))
		}
	}

	if ls := entry.Lookup(TraceIDKey.Name()); len(ls) == 0 {
		for _, ctx := range ctxs {
			tid, sid := GetTraceID(ctx)
			if len(tid) > 0 {
				entry.Labels = append(entry.Labels, TraceID(tid))
				entry.Labels = append(entry.Labels, SpanID(sid))
				break
			}
		}
	}

	if len(a) > 0 {
		entry.Labels = append(entry.Labels, Message("", a...))
	}

	if entry.Time.IsZero() {
		entry.Time = time.Now()
	}
	dispatcher.Dispatch(entry)
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
