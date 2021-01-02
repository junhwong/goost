package apm

import (
	"bytes"
	"fmt"
	"sync"
	"time"

	"github.com/junhwong/goost/pkg/field"
)

var bufferPool *sync.Pool

func init() {
	bufferPool = &sync.Pool{
		New: func() interface{} {
			return new(bytes.Buffer)
		},
	}
}

type Fields = field.Fields
type Field = field.Field

// Entry 单条日志记录
type Entry struct {
	equeue  func(entry *Entry)
	Time    time.Time // 时间戳
	Level   Level     // 日志级别
	Message string    // 消息
	Data    Fields    // 附加数据
}

func NewEntry(equeue func(entry *Entry)) *Entry {
	return &Entry{
		equeue: equeue,
		Data:   make(Fields, 5),
	}
}

func (entry *Entry) log(lvl Level, args []interface{}) {
	if entry.equeue == nil {
		return
	}
	// defer entry.Logger.releaseEntry(entry)

	// if entry.Logger.level() > entry.Level && len(args) == 0 {
	// 	return
	// }
	entry.Level = lvl
	if entry.Time.IsZero() {
		entry.Time = time.Now()
	}
	format, fmtok := "", false
	if len(args) > 0 {
		format, fmtok = args[0].(string)
		if fmtok {
			if format == "" {
				fmtok = false
			}
			args = args[1:]
		}
	}

	a := []interface{}{}
	for _, f := range args {
		if fd, ok := f.(*Field); ok {
			entry.Data[fd.Key] = fd
		} else if f != nil {
			a = append(a, f)
		}
	}

	switch {
	case fmtok && len(a) > 0:
		entry.Message = fmt.Sprintf(format, a...)
	case len(a) > 0:
		entry.Message = fmt.Sprint(a...)
	default:
		entry.Message = format
	}

	entry.equeue(entry)

}

func (e *Entry) Debug(args ...interface{}) { e.log(DebugLevel, args) }
func (e *Entry) Info(args ...interface{})  { e.log(InfoLevel, args) }
func (e *Entry) Warn(args ...interface{})  { e.log(WarningLevel, args) }
func (e *Entry) Error(args ...interface{}) { e.log(ErrorLevel, args) }
func (e *Entry) Fatal(args ...interface{}) { e.log(FatalLevel, args) }
func (e *Entry) Trace(args ...interface{}) { e.log(TraceLevel, args) }
