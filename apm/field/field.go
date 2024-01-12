package field

import (
	"fmt"
	"net"
	"time"
)

type Field struct {
	*Schema
	*Value
	Index  int
	Items  FieldSet
	Parent *Field
}

const (
	_         int32 = 1 << iota //
	NullFlag                    // 空值
	ListFlag                    // 列表
	TableFlag                   // 表格
	// RES                       // 是资源
	// ATTR                      // 是属性
	// BDG                       // 是传播
	// BDY                       // 是内容
	// SRC                       // 源
)

// func IsKey(f *Field) bool { return (Flags(f.GetFlags()) & FlagKey) == FlagKey }

// 设置类型.
func (f *Field) SetKind(t Kind, isList bool) *Field {
	f.Type = t
	f.Flags &^= ListFlag
	if isList {
		f.Flags |= ListFlag
	}
	if t != MapKind {
		f.ItemsSchema = nil
	}
	return f
}

// 是否是列表
func (f *Field) IsNull() bool {
	return (f.GetFlags() & NullFlag) == NullFlag
}

func (f *Field) SetNull(b bool) *Field {
	f.Flags &^= NullFlag
	if b {
		f.Flags |= NullFlag
	}
	return f
}

// 是否是列表
func (f *Field) IsList() bool {
	return (f.GetFlags() & ListFlag) == ListFlag
}

func (f *Field) reset() {
	f.SetNull(false)
	f.IntValue = nil
	f.UintValue = nil
	f.FloatValue = nil
	f.StringValue = nil
	f.BytesValue = nil
	f.ItemsValue = nil
	f.Items = nil
}

func (f *Field) SetString(v string) *Field {
	f.SetKind(StringKind, false)
	f.reset()

	f.StringValue = &v
	return f
}
func (f *Field) GetString() string {
	if !f.isKind(StringKind) {
		return ""
	}
	return f.GetStringValue()
}

func (f *Field) SetBool(v bool) *Field {
	f.SetKind(BoolKind, false)
	f.reset()

	var b int64
	if v {
		b = 1
	}
	f.IntValue = &b
	return f
}

func (f *Field) SetInt(v int64) *Field {
	f.SetKind(IntKind, false)
	f.reset()

	f.IntValue = &v
	return f
}

func (f *Field) SetUint(v uint64) *Field {
	f.SetKind(UintKind, false)
	f.reset()

	f.UintValue = &v
	return f
}

func (f *Field) SetFloat(v float64) *Field {
	f.SetKind(FloatKind, false)
	f.reset()
	f.FloatValue = &v
	return f
}

func (f *Field) SetTime(v time.Time) *Field {
	f.SetKind(TimeKind, false)
	f.reset()
	f.SetNull(v.IsZero())

	i := uint64(v.UnixNano())
	f.UintValue = &i

	return f
}
func (f *Field) GetTime() time.Time {
	if !f.isKind(TimeKind) {
		return time.Time{}
	}

	return time.Unix(0, int64(f.GetUintValue()))
}

func (f *Field) SetDuration(v time.Duration) *Field {
	f.SetKind(DurationKind, false)
	f.reset()
	i := int64(v)
	f.IntValue = &i
	return f
}

func (f *Field) SetIP(v net.IP) *Field {
	f.SetKind(IPKind, false)
	f.reset()

	l := len(v)
	f.SetNull(!(l == net.IPv4len || l == net.IPv6len))

	f.BytesValue = v

	return f
}
func (f *Field) GetIP() net.IP {
	if !f.isKind(IPKind) {
		return nil
	}
	if l := len(f.BytesValue); !(l == net.IPv4len || l == net.IPv6len) {
		return nil
	}
	return f.BytesValue
}

func (f *Field) SetLevel(v Level) *Field {
	f.SetKind(LevelKind, false)
	f.reset()

	i := uint64(v)
	f.UintValue = &i

	return f
}
func (f *Field) GetLevel() Level {
	if !f.isKind(LevelKind) {
		return LevelUnset
	}

	return LevelFromInt(int(f.GetUintValue()))
}

func (f *Field) GetDuration() time.Duration {
	if !f.isKind(DurationKind) {
		return 0
	}
	return time.Duration(f.GetIntValue())
}

func (f *Field) GetBool() bool {
	if !f.isKind(BoolKind) {
		return false
	}
	return f.GetIntValue() != 0
}

func (f *Field) SetBytes(v []byte) *Field {
	f.SetKind(BytesKind, false)
	f.reset()
	f.SetNull(v == nil)

	f.BytesValue = v

	return f
}

func (f *Field) isKind(t Kind) bool {
	if f == nil || f.IsList() { // todo || f.IsNull()
		return false
	}
	return f.Type == t
}

func (f *Field) GetBytes() []byte {
	if !f.isKind(BytesKind) {
		return nil
	}
	return f.GetBytesValue()
}

