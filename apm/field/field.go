package field

import (
	"fmt"
	"net"
	"sort"
	"time"

	"github.com/junhwong/goost/runtime"
)

type FieldSet []*Field

func (x FieldSet) Len() int           { return len(x) }
func (x FieldSet) Less(i, j int) bool { return x[i].GetKey() < x[j].GetKey() } // 字典序
func (x FieldSet) Swap(i, j int)      { x[i], x[j] = x[j], x[i] }
func (x FieldSet) Sort() {
	if len(x) == 0 {
		return
	}
	sort.Sort(x)
}

func (fs *FieldSet) Set(f *Field) *Field {
	f, _ = fs.Put(f)
	return f
}
func (fs *FieldSet) Put(f *Field) (crt, old *Field) {
	crt = f
	i := fs.At(f.GetKey())
	if i < 0 {
		*fs = append(*fs, f)
		return
	}
	tmp := *fs
	old = tmp[i]
	tmp[i] = f
	return
}

func (fs FieldSet) Get(k string) *Field {
	for _, v := range fs {
		if v.GetKey() == k {
			return v
		}
	}
	return nil
}
func (fs FieldSet) At(k string) int {
	for i, v := range fs {
		if v.GetKey() == k {
			return i
		}
	}
	return -1
}

func (fs *FieldSet) Remove(k string) *Field {
	i := fs.At(k)
	if i < 0 {
		return nil
	}

	tmp := *fs
	f := tmp[i]
	n := len(tmp) - 1
	for j := i; j < n; j++ {
		tmp[j] = tmp[j+1]
	}
	*fs = tmp[:n]
	return f
}

// 清除重复
func (fs *FieldSet) Unique() FieldSet {
	if fs == nil {
		return nil
	}
	tmp := FieldSet{}
	for _, f := range *fs {
		tmp.Set(f)
	}
	*fs = tmp
	return tmp
}

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

type WrapField struct {
	*Field `json:",inline"`
	Value  any      `json:"value"`
	As     string   `json:"as"`
	Zone   string   `json:"zone"`   // 时区
	Layout []string `json:"layout"` // 时间格式
}

func (f *WrapField) SetString(v string) *WrapField {
	SetString(f.Field, v)
	return f
}
func (f *WrapField) SetBool(v bool) *WrapField {
	SetBool(f.Field, v)
	return f
}

func (f *WrapField) SetInt(v int64) *WrapField {
	SetInt(f.Field, v)
	return f
}
func (f *WrapField) SetUint(v uint64) *WrapField {
	SetUint(f.Field, v)
	return f
}

func (f *WrapField) SetFloat(v float64) *WrapField {
	SetFloat(f.Field, v)
	return f
}
func (f *WrapField) SetTime(v time.Time) *WrapField {
	SetTime(f.Field, v)
	return f
}

func (f WrapField) GetTimeValue() time.Time {
	return GetTimeValue(f.Field)
}

func (f *WrapField) SetDuration(v time.Duration) *WrapField {
	SetDuration(f.Field, v)
	return f
}

func (f WrapField) GetDurationValue() time.Duration {
	return GetDurationValue(f.Field)
}

func SetString(f *Field, v string) *Field {
	f.Type = StringKind
	f.StringValue = &v
	return f
}

func SetBool(f *Field, v bool) *Field {
	f.Type = BoolKind
	var b int64
	if v {
		b = 1
	}
	f.IntValue = &b
	return f
}

func SetInt(f *Field, v int64) *Field {
	f.Type = IntKind
	f.IntValue = &v
	return f
}

func SetUint(f *Field, v uint64) *Field {
	f.Type = UintKind
	f.UintValue = &v
	return f
}

func SetFloat(f *Field, v float64) *Field {
	f.Type = FloatKind
	f.FloatValue = &v
	return f
}

