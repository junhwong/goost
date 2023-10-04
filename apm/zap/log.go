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

func New(adapter apm.Adapter, fs ...*field.Field) *zap.Logger {
	core := &ioCore{
		log: adapter,
		// fields: []zapcore.Field{zap.String(apm.LogAdapterKey.Name(), "zap")},
	}
	core.fs.Set(apm.LogAdapter("zap"))
	for _, v := range fs {
		if v != nil {
			core.fs.Set(v)
		}
	}

	log := zap.New(core, zap.AddCaller())

	return log
}

// Core
type ioCore struct {
	log   apm.Adapter
	level zapcore.Level
	fs    field.FieldSet
}

func (c *ioCore) Enabled(l zapcore.Level) bool { return true }
func (c *ioCore) With(fields []zapcore.Field) zapcore.Core {
	cc := &ioCore{
		log:   c.log,
		level: c.level,
		fs:    c.fs,
	}
	for _, v := range transform(fields) {
		if v != nil {
			cc.fs.Set(v)
		}
	}
	return cc
}

func (c *ioCore) Check(ent zapcore.Entry, ce *zapcore.CheckedEntry) *zapcore.CheckedEntry {
	if c.Enabled(ent.Level) {
		return ce.AddCore(ent, c)
	}
	return ce
}

func (c *ioCore) Write(ent zapcore.Entry, fields []zapcore.Field) error {
	info := &apm.CallerInfo{
		File:   ent.Caller.File,
		Line:   ent.Caller.Line,
		Method: ent.Caller.Function,
		Ok:     true,
	}
	var fs field.FieldSet
	for _, v := range c.fs {
		if v != nil {
			fs.Set(v)
		}
	}
	for _, v := range transform(fields) {
		if v != nil {
			fs.Set(v)
		}
	}

	apmLvl := field.LevelDebug
	switch ent.Level {
	case zapcore.InfoLevel:
		apmLvl = field.LevelInfo
	case zapcore.WarnLevel:
		apmLvl = field.LevelWarn
	case zapcore.ErrorLevel:
		apmLvl = field.LevelError
	case zapcore.DPanicLevel:
		apmLvl = field.LevelWarn
	case zapcore.PanicLevel:
		apmLvl = field.LevelError
	case zapcore.FatalLevel:
		apmLvl = field.LevelFatal
	}

	c.log.Dispatch(&apm.FieldsEntry{
		Time:       ent.Time,
		Level:      apmLvl,
		Fields:     fs,
		CallerInfo: info,
	})
	// buf, err := c.enc.EncodeEntry(ent, fields)

	return nil
}

func (c *ioCore) Sync() error {
	// return c.out.Sync()
	return nil
}

func transform(fields []zapcore.Field) (fs field.FieldSet) {

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
			fs = append(fs, field.Any(f.Key, string(f.Interface.([]byte))))
		case zapcore.BoolType:
			fs = append(fs, field.Any(f.Key, f.Integer == 1))
		case zapcore.ByteStringType:
			fs = append(fs, field.Any(f.Key, string(f.Interface.([]byte))))
		case zapcore.Complex128Type:
			fs = append(fs, field.Any(f.Key, f.Interface))
		case zapcore.Complex64Type:
			fs = append(fs, field.Any(f.Key, f.Interface))
		case zapcore.DurationType:
			fs = append(fs, field.Any(f.Key, time.Duration(f.Integer)))
		case zapcore.Float64Type:
			fs = append(fs, field.Any(f.Key, math.Float64frombits(uint64(f.Integer))))
		case zapcore.Float32Type:
			fs = append(fs, field.Any(f.Key, math.Float32frombits(uint32(f.Integer))))
		case zapcore.Int64Type:
			fs = append(fs, field.Any(f.Key, f.Integer))
		case zapcore.Int32Type:
			fs = append(fs, field.Any(f.Key, f.Integer))
		case zapcore.Int16Type:
			fs = append(fs, field.Any(f.Key, f.Integer))
		case zapcore.Int8Type:
			fs = append(fs, field.Any(f.Key, f.Integer))
		case zapcore.StringType:
			fs = append(fs, field.Any(f.Key, f.String))
		case zapcore.TimeType:
			t := time.UnixMicro(f.Integer / 1e3)
			t = t.Local().In(f.Interface.(*time.Location))
			fs = append(fs, field.Any(f.Key, t))
		case zapcore.TimeFullType:
		case zapcore.Uint64Type:
			fs = append(fs, field.Any(f.Key, uint64(f.Integer)))
		case zapcore.Uint32Type:
			fs = append(fs, field.Any(f.Key, uint64(f.Integer)))
		case zapcore.Uint16Type:
			fs = append(fs, field.Any(f.Key, uint64(f.Integer)))
		case zapcore.Uint8Type:
			fs = append(fs, field.Any(f.Key, uint64(f.Integer)))
		case zapcore.UintptrType:
			fs = append(fs, field.Any(f.Key, uint64(f.Integer)))
		case zapcore.ReflectType:
		case zapcore.NamespaceType:
		case zapcore.StringerType:
			fs = append(fs, field.Any(f.Key, f.Interface.(fmt.Stringer).String()))
		case zapcore.ErrorType:
			fs = append(fs, field.Any(f.Key, f.Interface))
		case zapcore.SkipType:
		case zapcore.InlineMarshalerType:
		}
	}
	return
}