func (f *Field) SetList(t Kind, v []*Field) error {
	f.SetKind(t, true)
	f.reset()
	f.SetNull(v == nil)

	for i, v2 := range v {
		// v2.Schema.ParentSchema = f.Schema
		// v2.Value.ParentValue = f.Value
		if v2.Type != t {
			return fmt.Errorf("元素%d的类型不匹配: %v,%v", i, t, v2.Type)
		}
		f.Append(v2)
	}

	return nil
}

func (f *Field) SetNest(v []*Field, isTable bool) *Field {
	f.SetKind(MapKind, false)
	f.reset()
	f.SetNull(v == nil)

	for _, v2 := range v {
		f.Set(v2)
	}

	return f
}

func (f *Field) Set(n *Field) {
	if !f.isKind(MapKind) {
		panic(fmt.Errorf("标识re: %v,%v", f.Type, n.Type))
	}
	f.SetNull(false)
	n.Parent = f
	for i, s := range f.ItemsSchema {
		if s.Key == n.Key {
			n.Index = i
			f.ItemsSchema[i] = n.Schema
			f.ItemsValue[i] = n.Value
			f.Items[i] = n
			return
		}
	}
	n.Index = len(f.Items)
	f.ItemsSchema = append(f.ItemsSchema, n.Schema)
	f.ItemsValue = append(f.ItemsValue, n.Value)
	f.Items = append(f.Items, n)
}

func (f *Field) Append(n *Field) error {
	if !f.IsList() || !n.isKind(f.Type) {
		panic(fmt.Errorf("元素的类型不匹配: %v,%v", f.Type, n.Type))
	}
	f.SetNull(false)
	n.Parent = f
	n.Index = len(f.Items)
	f.ItemsValue = append(f.ItemsValue, n.Value)
	f.Items = append(f.Items, n)
	return nil
}

func (f *Field) Remove() *Field {
	if f == nil || f.Parent == nil {
		return f
	}
	f = f.Parent
	if !(f.Type == MapKind || f.IsList()) {
		return f
	}

	if f.Type == MapKind {
		is, _, ok := RemoveAt(f.Index, f.ItemsSchema)
		if ok {
			f.ItemsSchema = is
		}
	}

	itvs, _, ok := RemoveAt(f.Index, f.ItemsValue)
	if ok {
		f.ItemsValue = itvs
	}

	its, found, ok := RemoveAt(f.Index, f.Items)
	if ok {
		f.Items = its
	}
	for i, v := range f.Items {
		v.Index = i
	}
	return found
}

func RemoveAt[T any](i int, tmp []T) ([]T, T, bool) {
	t := len(tmp)
	var found T

	if i < 0 || i >= t {
		return tmp, found, false
	}
	found = tmp[i]
	for j := i; j < t-1; j++ { // 将后面的元素提前
		tmp[j] = tmp[j+1]
	}
	return tmp[:t-1], found, true
}

func GetObject(f *Field) any {
	if f == nil || f.Type == InvalidKind {
		return nil
	}
	if f.IsNull() {
		return nil
	}
	if f.IsList() {
		var objs []any
		for _, f2 := range f.Items {
			if f2.Type == InvalidKind {
				continue
			}
			objs = append(objs, GetObject(f2))
		}
		return objs
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
		return f.GetBytes()
	case TimeKind:
		return f.GetTime()
	case DurationKind:
		return f.GetDuration()
	case IPKind:
		return f.GetIP()
	case LevelKind:
		return f.GetLevel()
	case MapKind:
		obj := map[string]any{}
		for _, f2 := range f.Items {
			if f2.Type == InvalidKind {
				continue
			}
			obj[f2.GetKey()] = GetObject(f2)
		}
		return obj
	default:
		panic(fmt.Sprintf("未定义:%#v", f))
	}
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
			v, err := ParseTime(f.GetString(), layouts, loc)
			if err == nil {
				f.SetTime(v)
			} else {
				f.SetNull(true)
				fmt.Printf("field:转换为时间戳失败: %v\n", err)
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
	// sch := *f.Schema
	// dst.Schema = &sch
	// if f.Schema.ItemsValue != nil {
	// 	ItemsValueCopy := make([]*Schema, len(f.Schema.ItemsValue))
	// 	copy(ItemsValueCopy, f.Schema.ItemsValue)
	// 	dst.Schema.ItemsValue = ItemsValueCopy

	// }

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
		v = "{\n"
		for i, f2 := range f.Items {
			if i != 0 {
				v += ",\n\n"
			}
			v += fmt.Sprintf("\t%#v", f2)
		}
		v += "\n}"
	case BytesKind:
		v = "<bytes>"
	default:
		if f.IsList() {
			v = "[\n"
			for i, f2 := range f.Items {
				if i != 0 {
					v += ",\n\n"
				}
				v += fmt.Sprintf("\t%#v", f2)
			}
			v += "\n]"
			break
		}
		v = fmt.Sprintf("%v", GetObject(f))
	}
	return fmt.Sprintf("Field(key:%v type: %v value: %v)", f.Key, f.Type, v)
}
