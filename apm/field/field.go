package field

import (
	"fmt"
	"net"
	"time"
)

// 字段标志位.
type Flags int32

const (
	FlagKey Flags = 1 << iota // 索引

	// RES                       // 是资源
	// ATTR                      // 是属性
	// BDG                       // 是传播
	// BDY                       // 是内容
	// SRC                       // 源
)

func IsKey(f *Field) bool { return (Flags(f.GetFlags()) & FlagKey) == FlagKey }

func (f *Field) SetString(v string) *Field {
	f.Type = StringKind
	f.StringValue = &v
	return f
}
func (f *Field) GetString() string {
	if f == nil || f.Type != StringKind {
		return ""
	}
	return f.GetStringValue()
}

func (f *Field) SetBool(v bool) *Field {
	f.Type = BoolKind
	var b int64
	if v {
		b = 1
	}
	f.IntValue = &b
	return f
}

func (f *Field) SetInt(v int64) *Field {
	f.Type = IntKind
	f.IntValue = &v
	return f
}

func (f *Field) SetUint(v uint64) *Field {
	f.Type = UintKind
	f.UintValue = &v
	return f
}

func (f *Field) SetFloat(v float64) *Field {
	f.Type = FloatKind
	f.FloatValue = &v
	return f
}

func (f *Field) SetTime(v time.Time) *Field {
	f.Type = TimeKind
	i := uint64(v.UnixNano())
	f.UintValue = &i
	return f
}
func (f *Field) GetTime() time.Time {
	if f == nil || f.Type != TimeKind {
		return time.Time{}
	}
	return time.Unix(0, int64(f.GetUintValue()))
}

func (f *Field) SetDuration(v time.Duration) *Field {
	f.Type = DurationKind
	i := int64(v)
	f.IntValue = &i
	return f
}

func (f *Field) SetIP(v net.IP) *Field {
	if l := len(v); !(l == net.IPv4len || l == net.IPv6len) {
		return f
	}
	f.Type = IPKind
	f.BytesValue = v
	return f
}
func (f *Field) GetIP() net.IP {
	if f == nil || f.Type != IPKind {
		return nil
	}
	if l := len(f.BytesValue); !(l == net.IPv4len || l == net.IPv6len) {
		return nil
	}
	return f.BytesValue
}

func (f *Field) SetLevel(v Level) *Field {
	f.Type = LevelKind
	i := uint64(v)
	f.UintValue = &i
	return f
}
func (f *Field) GetLevel() Level {
	if f == nil || f.Type != LevelKind {
		return LevelUnset
	}
	return LevelFromInt(int(f.GetUintValue()))
}

func (f *Field) GetDuration() time.Duration {
	if f == nil || f.Type != DurationKind {
		return 0
	}
	return time.Duration(f.GetIntValue())
}

func (f *Field) GetBool() bool {
	if f == nil || f.Type != BoolKind {
		return false
	}
	return f.GetIntValue() != 0
}

func (f *Field) SetBytes(v []byte) *Field {
	f.Type = BytesKind
	f.BytesValue = v
	return f
}

func (f *Field) GetBytes() []byte {
	if f == nil || f.Type != BytesKind {
		return nil
	}
	return f.GetBytesValue()
}

func (f *Field) SetMap(v FieldSet) *Field {
	f.Type = MapKind
	f.ItemsValue = v
	return f
}

func GetObject(f *Field) any {
	if f == nil {
		return nil
	}
	switch f.Type {
	case StringKind:
		return f.GetString()
	case IntKind:
		return f.GetIntValue()
	case UintKind:
		return f.GetUintValue()
	case FloatKind:
		return f.GetFloatValue()
	case BoolKind:
		return f.GetBool()
	case BytesKind:
		return f.GetBytesValue()
	case TimeKind:
		return f.GetTime()
	case DurationKind:
		return f.GetDuration()
	case IPKind:
		return f.GetIP()
	case LevelKind:
		return f.GetLevel()
	}
	return nil
}

