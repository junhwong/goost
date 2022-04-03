package apm

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/junhwong/goost/apm/level"
	"github.com/junhwong/goost/pkg/field"
	"github.com/spf13/cast"
)

// Formatter 表示一个格式化器。
type Formatter interface {
	// Format 格式化一条日志。
	//
	// 注意：不要缓存 `entry`, `dest` 对象，因为它们是池化对象。
	Format(entry Entry, dest *bytes.Buffer) (err error)
}

var _ Formatter = (*JsonFormatter)(nil)

type JsonFormatter struct {
	layout string
}

func (f *JsonFormatter) Format(entry Entry, dest *bytes.Buffer) (err error) {
	settings := FormatSettings{TrimFieldPrefix: []string{"apm."}}
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
	fs := field.Fields(entry)
	lvl := level.String(entry.GetLevel())
	if lvl != "" {
		fprintf(`"level":%q`, lvl)
	}
	fs.Del(LevelKey)
	if f := fs.Del(TimeKey); f != nil && f.Valid() {
		t, err := cast.ToTimeE(f.Value)
		if err == nil && !t.IsZero() {
			fprintf(`,"time":%q`, t.Format(time.RFC3339Nano))
		}
	}

	if f := fs.Del(MessageKey); f != nil && f.Valid() {
		fprintf(`,"message":%q`, f.Value)
	}

	for _, it := range entry {
		if !it.Valid() {
			fmt.Println("apm: skip") // TODO devop log
			continue
		}
		val := it.Value // it.GetValue()
		// if val == nil {
		// 	continue
		// }
		var data []byte
		data, err = json.Marshal(val)
		if err != nil {
			return
		}

		name := it.Key.Name()
		for _, prefix := range settings.TrimFieldPrefix {
			name = strings.TrimSpace(strings.TrimPrefix(name, prefix))
		}
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

func (f *JsonFormatter) FormatWith(entry Entry, cb func(err error, buf *bytes.Buffer)) {
	buf := bufferPool.Get().(*bytes.Buffer)
	defer bufferPool.Put(buf)
	buf.Reset()
	cb(f.Format(entry, buf), buf)

}
