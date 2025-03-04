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
	TimeLayout       string // 日期格式
	OmitEmpty        bool   // 是否忽略空值
	DurationToString bool   // 是否转字符串
	EscapeHTML       bool   // 是否转义html
	Pretty           bool   // 是否美化
	NameFilter       func(string) string
	NameSort         func(a, b *Field) int // group 排序函数

	err   error
	w     io.Writer
	n     int64
	skip  []*Field
	ident int
}

func (m JsonMarshaler) Marshal(f *Field, w io.Writer, skip ...*Field) (int64, error) {
	m.w = w
	m.err = nil
	m.n = 0
	m.skip = skip
	if m.NameSort == nil {
		m.NameSort = NameLess
	}
	m.write(f, func() {})
	return m.n, m.err
}

func (m JsonMarshaler) MarshalGroup(fs []*Field, w io.Writer, skip ...*Field) (int64, error) {
	m.w = w
	m.err = nil
	m.n = 0
	m.skip = skip
	if m.NameSort == nil {
		m.NameSort = NameLess
	}
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
	items := make([]*Field, len(fs))
	copy(items, fs)
	slices.SortFunc(items, m.NameSort)

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
			if m.Pretty {
				m.writeBytes([]byte(strings.Repeat(" ", m.ident)))
			}
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
		b := strconv.AppendFloat(nil, v, 'g', -1, 64)                                       // -1
		if bytes.LastIndex(b, []byte{'.'}) == -1 && bytes.LastIndex(b, []byte{'e'}) == -1 { // 保留数据类型
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
		data, err := json.Marshal(v)
		if err != nil {
			m.err = err
			return
		}
		if n := len(data) - 1; data[n] == '\n' { // json.NewEncoder 会增加一个换行
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
		if !m.DurationToString {
			m.writeBytes(strconv.AppendInt(nil, v.Nanoseconds(), 10))
			return
		}
		m.writeBytes([]byte(`"`))
		m.writeBytes([]byte(v.String()))
		m.writeBytes([]byte(`"`))

	case BytesKind:
		v := f.GetBytes()
		befor()
		m.writeBytes([]byte(`"base64:`))
		benc := base64.NewEncoder(base64.StdEncoding, m.w)
		_, m.err = benc.Write(v)
		if m.err != nil {
			return
		}
		m.err = benc.Close()
		m.writeBytes([]byte(`"`))
	default:
		panic("todo:" + f.Type.String())
	}
}

func NameLess(a, b *Field) int {
	if a.Name == b.Name {
		return 0
	}
	if a.Name < b.Name {
		return -1
	}
	return 1
}