// 从值倒推类型(只能是基本类型)
func InvertType(f *Field) *Field {
	switch f.Type {
	case IntKind:
		if f.IntValue == nil {
			f.Type = InvalidKind
		}
	case UintKind:
		if f.UintValue == nil {
			f.Type = InvalidKind
		}
	case FloatKind:
		if f.FloatValue == nil {
			f.Type = InvalidKind
		}
	case StringKind:
		if f.StringValue == nil {
			f.Type = InvalidKind
		}
	case BoolKind:
		if f.IntValue == nil {
			f.Type = InvalidKind
		}
	case BytesKind:
		if f.BytesValue == nil {
			f.Type = InvalidKind
		}
	case TimeKind:
		if f.UintValue == nil {
			f.Type = InvalidKind
		}
	case DurationKind:
		if f.IntValue == nil {
			f.Type = InvalidKind
		}
	default:
		f.Type = InvalidKind
	}
	if f.Type != InvalidKind {
		return f
	}
	switch {
	case f.BytesValue != nil:
		f.Type = BytesKind
	case f.FloatValue != nil:
		f.Type = FloatKind
	case f.IntValue != nil:
		f.Type = IntKind
	case f.UintValue != nil:
		f.Type = UintKind
	case f.StringValue != nil:
		f.Type = StringKind
	}
	return f
}

// 转换类型. 转换失败将不会改变
func As(f *Field, t Kind, layouts []string, loc *time.Location) error {
	if f.Type == t {
		if t == TimeKind && loc != nil {
			f.SetTime(f.GetTime().In(loc))
		}
		return nil
	}
	switch t {
	case IntKind:
		panic("todo convert")
	case UintKind:
		panic("todo convert")
	case FloatKind:
		panic("todo convert")
	case StringKind:
		panic("todo convert")
	case BoolKind:
		panic("todo convert")
	case BytesKind:
		switch f.Type {
		case StringKind:
			f.StringValue = nil
			f.SetBytes([]byte(f.GetStringValue()))
			return nil
		case BytesKind:
			return nil
		default:
			panic("todo convert")
		}

	case TimeKind:
		switch f.GetType() {
		case IntKind:
			panic("todo convert")
		case UintKind:
			panic("todo convert")
		case FloatKind:
			panic("todo convert")
		case StringKind: // 字符串转日期
			v, err := ParseTime(f.GetStringValue(), layouts, loc)
			if err == nil {
				f.SetTime(v)
			}
			return err
		default:
			panic("todo convert")
		}
	case DurationKind:
		panic("todo convert1")
	}

	panic(fmt.Sprintf("todo convert %v->%v", f.GetType(), t))
}

func Clone(f *Field) *Field {
	dst := New(f.GetKey())
	dst.Type = f.Type
	dst.Flags = f.Flags
	if f.ItemsValue != nil {
		dst.ItemsValue = make([]*Field, 0, len(f.ItemsValue))
		for _, f2 := range f.ItemsValue {
			dst.ItemsValue = append(dst.ItemsValue, Clone(f2))
		}
	}
	dst.BytesValue = f.BytesValue
	dst.FloatValue = f.FloatValue
	dst.IntValue = f.IntValue
	dst.UintValue = f.UintValue
	dst.StringValue = f.StringValue

	return dst
}

func (f *Field) GoString() string {
	var v string
	switch f.GetType() {
	case MapKind:
		v = "{"
		for i, f2 := range f.ItemsValue {
			if i != 0 {
				v += ","
			}
			v += fmt.Sprintf("%#v", f2)
		}
		v += "}"
	case ArrayKind:
		v = "["
		for i, f2 := range f.ItemsValue {
			if i != 0 {
				v += ","
			}
			v += fmt.Sprintf("%#v", f2)
		}
		v += "]"
	case BytesKind:
		v = "<bytes>"
	default:
		v = fmt.Sprintf("%v", GetObject(f))
	}
	return fmt.Sprintf("Field(key:%v type: %v value: %v)", f.Key, f.Type, v)
}
