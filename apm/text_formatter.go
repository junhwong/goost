package apm

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/junhwong/goost/apm/level"
	"github.com/spf13/cast"
)

func NewTextFormatter(timeLayout ...string) *TextFormatter {
	f := &TextFormatter{}
	for _, v := range timeLayout {
		f.timeLayout = v
	}
	if f.timeLayout == "" {
		f.timeLayout = "15:04:05"
	}
	return f
}

var _ Formatter = (*TextFormatter)(nil)

// JSON 格式
type TextFormatter struct {
	timeLayout string
}

func (jf *TextFormatter) Format(entry Entry, dest *bytes.Buffer) (err error) {
	writeByte := func(c byte) {
		if err != nil {
			return
		}
		err = dest.WriteByte(c)
	}
	fprintf := func(format string, a ...interface{}) {
		if err != nil {
			return
		}
		_, err = fmt.Fprintf(dest, format, a...)
	}

	fprintf(`%s`, level.Short(GetLevel(entry)))

	if val := entry.Get(TimeKey); val != nil {
		t, err := cast.ToTimeE(val)
		if err == nil && !t.IsZero() {
			// TODO: 时区
			fprintf(`%s`, t.Format(jf.timeLayout))
		}
	}
	writeByte('|')
	if val := entry.Get(MessageKey); val != nil {
		fprintf(`%s`, val)
	}

	fs := []string{}
	for key, val := range entry {
		if key == nil || val == nil {
			continue
		}

		switch key {
		case TimeKey, MessageKey, LevelKey: // 已经处理
			continue
		case TraceIDKey, TracebackCallerKey: // TODO: 开发者选项
			continue
		case TracebackPathKey, TracebackLineNoKey: // TODO: 调用者选项
			continue
		}

		if key == TracebackCallerKey {
			val = fmt.Sprintf("%s:%v", val, entry.Get(TracebackLineNoKey, 0))
		}

		var data []byte

		if data, err = json.Marshal(val); err != nil {
			return
		}

		name := key.Name() // TrimFieldNamePrefix(it.Key.Name())

		// if len(name) == 0 {
		// 	fmt.Println("apm: skip entry: name") // TODO devop log
		// 	continue
		// }

		switch name {
		case "level", "time", "message":
			name = "data." + name
		}
		fs = append(fs, fmt.Sprintf("%s=%s", name, data))
	}
	if len(fs) > 0 {
		fprintf(` {%s}`, strings.Join(fs, ","))
	}
	writeByte('\n')
	return
}
