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

	if f := entry.Get(TimeKey); f != nil && f.Valid() {
		t, err := cast.ToTimeE(f.Value)
		if err == nil && !t.IsZero() {
			fprintf(`,"time":%q`, t.Format(jf.timeLayout))
		}
	}

	if f := entry.Get(MessageKey); f != nil && f.Valid() {
		fprintf(`,"message":%q`, f.Value)
	}

	for _, it := range entry {
		if !it.Valid() {
			continue
		}
		if it.Key == TimeKey || it.Key == MessageKey || it.Key == LevelKey {
			continue
		}

		var data []byte

		if data, err = json.Marshal(it.Value); err != nil {
			return
		}

		name := TrimFieldNamePrefix(it.Key.Name())

		if len(name) == 0 {
			fmt.Println("apm: skip entry: name") // TODO devop log
			continue
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
