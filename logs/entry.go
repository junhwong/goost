package logs

import (
	"bytes"
	"fmt"
	"os"
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
	Logger  *Logger   // 处理器
	Time    time.Time // 时间戳
	Level   Level     // 日志级别
	Message string    // 消息
	Data    Fields    // 附加数据
}

func NewEntry(logger *Logger) *Entry {
	return &Entry{
		Logger: logger,
		Data:   make(Fields, 5),
	}
}

func (entry Entry) log(lvl Level, args []interface{}) {
	if entry.Logger == nil {
		return
	}
	defer entry.Logger.releaseEntry(&entry)

	if entry.Logger.level() > entry.Level && len(args) == 0 {
		return
	}

	format, fmtok := args[0].(string)

	if fmtok {
		if format == "" {
			fmtok = false
		}
		args = args[1:]
	}
	entry.Level = lvl
	entry.Time = time.Now()

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

	buf := bufferPool.Get().(*bytes.Buffer)
	defer bufferPool.Put(buf)
	buf.Reset()
	err := entry.Logger.format(&entry, buf)
	if err != nil {
		entry.Logger.mu.Lock()
		fmt.Fprintf(os.Stderr, "Failed to obtain reader, %v\n", err)
		entry.Logger.mu.Unlock()
		return
	}

	err = entry.Logger.write(buf)
	if err != nil {
		entry.Logger.mu.Lock()
		fmt.Fprintf(os.Stderr, "Failed to write to log, %v\n", err)
		entry.Logger.mu.Unlock()
		return
	}
}

func (e *Entry) Debug(args ...interface{}) { e.log(DebugLevel, args) }
func (e *Entry) Info(args ...interface{})  { e.log(InfoLevel, args) }
func (e *Entry) Warn(args ...interface{})  { e.log(WarningLevel, args) }
func (e *Entry) Error(args ...interface{}) { e.log(ErrorLevel, args) }
func (e *Entry) Fatal(args ...interface{}) { e.log(FatalLevel, args) }
func (e *Entry) Trace(args ...interface{}) { e.log(TraceLevel, args) }
