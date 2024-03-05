package field

import (
	"fmt"
	"net"
	"slices"
	"time"

	"github.com/junhwong/goost/apm/field/loglevel"
	"github.com/spf13/cast"
)

type Field struct {
	*Schema
	*Value
	Index  int      // 索引
	Items  []*Field // 子元素
	Parent *Field   // 父元素
}

const (
	_          int32 = 1 << iota //
	ColumnFlag                   // 列, 但子类型必须完全相同. 如表格中的列
	TableFlag                    // 表格, 子元素必须完全是Column
)

// 设置类型.
func (f *Field) SetKind(t Type, isColumn, isTable bool) *Field {
	f.Type = t

	f.Flags &^= ColumnFlag
	if isColumn {
		f.Flags |= ColumnFlag
	}
	f.Flags &^= TableFlag
	if isTable {
		f.Flags |= TableFlag
	}
	return f
}

func (f *Field) setPK(t Type) {
	f.SetKind(t, false, false)
}

// 是否null
func (f *Field) IsNull() bool {
	return f == nil || f.Value == nil || f.NullValue
}

func (f *Field) SetNull(b bool) *Field {
	f.NullValue = b
	return f
}

// 是否是集合
func (f *Field) IsCollection() bool {
	return f.Type == GroupKind || f.Type == ArrayKind || f.IsColumn() || f.IsTable()
}

func (f *Field) IsArray() bool {
	return f.Type == ArrayKind || f.IsColumn()
}

// 是否是列或统一类型的数组.
func (f *Field) IsColumn() bool {
	return (f.GetFlags() & ColumnFlag) == ColumnFlag
}

// 是否是表格类型. 行:列
func (f *Field) IsTable() bool {
	return (f.GetFlags() & TableFlag) == TableFlag
}

func (f *Field) resetValue() {
	f.SetNull(false)
	f.IntValue = nil
	f.UintValue = nil
	f.FloatValue = nil
	f.StringValue = nil
	f.BytesValue = nil
	f.ItemsValue = nil
	f.Items = nil
	f.ItemsSchema = nil
}

func (f *Field) SetString(v string) *Field {
	f.resetValue()
	f.setPK(StringKind)

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
	f.resetValue()
	f.setPK(BoolKind)

	var b int64
	if v {
		b = 1
	}
	f.IntValue = &b
	return f
}

func (f *Field) SetInt(v int64) *Field {
	f.resetValue()
	f.setPK(IntKind)

	f.IntValue = &v
	return f
}

func (f *Field) GetInt() int64 {
	if !f.isKind(IntKind) {
		return 0
	}
	return f.GetIntValue()
}

func (f *Field) SetUint(v uint64) *Field {
	f.resetValue()
	f.setPK(UintKind)

	f.UintValue = &v
	return f
}

func (f *Field) GetUint() uint64 {
	if !f.isKind(UintKind) {
		return 0
	}
	return f.GetUintValue()
}

func (f *Field) SetFloat(v float64) *Field {
	f.resetValue()
	f.setPK(FloatKind)

	f.FloatValue = &v
	return f
}
func (f *Field) GetFloat() float64 {
	if !f.isKind(FloatKind) {
		return 0
	}
	return f.GetFloatValue()
}

func (f *Field) SetTime(v time.Time) *Field {
	f.resetValue()
	f.setPK(TimeKind)

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
	f.resetValue()
	f.setPK(DurationKind)

	i := int64(v)
	f.IntValue = &i
	return f
}

