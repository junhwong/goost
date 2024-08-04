package field

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"io"
	"slices"
	"strconv"
	"strings"
	"time"
)

const (
	DefaultTimeLayout = time.RFC3339Nano
)

type JsonMarshaler struct {
	OmitEmpty  bool // todo 实现
	TimeLayout string
	EscapeHTML bool
	Pretty     bool
	NameFilter func(string) string
	NameLess   func(a, b *Field) int

	err   error
	w     io.Writer
	n     int64
	skip  []*Field
	ident int
	// buf  *bytes.Buffer
	// enc  *json.Encoder
	// benc io.WriteCloser
}

func (m JsonMarshaler) Marshal(f *Field, w io.Writer, skip ...*Field) (int64, error) {
	m.w = w
	m.err = nil
	m.n = 0
	m.skip = skip
	m.write(f, func() {})
	return m.n, m.err
}

func (m JsonMarshaler) MarshalGroup(fs []*Field, w io.Writer, skip ...*Field) (int64, error) {
	m.w = w
	m.err = nil
	m.n = 0
	m.skip = skip
	// m.enc = json.NewEncoder(m.buf)
	m.writeGroup(fs)
	return m.n, m.err
}

func (m *JsonMarshaler) writeByte(c byte) {
	if m.err != nil {
		return
	}
	n, err := m.w.Write([]byte{c})
	m.err = err
	m.n += int64(n)
}
func (m *JsonMarshaler) writeBytes(p []byte) {
	if m.err != nil {
		return
	}
	n, err := m.w.Write(p)
	m.err = err
	m.n += int64(n)
}

func (m *JsonMarshaler) writeGroup(fs []*Field) {
	if m.err != nil {
		return
	}
	items := fs
	if m.NameLess != nil {
		items = make([]*Field, len(items))
		copy(items, fs)
		slices.SortFunc(items, m.NameLess)
	}

	has := false
	m.writeByte('{')
	if m.Pretty {
		m.writeByte('\n')
	}
	m.ident += 2

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

		m.write(it, func() {
			if has {
				m.writeByte(',')
				if m.Pretty {
					m.writeByte('\n')
				}
			}
			if m.Pretty {
				m.writeBytes([]byte(strings.Repeat(" ", m.ident)))
			}
			has = true
			m.writeBytes([]byte(`"`))
			m.writeBytes([]byte(name))
			m.writeBytes([]byte(`":`))
			if m.Pretty {
				m.writeByte(' ')
			}
		})
	}
	m.ident -= 2
	if m.Pretty {
		m.writeByte('\n')
		m.writeBytes([]byte(strings.Repeat(" ", m.ident)))
	}
	m.writeByte('}')
}

func (m *JsonMarshaler) writeArray(fs []*Field) {
	if m.err != nil {
		return
	}
	has := false
	m.writeByte('[')
	if m.Pretty {
		m.writeByte('\n')
	}
	m.ident += 2
	for _, it := range fs {
		if it.Type == InvalidKind {
			continue
		}

		m.write(it, func() {
			if has {
				m.writeByte(',')
				if m.Pretty {
					m.writeByte('\n')
				}
			}
			has = true
			m.writeBytes([]byte(strings.Repeat(" ", m.ident)))
		})
	}
	m.ident -= 2
	if m.Pretty {
		m.writeByte('\n')
		m.writeBytes([]byte(strings.Repeat(" ", m.ident)))
	}
	m.writeByte(']')
}

func (m *JsonMarshaler) write(f *Field, befor func()) {
	if m.err != nil {
		return
	}

	if f.Type == InvalidKind {
		return
	}

	for _, it := range m.skip {
		if it == f {
			return
		}
	}

	if f.IsArray() {
		if f.IsNull() || len(f.Items) == 0 {
			if !m.OmitEmpty {
				befor()
				m.writeBytes([]byte("null"))
			}
			return
		}
		befor()
		m.writeArray(f.Items)
		return
	}

	if f.IsGroup() {
		if f.IsNull() || len(f.Items) == 0 {
			if !m.OmitEmpty {
				befor()
				m.writeBytes([]byte("null"))
			}
			return
		}
		befor()
		m.writeGroup(f.Items)
		return
	}

	switch f.Type {
	case IntKind:
		v := f.GetInt()
		befor()
		m.writeBytes(strconv.AppendInt(nil, v, 10))
	case UintKind:
		v := f.GetUint()
		befor()
		m.writeBytes(strconv.AppendUint(nil, v, 10))
	case FloatKind:
		v := f.GetFloat()
		b := strconv.AppendFloat(nil, v, 'f', -1, 64) // -1
		if bytes.LastIndex(b, []byte{'.'}) == -1 {    // 保留数据类型
			b = append(b, '.', '0')
		}
		befor()
		m.writeBytes(b)
	case BoolKind:
		v := f.GetBool()
		befor()
		m.writeBytes(strconv.AppendBool(nil, v))
	case StringKind:
		v := f.GetString()
		if v == "" {
			if !m.OmitEmpty {
				befor()
				m.writeBytes([]byte(`""`))
			}
			return
		}
		// if m.buf == nil {
		// 	m.buf = bytes.NewBuffer(nil)
		// }
		// if m.enc == nil {
		// 	m.enc = json.NewEncoder(m.buf)
		// }

		// m.enc.SetEscapeHTML(m.EscapeHTML)
		// m.buf.Reset()

		// if err := m.enc.Encode(v); err != nil {
		// 	m.err = err
		// 	return
		// }

		// var data []byte
		// data = m.buf.Bytes()
		data, err := json.Marshal(v)
		if err != nil {
			m.err = err
			return
		}
		n := len(data) - 1
		if data[n] == '\n' { // json.NewEncoder 会增加一个换行
			data = data[:n]
		}
		befor()
		m.writeBytes(data)

	case TimeKind:
		v := f.GetTime()
		befor()
		m.writeBytes([]byte(`"`))
		if m.TimeLayout != "" {
			m.writeBytes([]byte(v.Format(m.TimeLayout)))
		} else {
			m.writeBytes([]byte(v.Format(DefaultTimeLayout)))
		}
		m.writeBytes([]byte(`"`))
	case IPKind:
		v := f.GetIP()
		befor()
		m.writeBytes([]byte(`"`))
		m.writeBytes([]byte(v.String()))
		m.writeBytes([]byte(`"`))
	case LevelKind:
		v := f.GetLevel()
		befor()
		m.writeBytes([]byte(`"`))
		m.writeBytes([]byte(v.String()))
		m.writeBytes([]byte(`"`))
	case DurationKind:
		v := f.GetDuration()
		befor()
		m.writeBytes(strconv.AppendInt(nil, v.Nanoseconds(), 10))
	case BytesKind:
		v := f.GetBytes()
		befor()
		m.writeBytes([]byte(`"base64:`))
		// if m.buf == nil {
		// 	m.buf = bytes.NewBuffer(nil)
		// }
		// if m.benc == nil {
		// 	m.benc = base64.NewEncoder(base64.StdEncoding, m.buf)
		// }
		// m.buf.Reset()
		// _, m.err = m.benc.Write(v)
		// m.writeBytes(m.buf.Bytes())

		buf := bytes.NewBuffer(nil)
		benc := base64.NewEncoder(base64.StdEncoding, buf)
		buf.Reset()
		_, m.err = benc.Write(v)
		if m.err != nil {
			return
		}
		data := buf.Bytes()

		m.writeBytes(data)
		m.writeBytes([]byte(`"`))
	default:
		panic("todo:" + f.Type.String())
	}
}
