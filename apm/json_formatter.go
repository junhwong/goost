package apm

import (
	"bytes"
	"strings"
	"time"

	"github.com/junhwong/goost/apm/field"
)

var _ Formatter = (*JsonFormatter)(nil)

// JSON 格式
type JsonFormatter struct {
	field.JsonMarshaler

	TimeName    string
	LevelName   string
	MessageName string
	SkipFields  []string
}

func NewJsonFormatter() *JsonFormatter {
	f := &JsonFormatter{
		JsonMarshaler: field.JsonMarshaler{
			EscapeHTML: true,
			OmitEmpty:  true,
			TimeLayout: time.RFC3339Nano,
		},
		LevelName:   "level",
		MessageName: "msg",
		TimeName:    "time",
	}

	f.NameFilter = func(s string) string {
		switch s {
		case TimeKey.Name():
			return f.TimeName
		case LevelKey.Name():
			return f.LevelName
		case MessageKey.Name():
			return f.MessageName
		}

		if strings.HasPrefix(s, "__") {
			return "-"
		}

		for _, it := range f.SkipFields {
			if it == s {
				return "-"
			}
		}

		return s
	}

	return f
}

func (f *JsonFormatter) Format(entry *field.Field, dest *bytes.Buffer) (err error) {
	// m := f.JsonMarshaler // copy
	f.MarshalGroup(entry.Items, dest)
	dest.WriteByte('\n')
	return
}
