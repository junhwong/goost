package zap

import (
	"fmt"
	"math"
	"time"

	"github.com/junhwong/goost/apm"
	"github.com/junhwong/goost/apm/field"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func f() {

	l, err := zap.NewDevelopment()
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return
	}

	zap.Bools("", nil)
	l.Debug("hello")
	err = l.Sync()
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return
	}
	zap.CombineWriteSyncers()
}

func New(adapter apm.Adapter) *zap.Logger {
	log := zap.New(&ioCore{
		log:    adapter,
		fields: []zapcore.Field{zap.String("log_adapter", "zap")},
	}, zap.AddCaller())
	return log
}

// Core
type ioCore struct {
	log    apm.Adapter
	fields []zapcore.Field
	level  zapcore.Level
}

func (c *ioCore) Enabled(l zapcore.Level) bool { return true }
func (c *ioCore) With(fields []zapcore.Field) zapcore.Core {
	cc := &ioCore{
		log:    c.log,
		fields: make([]zapcore.Field, len(c.fields)),
		level:  c.level,
	}
	copy(cc.fields, c.fields)
	cc.fields = append(cc.fields, fields...)
	return cc
}

func (c *ioCore) Check(ent zapcore.Entry, ce *zapcore.CheckedEntry) *zapcore.CheckedEntry {
	if c.Enabled(ent.Level) {
		return ce.AddCore(ent, c)
	}
	return ce
}

type Field struct {
	Type  any
	Float float64
	Int   int16
	Uint  uint64
	Any   any
}

type FloatType interface {
	~float32 | ~float64
}
type IntType interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64
}
type UintType interface {
	~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64
}

type ComplexType interface {
	~complex64 | ~complex128
}

type BoolType interface {
	~bool
}

func Float[T FloatType](typeValue T) func(T) {
	return func(t T) {
		x := float64(t)
		fmt.Printf("x: %v\n", x)
	}
}

func Slice[T any](typeValue T) func(T) {
	return func(t T) {}
}
func Complex[T ComplexType](typeValue T) func(T) {
	return func(t T) {}
}

type Entry struct {
	ent    zapcore.Entry
	fields []zapcore.Field
	c      *ioCore
}

func (e Entry) GetLevel() (lvl apm.LogLevel) {
	apmLvl := apm.Debug
	switch e.ent.Level {
	case zapcore.InfoLevel:
		apmLvl = apm.Info
	case zapcore.WarnLevel:
		apmLvl = apm.Warn
	case zapcore.ErrorLevel:
		apmLvl = apm.Error
	case zapcore.DPanicLevel:
		apmLvl = apm.Warn
	case zapcore.PanicLevel:
		apmLvl = apm.Error
	case zapcore.FatalLevel:
		apmLvl = apm.Fatal
	}
	return apmLvl
}

