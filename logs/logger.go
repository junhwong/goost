package logs

import (
	"bytes"
	"io"
	"sync"
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
