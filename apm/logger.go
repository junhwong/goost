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
func (s *factoryEntry) setAttributes(a ...*field.Field) {
	for _, f := range a {
		s.Set(f)
	}
}
func (l *factoryEntry) setCalldepth(v int) { l.calldepth = v }
func (l *factoryEntry) getCalldepth() int  { return l.calldepth }

// func (l *factoryEntry) CalldepthInc() Interface {
// 	l.calldepth++
// 	return l
// }

//	func (l *factoryEntry) WithFields(fs ...*field.Field) Interface {
//		cl := l.new()
//		for _, f := range fs {
//			cl.Set(f)
//		}
//		return cl
//	}
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
	if l == nil || len(args) == 0 {
		return
	}

	var (
		entry *field.Field
		arr   []interface{}
		ctxs  []context.Context
	)

	for _, it := range args {
		if it == nil {
			continue
		}
		switch it := it.(type) {
		case field.Field:
			if it.IsNull() {
				continue
			}
			if entry == nil {
				entry = field.MakeRoot()
			}
			entry.Set(&it)
		case *field.Field:
			if it == nil || it.IsNull() {
				continue
			}
			if entry == nil {
				entry = field.MakeRoot()
			}
			entry.Set(it)
		default:
			if ctx, _ := it.(context.Context); ctx != nil {
				ctxs = append(ctxs, ctx)
			} else {
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
				entry.Set(field.Clone(f))
			}
		}
	}
	l.mu.Unlock()

	entry.Set(Time(time.Now()))

	do(level, entry, l.calldepth, ctxs, arr, func() {})
}

func MakeSource(file string, line int, fn string) *field.Field {
	f := field.Make("source")
	f.SetKind(field.GroupKind, false, false)
	f.Set(field.Make("file").SetString(file))
	f.Set(field.Make("line").SetInt(int64(line)))
	f.Set(field.Make("func").SetString(fn))
	return f
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
				// if level == loglevel.Trace {
				// 	entry.Set(SpanID(p.GetSpanID()))
				// }
				// break
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
