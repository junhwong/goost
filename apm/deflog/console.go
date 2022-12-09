package deflog

import (
	"bytes"
	"fmt"
	"io"
	"os"

	"github.com/junhwong/goost/apm"
)

var _ apm.Handler = (*ConsoleHandler)(nil)

// 控制台
type ConsoleHandler struct {
	Out       io.Writer
	Formatter apm.Formatter
	Level     apm.LogLevel
}

func (h *ConsoleHandler) Priority() int {
	return -9999
}

func (h *ConsoleHandler) Handle(entry apm.Entry, next func()) {
	defer next()
	lvl := entry.GetLevel()

	// TODO: 临时开发
	// if lvl == Trace {
	// 	return
	// }

	var out io.Writer = h.Out
	if lvl >= apm.Error && lvl < apm.Trace {
		out = os.Stderr
	}
	if out == nil {
		out = os.Stdout
	}
	// out = io.Discard

	err := apm.UseBuffer(func(buf *bytes.Buffer) error {
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
