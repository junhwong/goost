package apm

import (
	"context"
	"time"

	"github.com/junhwong/goost/apm/field"
	"github.com/junhwong/goost/apm/field/loglevel"
)

// LoggerInterface 日志记录接口
type LoggerInterface interface {
	Close() error
	Flush() error

	Log(*factoryEntry)
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
type WithOption interface {
	applyWithOption(*factoryEntry)
}

// 统一接口
type Interface interface {
	Logger
	SpanFactory
	// With(options ...WithOption) Interface
}

func (l *factoryEntry) With(options ...WithOption) Interface {
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

func (s *factoryEntry) SetAttributes(a ...*field.Field) {
	for _, f := range a {
		s.Set(f)
	}
}
func (l *factoryEntry) SetCalldepth(v int) { l.calldepth = v }
func (l *factoryEntry) GetCalldepth() int  { return l.calldepth }
func (l *factoryEntry) CalldepthInc() Interface {
	l.calldepth++
	return l
}
func (l *factoryEntry) WithFields(fs ...*field.Field) Interface {
	cl := l.new()
	for _, f := range fs {
		cl.Set(f)
	}
	return cl
}
func (l *factoryEntry) new() *factoryEntry {
	l.mu.Lock()
	r := &factoryEntry{
		calldepth: l.calldepth,
	}
	r.Field = field.Clone(l.Field)
	// r.Set(Time(time.Now()))
	l.mu.Unlock()
	return r
}

func (l *factoryEntry) Log(level loglevel.Level, args []interface{}) {
	if len(args) == 0 {
		return
	}

	var (
		entry *field.Field
		arr   []interface{}
		ctxs  []context.Context
	)

	for _, it := range args {
		switch it := it.(type) {
		case field.Field:
			if entry == nil {
				entry = field.MakeRoot()
			}
			entry.Set(&it)
		case *field.Field:
			if entry == nil {
				entry = field.MakeRoot()
			}
			entry.Set(it)
		case context.Context:
			ctxs = append(ctxs, it)
		default:
			if ctx, _ := it.(context.Context); ctx != nil {
				ctxs = append(ctxs, ctx)
			} else if len(arr) > 0 || it != nil {
				arr = append(arr, it)
			}
		}
	}

	if len(arr) == 0 && entry == nil {
		return
	}

	l.mu.Lock()
	if entry == nil {
		entry = field.Clone(l.Field)
	} else {
		for _, f := range l.Items {
			if field.GetLast(entry.Items, f.GetName()) == nil {
				entry.Set(f)
			}
		}
	}
	l.mu.Unlock()

	entry.Set(Time(time.Now()))

	do(level, entry, l.calldepth+1, ctxs, arr, func() {})
}

func do(level loglevel.Level, entry *field.Field, calldepth int, ctxs []context.Context, args []interface{}, befor func()) {
	info := entry.GetItem("source")
	if info == nil || !info.IsGroup() {
		ci := CallerInfo{}
		doCaller(calldepth, &ci)
		info := field.Make("source")
		info.SetKind(field.GroupKind, false, false)
		info.Set(field.Make("file").SetString(ci.File))
		info.Set(field.Make("line").SetInt(int64(ci.Line)))
		info.Set(field.Make("func").SetString(ci.Method))
		entry.Set(info)
	}

	if entry.GetItem(TraceIDKey.Name()) == nil {
		for _, ctx := range ctxs {
			p := SpanContextFrom(ctx)
			if p != nil {
				entry.Set(TraceIDField(p.GetTranceID()))
				if level == loglevel.Trace {
					entry.Set(SpanID(p.GetSpanID()))
				}
				break
			}
		}
	}
	entry.Set(Message("", args...))
	// if len(a) > 0 {
	// 	entry.Set(Message("", a...))
	// }
	entry.Set(LevelField(level))

	befor()

	Dispatch(entry)
}

func (l *factoryEntry) Debug(a ...interface{}) { l.Log(loglevel.Debug, a) }
func (l *factoryEntry) Info(a ...interface{})  { l.Log(loglevel.Info, a) }
func (l *factoryEntry) Warn(a ...interface{})  { l.Log(loglevel.Warn, a) }
func (l *factoryEntry) Error(a ...interface{}) { l.Log(loglevel.Error, a) }
func (l *factoryEntry) Fatal(a ...interface{}) { l.Log(loglevel.Fatal, a) }
