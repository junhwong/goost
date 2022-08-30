package apm

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"
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
		f.timeLayout = "20060102 15:04:05.000"
	}
	return f
}

var _ Formatter = (*TextFormatter)(nil)

// JSON 格式
type TextFormatter struct {
	timeLayout string
}

func cutstr(v interface{}, l int) string {
	s := cast.ToString(v)
	if ls := len(s); ls >= l {
		return s[ls-l:]
	}
	return s
}

var (
	supportColor = false
)

func getColor(lvl level.Level) (start, end string) {

	// https://en.wikipedia.org/wiki/ANSI_escape_code#Colors
	// https://juejin.cn/post/6920241597846126599
	switch lvl {
	case level.Debug:
		start = "\033[1;30;49m" // 34
		end = "\033[0m"
	case level.Info:
		start = "\033[1;32;49m"
		end = "\033[0m"
	case level.Warn:
		start = "\033[1;33;49m"
		end = "\033[0m"
	case level.Error:
		start = "\033[1;31;49m"
		end = "\033[0m"
	case level.Fatal:
		start = "\033[1;91;49m"
		end = "\033[0m"
	default:
	}
	if supportColor {
		return
	}
	return "", ""
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
	lvl := GetLevel(entry)
	cp, cs := getColor(lvl)
	fprintf(`%s%s%s`, cp, level.Short(lvl), "")
	start := strings.ReplaceAll(cp, "1;", "") //"\033[1;31;40m"
	// end := ""   //"\033[0m"
	if !supportColor {
		start = ""
		// end = ""
	}
	fprintf(`%s`, start)

	if val := entry.Get(TimeKey); val != nil {
		t, err := cast.ToTimeE(val)
		if err == nil && !t.IsZero() {
			// TODO: 时区
			// v := float64(t.UnixMilli()) / 1000
			// fprintf(`%0.6f`, v)
			fprintf(`%s`, t.Format(jf.timeLayout))
		}
	}
	writeByte('|')
	fprintf(`%s`, entry.Get(TracebackPathKey))
	fprintf(`:%v`, entry.Get(TracebackLineNoKey))
	// writeByte(']')

	fprintf(`%s`, cs)
	writeByte('\n')

	if val := entry.Get(MessageKey); val != nil {
		fprintf(`%s`, val)
	}

	keys := entry.Keys()
	fs := []string{}
	for _, key := range keys {
		val := entry[key]
		if key == nil || val == nil {
			continue
		}

		switch key {
		case TimeKey, MessageKey, LevelKey: // 已经处理
			continue
		// case TraceIDKey: // TODO: 开发者选项
		// 	continue
		case TracebackCallerKey, TracebackPathKey, TracebackLineNoKey: // TODO: 调用者选项
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
		fs = append(fs, fmt.Sprintf("%q:%s", name, data))
	}

	if len(fs) > 0 {
		start := "\033[2m"
		end := "\033[0m"
		if !supportColor {
			start = ""
			end = ""
		}
		fprintf(" %s{\n%s\n}%s", start, strings.Join(fs, ",\n"), end)
	}
	writeByte('\n')
	return
}

func (jf *TextFormatter) Format2(entry logEntry, dest *bytes.Buffer) (err error) {
	writeByte := func(c byte) {
		if err != nil {
			return
		}
		err = dest.WriteByte(c)
	}

	fprint := func(a ...string) {
		for _, it := range a {
			if err != nil {
				break
			}
			_, err = dest.WriteString(it)
		}
	}

	fprint(level.Short(entry.Level))

	if !entry.Time.IsZero() {
		// TODO: 时区
		fprint(entry.Time.Format(jf.timeLayout))
	}

	writeByte('|')
	fprint(entry.Caller.Method, ":", strconv.Itoa(entry.Caller.Line))

	// writeByte(']')
	writeByte('\n')

	if entry.Message != "" {
		fprint(entry.Message)
	}

	fs := []string{}
	for key, val := range entry.Fields {
		if key == nil || val == nil {
			continue
		}

		switch key {
		case TimeKey, MessageKey, LevelKey: // 已经处理
			continue
		// case TraceIDKey: // TODO: 开发者选项
		// 	continue
		case TracebackCallerKey, TracebackPathKey, TracebackLineNoKey: // TODO: 调用者选项
			continue
		}

		// if key == TracebackCallerKey {
		// 	val = fmt.Sprintf("%s:%v", val, entry.Get(TracebackLineNoKey, 0))
		// }

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
		fs = append(fs, `"`+name+`":`+string(data))
	}
	if len(fs) > 0 {
		fprint(" {", strings.Join(fs, ",\n"), "}")
	}
	writeByte('\n')
	return
}
