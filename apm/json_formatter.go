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

func (f *JsonFormatter) DoFormat(entry *field.Field, dest *bytes.Buffer, befor, after func()) (err error) {
	var skip []*field.Field
	for _, s := range f.SkipFields {
		fs, err := field.Find(entry, s)
		if err != nil {
			return err
		}
		skip = append(skip, fs...)
	}
	var items []*field.Field
LOOP:
	for _, it := range entry.Items {
		for _, s := range skip {
			if it == s {
				continue LOOP
			}
		}
		items = append(items, it)
	}
	if len(items) == 0 {
		return
	}
	befor()
	_, err = f.MarshalGroup(items, dest, skip...)
	if err != nil {
		return err
	}
	after()
	return
}

func (f *JsonFormatter) Format(entry *field.Field, dest *bytes.Buffer) (err error) {
	if err = f.DoFormat(entry, dest, func() {}, func() {}); err != nil {
		return err
	}
	err = dest.WriteByte('\n')
	return
}
