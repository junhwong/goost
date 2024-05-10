package field

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"slices"
	"strconv"
	"time"
)

type Marshaler struct {
	err error
	w   io.Writer
	n   int64

	OmitEmpty  bool
	NameFilter func(string) string
	NameLess   func(a, b *Field) int
}

func (m *Marshaler) Marshal(f *Field, w io.Writer) (int64, error) {
	m.w = w
	m.err = nil
	m.n = 0
	m.write(f)
	return m.n, m.err
}
func (m *Marshaler) writeByte(c byte) {
	if m.err != nil {
		return
	}
	n, err := m.w.Write([]byte{c})
	m.err = err
	m.n += int64(n)
}
func (m *Marshaler) writeBytes(p []byte) {
	if m.err != nil {
		return
	}
	n, err := m.w.Write(p)
	m.err = err
	m.n += int64(n)
}

func (m *Marshaler) write(f *Field) {
	if m.err != nil {
		return
	}

	if f.Type == InvalidKind {
		return
	}

	if f.IsGroup() {
		if f.IsNull() {
			m.writeBytes([]byte("null"))
			return
		}
		has := false
		m.writeByte('{')
		items := f.Items
		if m.NameLess != nil {
			items = make([]*Field, len(items))
			copy(items, f.Items)
			slices.SortFunc(items, m.NameLess)
		}
		for _, it := range items {
			if it.Type == InvalidKind {
				continue
			}
			name := it.GetName()
			if m.NameFilter != nil {
				name = m.NameFilter(name)
			}
			fmt.Printf("nameeeee: %v\n", name)
			if name == "-" {
				continue
			}
			if has {
				m.writeByte(',')
			}

			has = true
			m.writeBytes([]byte(`"`))
			m.writeBytes([]byte(name))
			m.writeBytes([]byte(`":`))
			m.write(it)
		}
		m.writeByte('}')
		return
	}

	if f.IsArray() {
		if f.IsNull() {
			m.writeBytes([]byte("null"))
			return
		}
		has := false
		m.writeByte('[')
		for _, it := range f.Items {
			if it.Type == InvalidKind {
				continue
			}
			if has {
				m.writeByte(',')
			}
			has = true
			m.write(it)
		}
		m.writeByte(']')
		return
	}

	switch f.Type {
	case IntKind:
		v := f.GetInt()
		m.writeBytes(strconv.AppendInt(nil, v, 10))
	case UintKind:
		v := f.GetUint()
		m.writeBytes(strconv.AppendUint(nil, v, 10))
	case FloatKind:
		v := f.GetFloat()
		b := strconv.AppendFloat(nil, v, 'f', -1, 64) // -1
		if bytes.LastIndex(b, []byte{'.'}) == -1 {    // 保留数据类型
			b = append(b, '.', '0')
		}
		m.writeBytes(b)
	case BoolKind:
		v := f.GetBool()
		m.writeBytes(strconv.AppendBool(nil, v))
	case StringKind:
		v := f.GetString()
		data, err := json.Marshal(v) // 防止转义
		if err != nil {
			m.err = err
			return
		}
		m.writeBytes(data)
		// m.writeBytes([]byte(`"`))
		// m.writeBytes([]byte(v))
		// m.writeBytes([]byte(`"`))
	case TimeKind:
		v := f.GetTime()
		m.writeBytes([]byte(`"`))
		m.writeBytes([]byte(v.Format(time.RFC3339Nano)))
		m.writeBytes([]byte(`"`))
	case IPKind:
		v := f.GetIP()
		m.writeBytes([]byte(`"`))
		m.writeBytes([]byte(v.String()))
		m.writeBytes([]byte(`"`))
	case LevelKind:
		v := f.GetLevel()
		m.writeBytes([]byte(`"`))
		m.writeBytes([]byte(v.String()))
		m.writeBytes([]byte(`"`))
	case BytesKind:
		v := f.GetBytes()
		m.writeBytes([]byte(`"base64:`))
		e := base64.NewEncoder(base64.StdEncoding, m.w)
		_, m.err = e.Write(v)
		e.Close()
		m.writeBytes([]byte(`"`))
	default:
		panic("todo:" + f.Type.String())
	}
}
