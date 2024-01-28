package apm

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/junhwong/goost/apm/field"
	"github.com/junhwong/goost/apm/field/loglevel"
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
func (l *FieldsEntry) With(options ...WithOption) Interface {
	if len(options) == 0 {
		return l
	}
	cl := l.new()
	for _, o := range options {
		if o != nil {
			o.applyWithOption(cl)
		}
	}
	return cl
}
func (s *FieldsEntry) SetAttributes(a ...*field.Field) {
	for _, f := range a {
		s.Set(f)
	}
}
func (l *FieldsEntry) SetCalldepth(v int) { l.calldepth = v }
func (l *FieldsEntry) GetCalldepth() int  { return l.calldepth }
func (l *FieldsEntry) CalldepthInc() Interface {
	l.calldepth++
	return l
}
func (l *FieldsEntry) WithFields(fs ...*field.Field) Interface {
	cl := l.new()
	for _, f := range fs {
		cl.Set(f)
	}
	return cl
}
func (l *FieldsEntry) new() *FieldsEntry {
	l.mu.Lock()
	r := &FieldsEntry{
		calldepth: l.calldepth,
	}
	r.Field = *field.Clone(&l.Field)
	r.Set(Time(time.Now()))
	l.mu.Unlock()
	return r
}

func (l *FieldsEntry) Log(level loglevel.Level, args []interface{}) {
	if len(args) == 0 {
		return
	}
	pass := false
	for _, v := range args {
		if v != nil {
			pass = true
			break
		}
	}
	if !pass {
		return
	}

	l = l.new()
	l.calldepth++
	l.Set(LevelField(level))

	l.do(args, func() {})
}

func (entry *FieldsEntry) do(args []interface{}, befor func()) {
	// var err error
	if entry.GetTime().IsZero() {
		panic("apm: entry.Time cannot be zero")
	}

	var (
		a    []interface{}
		ctxs []context.Context
		serr *StacktraceError
	)
	for _, f := range args {
		switch f := f.(type) {
		case field.Field:
			entry.Set(&f)
		case *field.Field:
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

	if entry.CallerInfo == nil {
		for _, ctx := range ctxs {
			info := CallerFrom(ctx)
			if info == nil {
				continue
			}
			entry.CallerInfo = info
			break
		}
	}

	if entry.CallerInfo == nil && entry.calldepth > -1 {
		entry.CallerInfo = &CallerInfo{}
		doCaller(entry.calldepth+1, entry.CallerInfo)
	}

	if serr != nil { // todo 额外处理
		caller := entry.CallerInfo.Caller()
		if entry.GetItem(ErrorMethodKey.Name()) == nil {
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

		if entry.GetItem(ErrorStackTraceKey.Name()) == nil {
			entry.Set(ErrorStackTrace("%s", serr.Stack))
		}
	}

	if entry.GetItem(TraceIDKey.Name()) == nil {
		for _, ctx := range ctxs {
			tid, sid := GetTraceID(ctx)
			if len(tid) > 0 {
				entry.Set(TraceIDField(tid))
				entry.Set(SpanID(sid))
				break
			}
		}
	}

	if len(a) > 0 {
		entry.Set(Message("", a...))
	}
	befor()

	if d := GetDispatcher(); d != nil {
		d.Dispatch(entry)
	} else {
		// todo
	}
}

func (l *FieldsEntry) Debug(a ...interface{}) { l.Log(loglevel.Debug, a) }
func (l *FieldsEntry) Info(a ...interface{})  { l.Log(loglevel.Info, a) }
func (l *FieldsEntry) Warn(a ...interface{})  { l.Log(loglevel.Warn, a) }
func (l *FieldsEntry) Error(a ...interface{}) { l.Log(loglevel.Error, a) }
func (l *FieldsEntry) Fatal(a ...interface{}) { l.Log(loglevel.Fatal, a) }