func (f *Field) SetIP(v net.IP) *Field {
	f.resetValue()
	f.setPK(IPKind)

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

func (f *Field) SetLevel(v loglevel.Level) *Field {
	f.resetValue()
	f.setPK(LevelKind)

	i := uint64(v)
	f.UintValue = &i

	return f
}
func (f *Field) GetLevel() loglevel.Level {
	if !f.isKind(LevelKind) {
		return loglevel.Unset
	}

	return loglevel.FromInt(int(f.GetUintValue()))
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
	f.resetValue()
	f.setPK(BytesKind)

	f.SetNull(v == nil)

	f.BytesValue = v

	return f
}

func (f *Field) isKind(t Type) bool {
	if f == nil {
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

func (f *Field) GetItem(k string) *Field {
	if f.Type != GroupKind {
		panic(fmt.Errorf("类型不匹配:必须是%v,%v", GroupKind, f.Type))
	}
	return Get(f.Items, k)
}

func (f *Field) RemoveItem(k string) (dst *Field) {
	for {
		found := f.GetItem(k)
		if found == nil {
			break
		}
		dst = found.Remove()
	}
	return
}

func (f *Field) SetArray(t Type, v []*Field, isColumn ...bool) error {
	b := false
	if len(isColumn) > 0 {
		b = isColumn[len(isColumn)-1]
	}
	f.resetValue()
	if !b {
		t = ArrayKind
	}
	f.SetKind(t, b, false)
	f.SetNull(v == nil)
	if f.IsNull() {
		return nil
	}

	for _, v2 := range v {
		if err := f.Append(v2); err != nil {
			return err
		}
	}

	return nil
}

func (f *Field) SetGroup(v []*Field, isTable ...bool) error {
	b := false
	if len(isTable) > 0 {
		b = isTable[len(isTable)-1]
	}

	// old := f.Items

	f.resetValue()
	f.SetKind(GroupKind, false, b)
	f.SetNull(v == nil)
	if f.IsNull() {
		return nil
	}

	// for _, it := range old {
	// 	f.Set(it)
	// }

	for _, it := range v {
		f.Set(it)
	}

	return nil
}

func (f *Field) Set(n *Field) {
	if n == nil {
		return
	}
	if n.Name == "" { // todo 验证名称
		panic(fmt.Errorf("元素必须包含名称"))
	}
	if f.Type != GroupKind {
		panic(fmt.Errorf("类型不匹配:必须是%v,%v", GroupKind, f.Type))
	}
	if f.IsTable() && !n.IsColumn() {
		panic(fmt.Errorf("表格元素必须是Serial类型"))
	}
	f.SetNull(false)
	n.Parent = f
	for i, s := range f.ItemsSchema {
		if s.Name == n.Name {
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
	if !(f.Type == ArrayKind || (f.IsColumn() && n.isKind(f.Type))) {
		panic(fmt.Errorf("元素的类型不匹配: %v,%v", f.Type, n.Type))
	}

	f.SetNull(false)
	n.Parent = f
	n.Index = len(f.Items)
	f.ItemsSchema = append(f.ItemsSchema, n.Schema)
	f.ItemsValue = append(f.ItemsValue, n.Value)
	f.Items = append(f.Items, n)
	return nil
}

// 从树中移除自身
func (f *Field) Remove() *Field {
	if f == nil || f.Parent == nil {
		return f
	}
	self := f
	index := self.Index
	f = self.Parent
	self.Parent = nil
	self.Index = -1

	if !(f.Type == GroupKind || f.IsCollection()) {
		return self
	}

	if f.Type == GroupKind {
		is, _, ok := RemoveAt(index, f.ItemsSchema)
		if ok {
			f.ItemsSchema = is
		}
	}

	itvs, _, ok := RemoveAt(index, f.ItemsValue)
	if ok {
		f.ItemsValue = itvs
	}

	its, _, ok := RemoveAt(index, f.Items)
	if ok {
		f.Items = its
	}
	for i, v := range f.Items {
		v.Index = i
	}

	return self
}

func (f *Field) Sort(less func(a, b *Field) int) {
	if f == nil || !f.IsCollection() {
		return
	}
	slices.SortFunc(f.Items, less)
	for i, it := range f.Items {
		if it.Index == i {
			continue
		}
		it.Index = i
		f.ItemsSchema[it.Index] = it.Schema
		f.ItemsValue[it.Index] = it.Value
	}
}

func RemoveAt[T any](i int, tmp []T) ([]T, T, bool) {
	// slices.DeleteFunc(tmp, func(t T) bool {
	// 	return true
	// })
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

func GetValue(f *Field) any {
	if f == nil || f.Type == InvalidKind {
		return nil
	}
	if f.IsNull() {
		return nil
	}
	if f.IsCollection() {
		var objs []any
		for _, f2 := range f.Items {
			if f2.Type == InvalidKind {
				continue
			}
			objs = append(objs, GetValue(f2))
		}
		return objs
	}

	switch f.Type {
	case GroupKind:
		obj := map[string]any{}
		for _, f2 := range f.Items {
			if f2.Type == InvalidKind {
				continue
			}
			obj[f2.Name] = GetValue(f2)
		}
		return obj
	default:
		return GetPrimitiveValue(f)
	}
}

func GetPrimitiveValue(f *Field) any {
	if f == nil || f.Type == InvalidKind {
		return nil
	}
	if f.IsNull() {
		return nil
	}
	if f.IsCollection() {
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
		return f.GetBytes()
	case TimeKind:
		return f.GetTime()
	case DurationKind:
		return f.GetDuration()
	case IPKind:
		return f.GetIP()
	case LevelKind:
		return f.GetLevel()
	case GroupKind:
		return nil
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
func As(f *Field, t Type, layouts []string, loc *time.Location) error {
	if f.Parent != nil && f.IsColumn() && f.Parent.Type != t {
		return fmt.Errorf("必须与父级类型一致")
	}
	if f.IsColumn() {
		f.Type = t
		for _, it := range f.Items {
			if err := As(it, t, layouts, loc); err != nil {
				return err
			}
		}
		return nil
	}
	if f.Type == t {
		if t == TimeKind && loc != nil {
			f.SetTime(f.GetTime().In(loc))
		}
		return nil
	}
	switch t {
	case IntKind:
		obj := GetPrimitiveValue(f)
		v, err := cast.ToInt64E(obj)
		if err != nil {
			return err
		}
		f.SetInt(v)
		return nil
	case UintKind:
		obj := GetPrimitiveValue(f)
		v, err := cast.ToUint64E(obj)
		if err != nil {
			return err
		}
		f.SetUint(v)
		return nil
	case FloatKind:
		obj := GetPrimitiveValue(f)
		v, err := cast.ToFloat64E(obj)
		if err != nil {
			return err
		}
		f.SetFloat(v)
		return nil
	case StringKind:
		switch f.Type {
		case StringKind:
		default:
			obj := GetPrimitiveValue(f)
			v, err := cast.ToStringE(obj)
			if err != nil {
				return err
			}
			f.SetString(v)
		}
		return nil
	case BoolKind:
		obj := GetPrimitiveValue(f)
		v, err := cast.ToBoolE(obj)
		if err != nil {
			return err
		}
		f.SetBool(v)
		return nil
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
	if f == nil {
		return nil
	}
	dst := New(f.Name)
	dst.Index = f.Index
	dst.Parent = f.Parent
	dst.Type = f.Type
	dst.Flags = f.Flags
	dst.NullValue = f.NullValue
	if f.NullValue {
		return dst
	}
	dst.IntValue = f.IntValue
	dst.UintValue = f.UintValue
	dst.FloatValue = f.FloatValue
	dst.StringValue = f.StringValue
	dst.BytesValue = f.BytesValue

	if len(f.Items) == 0 {
		return dst
	}
	dst.Items = make([]*Field, 0, len(f.Items))
	dst.ItemsSchema = make([]*Schema, 0, len(f.ItemsSchema))
	dst.ItemsValue = make([]*Value, 0, len(f.ItemsValue))

	for i, f2 := range f.Items {
		f2 := Clone(f2)
		f2.Index = i
		f2.Parent = dst
		dst.Items = append(dst.Items, f2)
		dst.ItemsSchema = append(dst.ItemsSchema, f2.Schema)
		dst.ItemsValue = append(dst.ItemsValue, f2.Value)
	}

	return dst
}

func CloneInto(src, dst *Field) *Field {
	if src == nil {
		return nil
	}
	dst.Type = src.Type
	dst.Flags = src.Flags
	dst.NullValue = src.NullValue
	if src.NullValue {
		return dst
	}
	if v := src.IntValue; v != nil {
		v := *v
		src.IntValue = &v
	}
	if v := src.UintValue; v != nil {
		v := *v
		src.UintValue = &v
	}
	if v := src.FloatValue; v != nil {
		v := *v
		src.FloatValue = &v
	}
	if v := src.IntValue; v != nil {
		v := *v
		src.IntValue = &v
	}
	if v := src.StringValue; v != nil {
		v := *v
		src.StringValue = &v
	}
	if v := src.BytesValue; len(v) != 0 {
		vCopy := make([]byte, len(v))
		copy(vCopy, v)
		src.BytesValue = vCopy
	}

	if len(src.Items) == 0 {
		return dst
	}
	dst.Items = make([]*Field, 0, len(src.Items))
	dst.ItemsSchema = make([]*Schema, 0, len(src.ItemsSchema))
	dst.ItemsValue = make([]*Value, 0, len(src.ItemsValue))

	for i, f2 := range src.Items {
		f2 := Clone(f2)
		f2.Index = i
		f2.Parent = dst
		dst.Items = append(dst.Items, f2)
		dst.ItemsSchema = append(dst.ItemsSchema, f2.Schema)
		dst.ItemsValue = append(dst.ItemsValue, f2.Value)
	}

	return dst
}

func (f *Field) GoString() string {
	var v string

	switch f.GetType() {
	case GroupKind:
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
		if f.IsCollection() {
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
		v = fmt.Sprintf("%v", GetValue(f))
	}
	return fmt.Sprintf("Field(Name:%v type: %v value: %v)", f.Name, f.Type, v)
}