// type Field struct {
// 	Type  any
// 	Float float64
// 	Int   int16
// 	Uint  uint64
// 	Any   any
// }

// type FloatType interface {
// 	~float32 | ~float64
// }
// type IntType interface {
// 	~int | ~int8 | ~int16 | ~int32 | ~int64
// }
// type UintType interface {
// 	~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64
// }

// type ComplexType interface {
// 	~complex64 | ~complex128
// }

// type BoolType interface {
// 	~bool
// }

// func Float[T FloatType](typeValue T) func(T) {
// 	return func(t T) {
// 		x := float64(t)
// 		fmt.Printf("x: %v\n", x)
// 	}
// }

// func Slice[T any](typeValue T) func(T) {
// 	return func(t T) {}
// }
// func Complex[T ComplexType](typeValue T) func(T) {
// 	return func(t T) {}
// }

// type Entry struct {
// 	ent    zapcore.Entry
// 	fields []zapcore.Field
// 	c      *ioCore
// }

// func (e Entry) GetLevel() (lvl apm.Level) {
// 	apmLvl := field.LevelDebug
// 	switch e.ent.Level {
// 	case zapcore.InfoLevel:
// 		apmLvl = field.LevelInfo
// 	case zapcore.WarnLevel:
// 		apmLvl = field.LevelWarn
// 	case zapcore.ErrorLevel:
// 		apmLvl = field.LevelError
// 	case zapcore.DPanicLevel:
// 		apmLvl = field.LevelWarn
// 	case zapcore.PanicLevel:
// 		apmLvl = field.LevelError
// 	case zapcore.FatalLevel:
// 		apmLvl = field.LevelFatal
// 	}
// 	return apmLvl
// }
// func (e Entry) GetTime() (v time.Time) {
// 	return e.ent.Time
// }
// func (e Entry) GetMessage() (v string) {
// 	return e.ent.Message
// }

// func (e Entry) GetFields() field.FieldSet {

// 	info := apm.CallerInfo{
// 		File:   e.ent.Caller.File,
// 		Line:   e.ent.Caller.Line,
// 		Method: e.ent.Caller.Function,
// 	}
// 	caller := info.Caller()
// 	// fmt.Printf("e.ent.Caller: %+v\n", e.ent.Caller.PC)
// 	// fmt.Printf("e.ent.Caller.File: %v\n", e.ent.Caller.File)
// 	// fmt.Printf("====caller: %v\n", caller)
// 	fs := field.FieldSet{apm.TracebackCaller(caller)}

// 	fs = append(fs, transform(e.c.fields)...)
// 	fs = append(fs, transform(e.fields)...)

// 	return fs
// }

// func (c *ioCore) clone() *ioCore {
// 	return &ioCore{
// 		LevelEnabler: c.LevelEnabler,
// 		enc:          c.enc.Clone(),
// 		out:          c.out,
// 	}
// }
