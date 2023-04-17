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

func (l *FieldsEntry) SetCalldepth(v int) { l.calldepth = v }
func (l *FieldsEntry) WithFields(fs ...*field.Field) Interface {
	cl := l.clone()
	for _, f := range fs {
		cl.Fields.Set(f)
	}
	return cl
}
func (l *FieldsEntry) clone() *FieldsEntry {
	l.mu.Lock()
	r := &FieldsEntry{
		calldepth: l.calldepth,
	}
	for _, f := range l.Fields {
		r.Fields.Set(field.Clone(f))
	}
	l.mu.Unlock()
	return r
}

func (l *FieldsEntry) Log(level Level, args []interface{}) {
	l = l.clone()
	l.calldepth++
	l.Level = level

	l.do(args, func() {})
}

func (entry *FieldsEntry) do(args []interface{}, befor func()) {
	// var err error
	var (
		a    []interface{}
		ctxs []context.Context
		serr *StacktraceError
	)
	for _, f := range args {
		switch f := f.(type) {
		case field.Field:
			entry.Fields.Set(&f)
		case *field.Field:
			entry.Fields.Set(f)
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

	if !entry.CallerInfo.Ok && entry.calldepth > -1 {
		doCaller(entry.calldepth+1, &entry.CallerInfo)
	}

	if serr != nil { // todo 额外处理
		caller := entry.CallerInfo.Caller()
		if entry.GetFields().Get(ErrorMethodKey.Name()) == nil {
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
				entry.Fields.Set(ErrorMethod(strings.Join(arr, ",")))
			}
		}

		if entry.GetFields().Get(ErrorStackTraceKey.Name()) == nil {
			entry.Fields.Set(ErrorStackTrace("%s", serr.Stack))
		}
	}

	if entry.GetFields().Get(TraceIDKey.Name()) == nil {
		for _, ctx := range ctxs {
			tid, sid := GetTraceID(ctx)
			if len(tid) > 0 {
				entry.Fields.Set(TraceIDField(tid))
				entry.Fields.Set(SpanID(sid))
				break
			}
		}
	}

	if len(a) > 0 {
		entry.Fields.Set(Message("", a...))
	}
	befor()
	if entry.Time.IsZero() {
		entry.Time = time.Now()
	}
	dispatcher.Dispatch(entry)
}

func (l *FieldsEntry) Debug(a ...interface{}) { l.Log(field.LevelDebug, a) }
func (l *FieldsEntry) Info(a ...interface{})  { l.Log(field.LevelInfo, a) }
func (l *FieldsEntry) Warn(a ...interface{})  { l.Log(field.LevelWarn, a) }
func (l *FieldsEntry) Error(a ...interface{}) { l.Log(field.LevelError, a) }
func (l *FieldsEntry) Fatal(a ...interface{}) { l.Log(field.LevelFatal, a) }
