package deflog

import (
	"bytes"
	"encoding/json"
	"fmt"
	"time"

	"github.com/junhwong/goost/apm"
	"github.com/spf13/cast"
)

func NewJsonFormatter(timeLayout ...string) *JsonFormatter {
	f := &JsonFormatter{}
	for _, v := range timeLayout {
		f.timeLayout = v
	}
	if f.timeLayout == "" {
		f.timeLayout = time.RFC3339Nano
	}
	return f
}

var _ apm.Formatter = (*JsonFormatter)(nil)

// JSON 格式
type JsonFormatter struct {
	timeLayout string
}

func (jf *JsonFormatter) Format(entry apm.Entry, dest *bytes.Buffer) (err error) {
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

	writeByte('{')

	fprintf(`"level":%q`, entry.GetLevel().String())
	fs := entry.GetFields()
	if val := fs.Get(apm.TimeKey); val != nil {
		t, err := cast.ToTimeE(val)
		if err == nil && !t.IsZero() {
			// TODO: 时区
			fprintf(`,"time":%q`, t.Format(jf.timeLayout))
		}
	}

	if val := fs.Get(apm.MessageKey); val != nil {
		fprintf(`,"message":%q`, val)
	}

	// TODO 折叠map
	for _, f := range fs {
		key, val := f.Unwrap()
		if key == nil || val == nil {
			continue
		}
		if key == apm.TimeKey || key == apm.MessageKey || key == apm.LevelKey {
			continue
		}

		if key == apm.TracebackPathKey || key == apm.TracebackLineNoKey {
			continue
		}

		if key == apm.TracebackCallerKey {
			val = fmt.Sprintf("%s:%v", val, fs.Get(apm.TracebackLineNoKey, 0))
		}

		var data []byte

		if data, err = json.Marshal(val); err != nil {
			return
		}

		name := key.Name() // TrimFieldNamePrefix(it.Key.Name())

		if len(name) == 0 {
			panic(fmt.Sprintln("apm: entry key name is required"))
		}

		switch name {
		case "level", "time", "message":
			name = "data." + name
		}

		fprintf(`,%q:%s`, name, data)
	}

	writeByte('}')
	writeByte('\n')
	return
}
