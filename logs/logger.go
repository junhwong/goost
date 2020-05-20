package logs

import (
	"bytes"
	"context"
	"io"
	"sync"

	"github.com/junhwong/goost/pkg/field/common"
)

// ILogger 纯日志接口
type ILogger interface {
	Debug(...interface{})
	Info(...interface{})
	Warn(...interface{})
	Error(...interface{})
	Fatal(...interface{})
	Trace(...interface{})
}

// type EntryHandler interface {
// 	Handle(*Entry) error
// 	IsEnabled(Level) bool
// }

// type LoggerHandler interface {
// 	EntryHandler
// }

// type ForkedLogger interface {
// 	Logger
// 	IsEnabled(Level) bool
// 	WithPrefix(prefix string) ForkedLogger
// 	WithContext(prefix string) ForkedLogger
// }

// type RootLogger interface {
// 	ForkedLogger
// 	SetHandler(LoggerHandler) error
// 	SetLevel(Level)
// }

// type DefaultLogger struct {
// 	Prefix string
// }

// func (log *DefaultLogger) newEntry() *Entry {
// 	return &Entry{
// 		Prefix: log.Prefix,
// 		// Handler: log,
// 		Tags:   make(Fields),
// 		Data:   make(Fields),
// 		Fields: make([]*Field, 0),
// 	}
// }

// func (log *DefaultLogger) Handle(entry *Entry) error {
// 	formatT(entry)
// 	return nil
// }

// func (log *DefaultLogger) WithPrefix(prefix string) ForkedLogger {
// 	return log
// }
// func (log *DefaultLogger) WithContext(prefix string) ForkedLogger {
// 	return log
// }
// func (log *DefaultLogger) SetHandler(lhd LoggerHandler) error {
// 	return nil
// }
// func (log *DefaultLogger) SetLevel(lhd Level) {
// }
// func (log *DefaultLogger) IsEnabled(Level) bool {
// 	return true
// }
// func (log *DefaultLogger) Debug(args ...interface{}) { log.newEntry().Debug(args...) }
// func (log *DefaultLogger) Info(args ...interface{})  { log.newEntry().Info(args...) }
// func (log *DefaultLogger) Warn(args ...interface{})  { log.newEntry().Warn(args...) }
// func (log *DefaultLogger) Error(args ...interface{}) { log.newEntry().Error(args...) }
// func (log *DefaultLogger) Fatal(args ...interface{}) { log.newEntry().Fatal(args...) }
// func (log *DefaultLogger) Trace(args ...interface{}) { log.newEntry().Trace(args...) }

// func formatT(entry *Entry) {
// 	fmt.Printf("[%s][%v][%s]\t%s\n", entry.Time.Local().Format(time.RFC3339Nano), entry.Level, entry.Prefix, entry.Message)
// }

// func formatLP(entry *Entry) {
// 	tags := ",level=" + entry.Level.String()
// 	for key, val := range entry.Tags {
// 		tags += ","
// 		tags += key
// 		tags += "="
// 		tags += fmt.Sprint(val)
// 	}
// 	fmt.Printf("%s%s %s\n", entry.Prefix, tags, entry.Time.Local().Format(time.RFC3339Nano))
// }

type Logger struct {
	Out       io.Writer
	Formatter Formatter

	entryPool sync.Pool
	mu        sync.Mutex
}

func (logger *Logger) newEntry() *Entry {
	entry, ok := logger.entryPool.Get().(*Entry)
	if ok {
		return entry
	}
	return NewEntry(logger)
}
func (logger *Logger) releaseEntry(entry *Entry) {
	entry.Logger = nil
	entry.Data = make(Fields, 5) // reset
	logger.entryPool.Put(entry)
}
func (logger *Logger) level() Level {
	return DebugLevel
}

func (logger *Logger) format(entry *Entry, buf *bytes.Buffer) error {

	return logger.Formatter.Format(entry, buf)
}
func (logger *Logger) write(buf *bytes.Buffer) error {
	_, err := buf.WriteTo(logger.Out)
	return err
}

func (logger *Logger) Debug(a ...interface{}) { logger.newEntry().Debug(a...) }
func (logger *Logger) Info(a ...interface{})  { logger.newEntry().Info(a...) }
func (logger *Logger) Warn(a ...interface{})  { logger.newEntry().Warn(a...) }
func (logger *Logger) Error(a ...interface{}) { logger.newEntry().Error(a...) }
func (logger *Logger) Fatal(a ...interface{}) { logger.newEntry().Fatal(a...) }
func (logger *Logger) Trace(a ...interface{}) { logger.newEntry().Trace(a...) }

type ISpanContext interface {
	ILogger
	GetTraceId() string
}

type SpanContext struct {
	Context      context.Context
	TranceID     string
	SpanID       string
	SpanParentID string
}

func (span *SpanContext) GetTraceId() string { return span.TranceID }

func (span *Span) newEntry() *Entry {
	entry := span.logger.newEntry()
	traceId := common.Message("bbq")
	entry.Data[traceId.Key] = traceId
	return entry
}

type Span struct {
	ctx    context.Context
	logger *Logger
}

func (span *Span) Debug(a ...interface{}) { span.newEntry().Debug(a...) }
func (span *Span) Info(a ...interface{})  { span.newEntry().Info(a...) }
func (span *Span) Warn(a ...interface{})  { span.newEntry().Warn(a...) }
func (span *Span) Error(a ...interface{}) { span.newEntry().Error(a...) }
func (span *Span) Fatal(a ...interface{}) { span.newEntry().Fatal(a...) }
func (span *Span) Trace(a ...interface{}) { span.newEntry().Trace(a...) }

func (logger *Logger) SpanWithContext(ctx context.Context) *Span {
	if ctx == nil {
		ctx = context.TODO()
	}
	ctx = context.WithValue(ctx, "trace-id", "")
	ctx = context.WithValue(ctx, "span-id", "")
	ctx = context.WithValue(ctx, "span-parentid", "")

	var ic ISpanContext
	ic.Info("")

	return &Span{ctx, logger}
}

func (logger *Logger) SpanFromContext(ctx context.Context) *Span {
	span, _ := ctx.Value("goost.logs.span").(*Span)
	if span == nil {
		span = logger.SpanWithContext(ctx)
	}
	return span
}