func transform(fields []zapcore.Field, fs field.Fields) {
	for _, f := range fields {
		// if f.Key == "start time" || f.Key == "time spent" || f.Key == "response type" {
		// 	fmt.Printf("skip: %v\n", f)
		// 	continue
		// }
		if f.Key == "caller" {
			fmt.Printf("f: %v\n", f)
			panic("")
		}
		switch f.Type {
		case zapcore.ArrayMarshalerType:
			switch v := f.Interface.(type) {
			case []bool:
				fmt.Printf("v: %v\n", v)
			}
		case zapcore.ObjectMarshalerType:
		case zapcore.BinaryType:
			fs.Set(field.Dynamic(f.Key, string(f.Interface.([]byte))))
		case zapcore.BoolType:
			fs.Set(field.Dynamic(f.Key, f.Integer == 1))
		case zapcore.ByteStringType:
			fs.Set(field.Dynamic(f.Key, string(f.Interface.([]byte))))
		case zapcore.Complex128Type:
			fs.Set(field.Dynamic(f.Key, f.Interface))
		case zapcore.Complex64Type:
			fs.Set(field.Dynamic(f.Key, f.Interface))
		case zapcore.DurationType:
			fs.Set(field.Dynamic(f.Key, time.Duration(f.Integer)))
		case zapcore.Float64Type:
			fs.Set(field.Dynamic(f.Key, math.Float64frombits(uint64(f.Integer))))
		case zapcore.Float32Type:
			fs.Set(field.Dynamic(f.Key, math.Float32frombits(uint32(f.Integer))))
		case zapcore.Int64Type:
			fs.Set(field.Dynamic(f.Key, f.Integer))
		case zapcore.Int32Type:
			fs.Set(field.Dynamic(f.Key, f.Integer))
		case zapcore.Int16Type:
			fs.Set(field.Dynamic(f.Key, f.Integer))
		case zapcore.Int8Type:
			fs.Set(field.Dynamic(f.Key, f.Integer))
		case zapcore.StringType:
			fs.Set(field.Dynamic(f.Key, f.String))
		case zapcore.TimeType:
			t := time.UnixMicro(f.Integer / 1e3)
			t = t.Local().In(f.Interface.(*time.Location))
			fs.Set(field.Dynamic(f.Key, t))
		case zapcore.TimeFullType:
		case zapcore.Uint64Type:
			fs.Set(field.Dynamic(f.Key, uint64(f.Integer)))
		case zapcore.Uint32Type:
			fs.Set(field.Dynamic(f.Key, uint64(f.Integer)))
		case zapcore.Uint16Type:
			fs.Set(field.Dynamic(f.Key, uint64(f.Integer)))
		case zapcore.Uint8Type:
			fs.Set(field.Dynamic(f.Key, uint64(f.Integer)))
		case zapcore.UintptrType:
			fs.Set(field.Dynamic(f.Key, uint64(f.Integer)))
		case zapcore.ReflectType:
		case zapcore.NamespaceType:
		case zapcore.StringerType:
			fs.Set(field.Dynamic(f.Key, f.Interface.(fmt.Stringer).String()))
		case zapcore.ErrorType:
			fs.Set(field.Dynamic(f.Key, f.Interface))
		case zapcore.SkipType:
		case zapcore.InlineMarshalerType:
		}
	}
}

func (e Entry) GetFields() apm.Fields {
	fs := apm.Fields{}
	fs.Set(apm.Message(e.ent.Message))
	fs.Set(apm.Time(e.ent.Time))
	info := apm.CallerInfo{
		File:   e.ent.Caller.File,
		Line:   e.ent.Caller.Line,
		Method: e.ent.Caller.Function,
	}
	caller := info.Caller()
	// fmt.Printf("e.ent.Caller: %+v\n", e.ent.Caller.PC)
	// fmt.Printf("e.ent.Caller.File: %v\n", e.ent.Caller.File)
	// fmt.Printf("====caller: %v\n", caller)
	fs.Set(apm.TracebackCaller(caller))
	transform(e.c.fields, fs)
	transform(e.fields, fs)

	return fs
}

func (c *ioCore) Write(ent zapcore.Entry, fields []zapcore.Field) error {

	c.log.Dispatch(Entry{
		c:      c,
		ent:    ent,
		fields: fields,
	})
	// buf, err := c.enc.EncodeEntry(ent, fields)
	// if err != nil {
	// 	return err
	// }
	// _, err = c.out.Write(buf.Bytes())
	// buf.Free()
	// if err != nil {
	// 	return err
	// }
	// if ent.Level > zapcore.ErrorLevel {
	// 	// Since we may be crashing the program, sync the output. Ignore Sync
	// 	// errors, pending a clean solution to issue #370.
	// 	c.Sync()
	// }
	return nil
}

func (c *ioCore) Sync() error {
	// return c.out.Sync()
	return nil
}

// func (c *ioCore) clone() *ioCore {
// 	return &ioCore{
// 		LevelEnabler: c.LevelEnabler,
// 		enc:          c.enc.Clone(),
// 		out:          c.out,
// 	}
// }