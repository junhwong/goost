package apm

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/junhwong/goost/apm/field"
	"github.com/junhwong/goost/apm/field/loglevel"
)

var _ Formatter = (*TextFormatter)(nil)

// JSON 格式
type TextFormatter struct {
	TimeLayout string
	SkipFields []string
	Color      bool // 是否打印颜色
	Skipped    func(*field.Field)
}

func getColor(lvl loglevel.Level, supportColor bool) (start, end string) {
	if !supportColor {
		return
	}
	// https://en.wikipedia.org/wiki/ANSI_escape_code#Colors
	// https://juejin.cn/post/6920241597846126599
	switch lvl {
	case loglevel.Debug:
		start = "\033[1;30;49m" // 34
		end = "\033[0m"
	case loglevel.Info:
		start = "\033[1;32;49m"
		end = "\033[0m"
	case loglevel.Warn:
		start = "\033[1;33;49m"
		end = "\033[0m"
	case loglevel.Error:
		start = "\033[1;31;49m"
		end = "\033[0m"
	case loglevel.Fatal:
		start = "\033[1;91;49m"
		end = "\033[0m"
	case loglevel.Trace:
		start = "\033[1;30;2m"
		end = "\033[0m"
	default:
	}

	return
}
func (tf *TextFormatter) getFileName(s string) string {
	arr := strings.FieldsFunc(s, func(r rune) bool {
		switch r {
		case '\\', '/':
			return true
		}
		return false
	})
	if len(arr) >= 2 {
		return strings.Join(arr[len(arr)-2:], "/")
	}
	if len(arr) == 1 {
		return arr[0]
	}
	return s
}

func (tf *TextFormatter) Format(entry *field.Field, dest *bytes.Buffer) (err error) {
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
	skipFields := map[string]string{
		TimeKey.Name():    "",
		MessageKey.Name(): "",
		LevelKey.Name():   "",
		TraceIDKey.Name(): "",
		"source":          "",
	}
	for _, v := range tf.SkipFields {
		skipFields[v] = ""
	}
	supportColor := tf.Color

	lvl := GetLevel(entry)
	cp, cs := getColor(lvl, supportColor)
	fprintf(`%s%s%s`, cp, lvl.Short(), "")
	start := strings.ReplaceAll(cp, "1;", "")
	if !supportColor {
		start = ""
	}
	fprintf(`%s`, start)

	if t := getTime(entry); !t.IsZero() {
		// TODO: 时区
		layout := tf.TimeLayout
		if layout == "" {
			layout = "15:04:05.000" // 20060102 15:04:05.000
		}
		fprintf(`%s`, t.Format(layout))
	}

	// writeByte('|')
	writeByte(' ')

	if f := entry.GetItem(TraceIDKey.Name()); f != nil {
		fprintf(`%s`, f.GetString())
		writeByte(' ')
	} else {
		fprintf(`%s`, "-")
		writeByte(' ')
	}

	if f := entry.GetItem("source"); f != nil && f.IsGroup() {
		if ff := f.GetItem("file"); ff != nil {
			fprintf(tf.getFileName(ff.GetString()))
			if fl := f.GetItem("line"); fl != nil {
				writeByte(':')
				fprintf(`%v`, fl.GetInt())
			}
		} else {
			fprintf(`%s`, "-")
		}
	} else {
		fprintf(`%s`, "-")
	}

	fprintf(`] `)
	fprintf(`%s`, cs)

	if val := getMessage(entry); val != "" {
		fprintf(`%s`, val)
	}

	start = "\033[2m"
	end := "\033[0m"
	if !supportColor {
		start = ""
		end = ""
	}

	fj := NewJsonFormatter()
	fj.Pretty = true
	fj.DurationToString = true
	fj.SkipFields = []string{
		"$." + TimeKey.Name(),
		"$." + MessageKey.Name(),
		"$." + LevelKey.Name(),
		"$.trace_id",
		"$.source",
	}
	if err := fj.DoFormat(entry, dest, func() {
		fprintf("\n%s", start)
	}, func() {
		fprintf(end)
	}); err != nil {
		fmt.Printf("err: %v\n", err)
	}
	writeByte('\n')
	return
}