func SetTime(f *Field, v time.Time) *Field {
	f.Type = TimeKind
	i := uint64(v.UnixNano())
	f.UintValue = &i
	return f
}
func GetTimeValue(f *Field) time.Time {
	if f == nil || f.Type != TimeKind {
		return time.Time{}
	}
	return time.Unix(0, int64(f.GetUintValue()))
}

func SetDuration(f *Field, v time.Duration) *Field {
	f.Type = DurationKind
	i := int64(v)
	f.IntValue = &i
	return f
}

func SetIP(f *Field, v net.IP) *Field {
	if l := len(v); !(l == net.IPv4len || l == net.IPv6len) {
		return f
	}
	f.Type = IPKind
	f.BytesValue = v
	return f
}
func GetIPValue(f *Field) net.IP {
	if f == nil || f.Type != IPKind {
		return nil
	}
	if l := len(f.BytesValue); !(l == net.IPv4len || l == net.IPv6len) {
		return nil
	}
	return f.BytesValue
}

func SetLevel(f *Field, v Level) *Field {
	f.Type = LevelKind
	i := uint64(v)
	f.UintValue = &i
	return f
}
func GetLevelValue(f *Field) Level {
	if f == nil || f.Type != LevelKind {
		return LevelUnset
	}
	return LevelFromInt(int(f.GetUintValue()))
}

func GetDurationValue(f *Field) time.Duration {
	if f == nil || f.Type != DurationKind {
		return 0
	}
	return time.Duration(f.GetIntValue())
}

func GetBoolValue(f *Field) bool {
	if f == nil || f.Type != BoolKind {
		return false
	}
	return f.GetIntValue() != 0
}

func SetBytes(f *Field, v []byte) *Field {
	f.Type = BytesKind
	f.BytesValue = v
	return f
}

func GetObject(f *Field) any {
	if f == nil {
		return nil
	}
	switch f.Type {
	case StringKind:
		return f.GetStringValue()
	case IntKind:
		return f.GetIntValue()
	case UintKind:
		return f.GetUintValue()
	case FloatKind:
		return f.GetFloatValue()
	case BoolKind:
		return f.GetIntValue() != 0
	case BytesKind:
		return f.GetBytesValue()
	case TimeKind:
		return GetTimeValue(f)
	case DurationKind:
		return GetDurationValue(f)
	case IPKind:
		return GetIPValue(f)
	case LevelKind:
		return GetLevelValue(f)

	}
	return nil
}

func (f *WrapField) String() string {
	return fmt.Sprintf("structField{%v, Value=%v, kind=%v}", f.Key, GetObject(f.Field), f.Type)
}

func (f *WrapField) Valid() bool {
	return f != nil && f.Type != InvalidKind
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
func As(f *Field, t Kind, layouts ...string) error {
	if f.Type == t {
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
			SetBytes(f, []byte(f.GetStringValue()))
			return nil
		case BytesKind:
			return nil
		default:
			panic("todo convert")
		}

	case TimeKind:
		if len(layouts) == 0 {
			return fmt.Errorf("unable convert to timestamp: %s", "layouts is required")
		}
		switch f.GetType() {
		case IntKind:
			panic("todo convert")
		case UintKind:
			panic("todo convert")
		case FloatKind:
			panic("todo convert")
		case StringKind:
			// 字符串转日期
			s := f.GetStringValue()
			for _, l := range layouts {
				v, err := time.Parse(l, s)
				if err != nil {
					fmt.Printf("err: %v\n", err)
					runtime.Debug(err)
					continue
				}
				f.StringValue = nil
				SetTime(f, v)
				return nil
			}
			// 转换失败
			return fmt.Errorf("unable convert to timestamp: %q", s)
		case TimeKind:
			return nil
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
	dst.BytesValue = f.BytesValue
	dst.Flags = f.Flags
	dst.FloatValue = f.FloatValue
	dst.IntValue = f.IntValue
	dst.StringValue = f.StringValue
	dst.Type = f.Type
	dst.UintValue = f.UintValue
	return dst
}
