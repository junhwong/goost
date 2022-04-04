package apm

import (
	"bytes"
	"encoding/json"
	"fmt"
	"time"

	"github.com/junhwong/goost/apm/level"
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

var _ Formatter = (*JsonFormatter)(nil)

// JSON 格式
type JsonFormatter struct {
	timeLayout string
}

func (jf *JsonFormatter) Format(entry Entry, dest *bytes.Buffer) (err error) {
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

	fprintf(`"level":%q`, level.String(GetLevel(entry)))

	if val := entry.Get(TimeKey); val != nil {
		t, err := cast.ToTimeE(val)
		if err == nil && !t.IsZero() {
			// TODO: 时区
			fprintf(`,"time":%q`, t.Format(jf.timeLayout))
		}
	}

	if val := entry.Get(MessageKey); val != nil {
		fprintf(`,"message":%q`, val)
	}

	for key, val := range entry {
		if key == nil || val == nil {
			continue
		}
		if key == TimeKey || key == MessageKey || key == LevelKey {
			continue
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

		fprintf(`,%q:%s`, name, data)
	}

	writeByte('}')
	writeByte('\n')
	return
}
