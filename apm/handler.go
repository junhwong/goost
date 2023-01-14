package apm

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"sort"
)

// 日志项处理器
type Handler interface {

	// 优先级. 值越大越优先
	Priority() int

	// 处理日志
	Handle(entry Entry, next func())
}

type handlerSlice []Handler

func (x handlerSlice) Len() int           { return len(x) }
func (x handlerSlice) Less(i, j int) bool { return x[i].Priority() > x[j].Priority() }
func (x handlerSlice) Swap(i, j int)      { x[i], x[j] = x[j], x[i] }
func (x handlerSlice) Sort()              { sort.Sort(x) }
func (x handlerSlice) handle(entry Entry) {
	size := x.Len()
	crt := 0
	var next func()
	next = func() {
		if crt >= size {
			return
		}
		h := x[crt]
		crt++
		h.Handle(entry, next)
	}
	next()
}

var _ Handler = (*SimpleHandler)(nil)

func Console() *SimpleHandler {
	text := &TextFormatter{}
	if os.Getenv("GOOST_APM_CONSOLE_COLOR") == "1" {
		text.Color = true
	}
	return &SimpleHandler{
		Out:             os.Stdout,
		Formatter:       text,
		HandlerPriority: -9999,
	}
}

// 控制台
type SimpleHandler struct {
	Out             io.Writer
	Formatter       Formatter
	Level           Level
	HandlerPriority int
}

func (h SimpleHandler) Priority() int {
	return h.HandlerPriority
}

func (h SimpleHandler) Handle(entry Entry, next func()) {
	defer next()
	lvl := entry.GetLevel()

	// TODO: 临时开发
	if lvl < h.Level {
		return
	}

	var out io.Writer = h.Out
	if lvl >= LevelError && lvl < LevelTrace {
		out = os.Stderr
	}
	if out == nil {
		out = os.Stdout
	}

	err := UseBuffer(func(buf *bytes.Buffer) error {
		if err := h.Formatter.Format(entry, buf); err != nil {
			return err
		}

		// TODO 检查是否全部写入？
		if _, err := buf.WriteTo(out); err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		fmt.Fprintf(os.Stderr, "apm: Failed to handle, %v: %+v\n", err, entry)
	}
}
