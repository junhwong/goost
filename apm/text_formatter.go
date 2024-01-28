package apm

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/junhwong/goost/apm/field"
	"github.com/junhwong/goost/apm/field/loglevel"
	"github.com/spf13/cast"
)

var _ Formatter = (*TextFormatter)(nil)

// JSON 格式
type TextFormatter struct {
	TimeLayout string
	SkipFields []string
	Color      bool // 是否打印颜色
	Skipped    func(*field.Field)
}

func cutstr(v interface{}, l int) string {
	s := cast.ToString(v)
	if ls := len(s); ls >= l {
		return s[ls-l:]
	}
	return s
}

func getColor(lvl loglevel.Level, supportColor bool) (start, end string) {
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
	default:
	}

	if supportColor {
		return
	}
	return "", ""
}
func (tf *TextFormatter) Format(entry Entry, dest *bytes.Buffer) (err error) {
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
	}
	for _, v := range tf.SkipFields {
		skipFields[v] = ""
	}
	supportColor := tf.Color

	lvl := entry.GetLevel()
	cp, cs := getColor(lvl, supportColor)
	fprintf(`%s%s%s`, cp, lvl.Short(), "")
	start := strings.ReplaceAll(cp, "1;", "") //"\033[1;31;40m"
	// end := ""   //"\033[0m"
	if !supportColor {
		start = ""
		// end = ""
	}
	fprintf(`%s`, start)

	if t := entry.GetTime(); !t.IsZero() {
		// TODO: 时区
		layout := tf.TimeLayout
		if layout == "" {
			layout = "15:04:05.000" // 20060102 15:04:05.000
		}

		fprintf(`%s`, t.Format(layout))
	}
	// writeByte('|')
	writeByte(' ')
	fs := entry.GetFields()

	if ci := entry.GetCallerInfo(); ci != nil {
		fprintf(`%s`, ci.Caller())
	}

	fprintf(`%s`, cs)
	// writeByte('\n')
	writeByte(' ')

	if val := entry.GetMessage(); val != "" {
		fprintf(`%s`, val)
	}

	// keys := fs.Keys()
	fsv := []string{}
	for _, f := range fs {
		// key, val := f.Unwrap()
		// fmt.Printf("key: %v\n", key)
		// key == nil || val == nil ||
		if _, ok := skipFields[f.GetName()]; ok {
			if fn := tf.Skipped; fn != nil {
				fn(f)
			}
			continue
		}
		// if key == nil {
		// 	fmt.Printf("skipped field nil key: %v\n", f)
		// 	continue
		// }
		// if val == nil {
		// 	fmt.Printf("skipped field key: %v\n", f)
		// 	continue
		// }
		// if skipFields.Has(key.Name()) {
		// 	continue
		// }

		switch f.GetName() {
		// case TraceIDKey: // TODO: 开发者选项
		// 	continue
		case TracebackCallerKey.Name(), ErrorStackTraceKey.Name(), TracebackPathKey.Name(),
			TracebackLineNoKey.Name(): // TODO: 调用者选项
			continue
		}
		name := f.GetName()

		// if name == "error.method" {
		// 	if s, _ := val.(string); len(s) > 0 {

		// 		var tmp []string
		// 		arr := strings.Split(s, ",")
		// 		i:=len(arr)-1
		// 		for i>-1{
		// 			s:=arr[i]
		// 			if strings.HasPrefix(s,"cobra@"){
		// 				continue
		// 			}
		// 			if strings.HasPrefix(s,"dig@"){
		// 				continue
		// 			}
		// 		}

		// 	}
		// }
		val := field.GetValue(f)
		if e, _ := val.(error); e != nil {
			val = e.Error()
		}
		data, err := json.Marshal(val)
		if err != nil {
			return err
		}

		if bytes.Equal(data, []byte{'{', '}'}) {
			continue
		}

		if len(data) > 1024 {
			data = data[:1024]
		}

		switch name {
		case "level", "time", "message":
			name = "data." + name
		}

		fsv = append(fsv, fmt.Sprintf("%q:%s", name, data))
	}

	if len(fsv) > 0 {
		start := "\033[2m"
		end := "\033[0m"
		if !supportColor {
			start = ""
			end = ""
		}
		fprintf(" %s\n{\n%s\n}%s", start, strings.Join(fsv, ",\n"), end)
	}
	writeByte('\n')
	return
}

// func (jf *TextFormatter) Format2(entry logEntry, dest *bytes.Buffer) (err error) {
// 	writeByte := func(c byte) {
// 		if err != nil {
// 			return
// 		}
// 		err = dest.WriteByte(c)
// 	}

// 	fprint := func(a ...string) {
// 		for _, it := range a {
// 			if err != nil {
// 				break
// 			}
// 			_, err = dest.WriteString(it)
// 		}
// 	}

// 	fprint(level.Short(entry.Level))

// 	if !entry.Time.IsZero() {
// 		// TODO: 时区
// 		fprint(entry.Time.Format(jf.timeLayout))
// 	}

// 	writeByte('|')
// 	fprint(entry.Caller.Method, ":", strconv.Itoa(entry.Caller.Line))

// 	// writeByte(']')
// 	writeByte('\n')

// 	if entry.Message != "" {
// 		fprint(entry.Message)
// 	}

// 	fs := []string{}
// 	for key, val := range entry.Fields {
// 		if key == nil || val == nil {
// 			continue
// 		}

// 		switch key {
// 		case TimeKey, MessageKey, LevelKey: // 已经处理
// 			continue
// 		// case TraceIDKey: // TODO: 开发者选项
// 		// 	continue
// 		case TracebackCallerKey, TracebackPathKey, TracebackLineNoKey: // TODO: 调用者选项
// 			continue
// 		}

// 		// if key == TracebackCallerKey {
// 		// 	val = fmt.Sprintf("%s:%v", val, entry.Get(TracebackLineNoKey, 0))
// 		// }

// 		var data []byte

// 		if data, err = json.Marshal(val); err != nil {
// 			return
// 		}

// 		name := key.Name() // TrimFieldNamePrefix(it.Key.Name())

// 		// if len(name) == 0 {
// 		// 	fmt.Println("apm: skip entry: name") // TODO devop log
// 		// 	continue
// 		// }

// 		switch name {
// 		case "level", "time", "message":
// 			name = "data." + name
// 		}
// 		fs = append(fs, `"`+name+`":`+string(data))
// 	}
// 	if len(fs) > 0 {
// 		fprint(" {", strings.Join(fs, ",\n"), "}")
// 	}
// 	writeByte('\n')
// 	return
// }
