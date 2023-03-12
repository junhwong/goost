package apm

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"sort"
	"sync"
)

// 日志项处理器
type Handler interface {

	// 优先级. 值越大越优先
	Priority() int

	// 刷新到输出
	Flush()

	// 处理日志
	Handle(entry Entry, next, end func())
}

type handlerSlice []Handler

func (x handlerSlice) Len() int           { return len(x) }
func (x handlerSlice) Less(i, j int) bool { return x[i].Priority() > x[j].Priority() }
func (x handlerSlice) Swap(i, j int)      { x[i], x[j] = x[j], x[i] }
func (x handlerSlice) Sort()              { sort.Sort(x) }
func (x handlerSlice) handle(entry Entry) {
	size := x.Len()
	crt := 0
	var once sync.Once
	var release = func() {
		once.Do(func() {
			crt += size
			// todo 将entry释放
		})
	}

	var next func()
	next = func() {
		if crt >= size {
			release()
			return
		}
		h := x[crt]
		crt++
		h.Handle(entry, next, release)
	}
	next()
}

func Console() (*SimpleHandler, *TextFormatter) {
	text := &TextFormatter{SkipFields: []string{"log.component"}}
	if os.Getenv("GOOST_APM_CONSOLE_COLOR") == "1" {
		text.Color = true
	}
	return &SimpleHandler{
		IsEnd:           true,
		Formatter:       text,
		HandlerPriority: -9000,
		Filter: func(entry Entry) bool {
			l := entry.GetLevel()
			return l >= LevelDebug && l < LevelTrace
		},
	}, text
}

var _ Handler = (*SimpleHandler)(nil)

type SimpleHandler struct {
	Out             Outer // 如果未提供则输出到控制台
	Formatter       Formatter
	Filter          func(Entry) bool
	HandlerPriority int
	IsEnd           bool
}

func (h SimpleHandler) Priority() int {
	return h.HandlerPriority
}

type Outer interface {
	io.Writer
	// Sync() error
}

func (h SimpleHandler) Flush() {}
func (h SimpleHandler) Handle(entry Entry, next, end func()) {
	if filter := h.Filter; filter != nil && !filter(entry) {
		if h.IsEnd {
			end()
			return
		}
		next()
		return
	}
	defer end()

	out := h.Out

	if out == nil {
		out = os.Stdout
		if lvl := entry.GetLevel(); lvl >= LevelError && lvl < LevelTrace {
			out = os.Stderr
		}
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
