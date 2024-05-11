package field

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"io"
	"slices"
	"strconv"
	"time"
)

type Marshaler struct {
	err  error
	w    io.Writer
	n    int64
	buf  *bytes.Buffer
	enc  *json.Encoder
	benc io.WriteCloser

	OmitEmpty  bool // todo 实现
	NameFilter func(string) string
	NameLess   func(a, b *Field) int
	EscapeHTML bool
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
		if v == "" {
			m.writeBytes([]byte(`""`))
			return
		}
		if m.buf == nil {
			m.buf = bytes.NewBuffer(nil)
		}
		if m.enc == nil {
			m.enc = json.NewEncoder(m.buf)
		}

		m.enc.SetEscapeHTML(m.EscapeHTML)
		m.buf.Reset()

		if err := m.enc.Encode(v); err != nil {
			m.err = err
			return
		}

		var data []byte
		data = m.buf.Bytes()
		n := len(data) - 1
		if data[n] == '\n' { // json.NewEncoder 会增加一个换行
			data = data[:n]
		}
		m.writeBytes(data)

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
		if m.buf == nil {
			m.buf = bytes.NewBuffer(nil)
		}
		if m.benc == nil {
			m.benc = base64.NewEncoder(base64.StdEncoding, m.buf)
		}
		m.buf.Reset()
		_, m.err = m.benc.Write(v)
		m.writeBytes(m.buf.Bytes())
		m.writeBytes([]byte(`"`))
	default:
		panic("todo:" + f.Type.String())
	}
}
