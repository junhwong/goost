package apm

import (
	"bytes"
	"fmt"
	"os"

	"github.com/junhwong/goost/apm/level"
)

var _ Handler = (*ConsoleHandler)(nil)

// 控制台
type ConsoleHandler struct {
	Formatter Formatter
}

func (h *ConsoleHandler) Priority() int {
	return -9999
}

func (h *ConsoleHandler) Handle(entry Entry, next func()) {
	defer next()
	lvl := GetLevel(entry)

	// TODO: 临时开发
	if lvl == level.Trace {
		return
	}

	err := UseBuffer(func(buf *bytes.Buffer) error {
		if err := h.Formatter.Format(entry, buf); err != nil {
			return err
		}

		out := os.Stdout
		if lvl >= level.Error && lvl < level.Trace {
			out = os.Stderr
		}

		// TODO 检查是否全部写入？
		if _, err := buf.WriteTo(out); err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		fmt.Fprintf(os.Stderr, "apm: Failed to handle to log, %+v\n", err)
	}
}
