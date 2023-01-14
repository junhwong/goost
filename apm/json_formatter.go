package apm

import (
	"bytes"
	"encoding/json"
	"fmt"
	"time"

	"github.com/junhwong/duzee-go/pkg/sets"
)

var _ Formatter = (*JsonFormatter)(nil)

// JSON 格式
type JsonFormatter struct {
	TimeLayout string
	SkipFields []string
}

func (f *JsonFormatter) Format(entry Entry, dest *bytes.Buffer) (err error) {
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
	skipFields := sets.NewString(TimeKey.Name(), MessageKey.Name(), LevelKey.Name())
	for _, v := range f.SkipFields {
		skipFields.Insert(v)
	}

	writeByte('{')

	fprintf(`"level":%q`, entry.GetLevel().String())

	if t := entry.GetTime(); !t.IsZero() {
		// TODO: 时区
		layout := f.TimeLayout
		if layout == "" {
			layout = time.RFC3339Nano
		}
		fprintf(`,"time":%q`, t.Format(layout))
	}

	if val := entry.GetMessage(); val != "" {
		fprintf(`,"message":%q`, val)
	}

	// TODO 折叠map
	fs := entry.GetFields()
	for _, f := range fs {
		key, val := f.Unwrap()
		if key == nil || val == nil {
			continue
		}
		if skipFields.Has(key.Name()) {
			continue
		}

		if key == TracebackPathKey || key == TracebackLineNoKey {
			continue
		}

		// if key == TracebackCallerKey {
		// 	val = fmt.Sprintf("%s:%v", val, fs.Get(TracebackLineNoKey, 0))
		// }

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
