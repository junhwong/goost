package apm

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"slices"
	"sort"
	"sync/atomic"

	"github.com/junhwong/goost/apm/field"
	"github.com/junhwong/goost/apm/field/loglevel"
)

// 日志项处理器
type Handler interface {

	// 优先级. 值越大越优先
	Priority() int

	// 刷新到输出
	Flush()

	// 处理日志
	Handle(entry *field.Field, next, end func())
}

type handlerSlice []Handler

func (x handlerSlice) Len() int           { return len(x) }
func (x handlerSlice) Less(i, j int) bool { return x[i].Priority() > x[j].Priority() }
func (x handlerSlice) Swap(i, j int)      { x[i], x[j] = x[j], x[i] }
func (x handlerSlice) Sort()              { sort.Sort(x) }

var _handlers atomic.Value

func GetHandlers() handlerSlice {
	obj := _handlers.Load()
	if obj == nil {
		return handlerSlice{}
	}
	return obj.(handlerSlice)
}
func addHandlers(dst handlerSlice, handlers []Handler) handlerSlice {
	for _, h := range handlers {
		if dst != nil {
			dst = append(dst, h)
		}
	}
	return dst
}
func AddHandlers(handlers ...Handler) {
	if len(handlers) == 0 {
		panic("handlers is empty")
	}
	gmu.Lock()
	defer gmu.Unlock()

	dst := addHandlers(GetHandlers(), handlers)
	dst.Sort()
	_handlers.Store(dst)
}

func removeHandlers(dst handlerSlice, handlers []Handler) handlerSlice {
	dst = slices.DeleteFunc(dst, func(h Handler) bool {
		for _, it := range handlers {
			if it == h {
				return true
			}
		}
		return false
	})
	return dst
}
func RemoveHandlers(handlers ...Handler) {
	if len(handlers) == 0 {
		panic("handlers is empty")
	}
	gmu.Lock()
	defer gmu.Unlock()

	dst := removeHandlers(GetHandlers(), handlers)
	dst.Sort()
	_handlers.Store(dst)
}

func SetHandlers(handlers ...Handler) {
	if len(handlers) == 0 {
		panic("handlers is empty")
	}

	gmu.Lock()
	defer gmu.Unlock()

	dst := addHandlers(make(handlerSlice, 0, len(handlers)), handlers)
	dst.Sort()
	_handlers.Store(dst)
}

func NewConsole() (*SimpleHandler, *TextFormatter) {
	text := &TextFormatter{SkipFields: []string{"log.component"}}
	if a := os.Getenv("GOOST_APM_CONSOLE_COLOR"); a == "1" {
		text.Color = true
	}
	return &SimpleHandler{
		IsEnd:           true,
		Formatter:       text,
		HandlerPriority: -9000,
		// Filter: func(entry Entry) bool {
		// 	l := entry.GetLevel()
		// 	return l >= field.LevelDebug && l < field.LevelTrace
		// },
	}, text
}

var _ Handler = (*SimpleHandler)(nil)

type SimpleHandler struct {
	Out             Outer // 如果未提供则输出到控制台
	Formatter       Formatter
	Filter          func(*field.Field) bool
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
func (h SimpleHandler) Handle(entry *field.Field, next, end func()) {
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
		if lvl := entry.GetLevel(); lvl >= loglevel.Error && lvl < loglevel.Trace2 {
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
