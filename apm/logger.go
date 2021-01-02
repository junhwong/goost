package apm

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"sync"
	"time"

	"github.com/junhwong/goost/runtime"
)

type LoggerInterface interface {
	Debug(...interface{})
	Info(...interface{})
	Warn(...interface{})
	Error(...interface{})
	Fatal(...interface{})
}

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
	queue     chan *Entry
	cancel    context.CancelFunc
}

func (logger *Logger) equeue(entry *Entry) {
	entry.equeue = nil
	logger.queue <- entry
}

func (logger *Logger) newEntry() *Entry {
	entry, ok := logger.entryPool.Get().(*Entry)
	if !ok {
		entry = NewEntry(logger.equeue)
	}
	entry.equeue = logger.equeue
	entry.Time = time.Now()
	entry.Data = make(Fields, 5) // reset
	return entry
}
func (logger *Logger) releaseEntry(entry *Entry) {
	entry.equeue = nil
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
func (logger *Logger) handle(entry *Entry) {
	if entry == nil {
		return
	}
	buf := bufferPool.Get().(*bytes.Buffer)
	defer bufferPool.Put(buf)
	buf.Reset()
	err := logger.format(entry, buf)
	if err != nil {
		logger.mu.Lock()
		fmt.Fprintf(os.Stderr, "Failed to obtain reader, %v\n", err)
		logger.mu.Unlock()
		return
	}

	err = logger.write(buf)
	if err != nil {
		logger.mu.Lock()
		fmt.Fprintf(os.Stderr, "Failed to write to log, %v\n", err)
		logger.mu.Unlock()
		return
	}
}
func (logger *Logger) flush() {
	for {
		select {
		case entry := <-logger.queue:
			logger.handle(entry)
		default:
			return
		}
	}
}

func (logger *Logger) Run(stopCh runtime.StopCh) {
	for {
		select {
		case entry := <-logger.queue:
			logger.handle(entry)
		case <-stopCh:
			return
		}
	}
}

func (logger *Logger) Debug(a ...interface{}) { logger.newEntry().Debug(a...) }
func (logger *Logger) Info(a ...interface{})  { logger.newEntry().Info(a...) }
func (logger *Logger) Warn(a ...interface{})  { logger.newEntry().Warn(a...) }
func (logger *Logger) Error(a ...interface{}) { logger.newEntry().Error(a...) }
func (logger *Logger) Fatal(a ...interface{}) { logger.newEntry().Fatal(a...) }
func (logger *Logger) Trace(a ...interface{}) { logger.newEntry().Trace(a...) }
