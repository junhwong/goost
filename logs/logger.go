package logs

import (
	"fmt"
	"time"
)

// Logger 纯日志接口
type Logger interface {
	Debug(...interface{})
	Info(...interface{})
	Warn(...interface{})
	Error(...interface{})
	Fatal(...interface{})
	Trace(...interface{})
}

type EntryHandler interface {
	Handle(*Entry) error
	IsEnabled(Level) bool
}

type LoggerHandler interface {
	EntryHandler
}

type ForkedLogger interface {
	Logger
	IsEnabled(Level) bool
	WithPrefix(prefix string) ForkedLogger
	WithContext(prefix string) ForkedLogger
}

type RootLogger interface {
	ForkedLogger
	SetHandler(LoggerHandler) error
	SetLevel(Level)
}

type DefaultLogger struct {
	Prefix string
}

func (log *DefaultLogger) newEntry() *Entry {
	return &Entry{
		Prefix: log.Prefix,
		// Handler: log,
		Tags:   make(Fields),
		Data:   make(Fields),
		Fields: make([]*Field, 0),
	}
}

func (log *DefaultLogger) Handle(entry *Entry) error {
	formatT(entry)
	return nil
}

func (log *DefaultLogger) WithPrefix(prefix string) ForkedLogger {
	return log
}
func (log *DefaultLogger) WithContext(prefix string) ForkedLogger {
	return log
}
func (log *DefaultLogger) SetHandler(lhd LoggerHandler) error {
	return nil
}
func (log *DefaultLogger) SetLevel(lhd Level) {
}
func (log *DefaultLogger) IsEnabled(Level) bool {
	return true
}
func (log *DefaultLogger) Debug(args ...interface{}) { log.newEntry().Debug(args...) }
func (log *DefaultLogger) Info(args ...interface{})  { log.newEntry().Info(args...) }
func (log *DefaultLogger) Warn(args ...interface{})  { log.newEntry().Warn(args...) }
func (log *DefaultLogger) Error(args ...interface{}) { log.newEntry().Error(args...) }
func (log *DefaultLogger) Fatal(args ...interface{}) { log.newEntry().Fatal(args...) }
func (log *DefaultLogger) Trace(args ...interface{}) { log.newEntry().Trace(args...) }

func formatT(entry *Entry) {
	fmt.Printf("[%s][%v][%s]\t%s\n", entry.Time.Local().Format(time.RFC3339Nano), entry.Level, entry.Prefix, entry.Message)
}

func formatLP(entry *Entry) {
	tags := ",level=" + entry.Level.String()
	for key, val := range entry.Tags {
		tags += ","
		tags += key
		tags += "="
		tags += fmt.Sprint(val)
	}
	fmt.Printf("%s%s %s\n", entry.Prefix, tags, entry.Time.Local().Format(time.RFC3339Nano))
}
