package field

import (
	"fmt"
	"net"
	"slices"
	stdtime "time"

	"github.com/junhwong/goost/apm/field/loglevel"
	"github.com/junhwong/goost/apm/field/times"
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

func (f *Field) SetName(n string) {
	f.Name = n
	if f.IsArray() {
		for _, it := range f.Items {
			it.SetName(n)
		}
	}
}

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

// 设置简单类型
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

// 是否是集合(array或group)
func (f *Field) IsCollection() bool {
	return f.IsArray() || f.IsGroup()
}

// 是否是数组(array或column)
func (f *Field) IsArray() bool {
	if f == nil {
		return false
	}
	// todo 解决类型冲突
	if f.IsGroup() {
		return false
	}
	return f.IsColumn() || f.Type == ArrayKind
}

// 是否是列(统一类型的数组).
func (f *Field) IsColumn() bool {
	if f == nil {
		return false
	}
	return (f.GetFlags() & ColumnFlag) == ColumnFlag
}

// 是否是表格类型. 行:列
func (f *Field) IsTable() bool {
	if f == nil {
		return false
	}
	return (f.GetFlags() & TableFlag) == TableFlag
}

// 是否是字典类型.
func (f *Field) IsGroup() bool {
	if f == nil {
		return false
	}
	return f.Type == GroupKind && !f.IsColumn()
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

func (f *Field) GetBool() bool {
	if !f.isKind(BoolKind) {
		return false
	}
	return f.GetIntValue() != 0
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

func (f *Field) SetFloat(v float64) *Field { // todo NaN、Infinity、-Infinity
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

func (f *Field) SetTime(v stdtime.Time) *Field {
	f.resetValue()
	f.setPK(TimeKind)

	f.SetNull(v.IsZero())

	// 1970
	// 1678-2262
	if v.Year() < 1678 { // 处理时间范围
		v = v.AddDate(1678-v.Year(), 0, 0)
	}
	if v.Year() >= 2262 { // 大于范围
		v = stdtime.Time{}
	}
	i := v.UnixNano()
	f.IntValue = &i

	return f
}
func (f *Field) GetTime() stdtime.Time {
	if f.IsNull() || !f.isKind(TimeKind) {
		return stdtime.Time{}
	}
	v := stdtime.Unix(0, f.GetIntValue())
	if v.Year() <= 1678 { // 处理时间范围
		v = v.AddDate(-v.Year(), 0, 0)
	}
	if v.Year() >= 2262 { // 大于范围
		return stdtime.Time{}
	}
	v = v.In(times.Local)

	return v
}

func (f *Field) SetDuration(v stdtime.Duration) *Field {
	f.resetValue()
	f.setPK(DurationKind)

	i := int64(v)
	f.IntValue = &i
	return f
}

func (f *Field) GetDuration() stdtime.Duration {
	if !f.isKind(DurationKind) {
		return 0
	}
	return stdtime.Duration(f.GetIntValue())
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

func (f *Field) SetBytes(v []byte) *Field {
	f.resetValue()
	f.setPK(BytesKind)

	f.SetNull(v == nil)

	f.BytesValue = v

	return f
}

func (f *Field) GetBytes() []byte {
	if !f.isKind(BytesKind) {
		return nil
	}
	return f.GetBytesValue()
}

func (f *Field) isKind(t Type) bool {
	if f == nil {
		return false
	}
	return f.Type == t
}

func (f *Field) GetItem(k string) *Field {
	if !f.IsGroup() {
		panic(fmt.Errorf("类型不匹配:必须是%v,%v", GroupKind, f.Type))
	}
	return GetLast(f.Items, k)
}

func (f *Field) RemoveItem(k string) (dst *Field) {
	if !f.IsGroup() {
		panic(fmt.Errorf("类型不匹配:必须是%v,%v", GroupKind, f.Type))
	}
	for _, it := range Get(f.Items, k) {
		dst = it.Remove()
	}
	return
}

func (f *Field) SetArray(v []*Field, isColumn ...bool) *Field {
	f.resetValue()
	b := false
	if len(isColumn) > 0 {
		b = isColumn[len(isColumn)-1]
	}
	f.SetKind(ArrayKind, b, false)
	f.SetNull(v == nil)
	if f.IsNull() {
		return f
	}
	if b {
		f.Type = v[0].Type
	}

	for _, it := range v {
		f.Append(it)
	}
	return f
}

// 将元素添加到数组末尾
func (f *Field) Append(n *Field) {
	if !f.IsArray() {
		panic(fmt.Errorf("元素的类型不是数组: %v", f.Type))
	}

	if !f.IsNull() && !(!f.IsColumn() || n.isKind(f.Type)) {
		panic(fmt.Errorf("元素的类型是列,但与列的类型不匹配: %v,%v,%v", f.Type, n.Type, f.IsColumn()))
	}
	if f.IsNull() && f.IsColumn() { // 补齐类型
		f.Type = n.Type
	}
	f.SetNull(false)
	n.Parent = f
	n.Index = len(f.Items)
	f.ItemsSchema = append(f.ItemsSchema, n.Schema)
	f.ItemsValue = append(f.ItemsValue, n.Value)
	f.Items = append(f.Items, n)
}

func (f *Field) SetGroup(v []*Field, isTable ...bool) *Field {
	b := false
	if len(isTable) > 0 {
		b = isTable[len(isTable)-1]
	}
	f.resetValue()
	f.SetKind(GroupKind, false, b)
	f.SetNull(v == nil) //len(v) == 0
	if f.IsNull() {
		return f
	}

	for _, it := range v {
		f.Set(it)
	}
	return f
}

func (f *Field) Set(n *Field) {
	if f == nil {
		panic(fmt.Errorf("f connat be nil"))
	}
	if n == nil {
		panic(fmt.Errorf("n connat be nil"))
	}
	if !f.IsGroup() {
		panic(fmt.Errorf("类型不匹配:必须是%v,%v %s", GroupKind, f.Type, f.Name))
	}
	f.SetNull(false)
	n.Parent = f
	for i, s := range f.Items {
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

// 从树中移除自身
func (f *Field) Remove() *Field {
	if f == nil || f.Parent == nil {
		return f
	}
	self := f
	index, f := self.Index, self.Parent
	self.Parent = nil
	self.Index = -1

	if !f.IsCollection() {
		return self
	}

	iss, _, ok := RemoveAt(index, f.ItemsSchema)
	if ok {
		f.ItemsSchema = iss
	}

	ivs, _, ok := RemoveAt(index, f.ItemsValue)
	if ok {
		f.ItemsValue = ivs
	}

	is, _, ok := RemoveAt(index, f.Items)
	if ok {
		f.Items = is
	}

	for i, v := range f.Items { // 重新设置索引
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
	if f.IsArray() { // todo 构造具体类型,如 []string
		var objs []any
		for _, it := range f.Items {
			if it.Type == InvalidKind {
				continue
			}
			objs = append(objs, GetValue(it))
		}
		return objs
	}

	if f.IsGroup() {
		obj := map[string]any{}
		for _, it := range f.Items {
			if it.Type == InvalidKind {
				continue
			}
			obj[it.Name] = GetValue(it)
		}
		return obj
	}

	return GetPrimitiveValue(f)
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

func GetNumberValue(f *Field) any {
	if f == nil || f.Type == InvalidKind {
		return nil
	}
	switch f.Type {
	case IntKind:
		return f.GetIntValue()
	case UintKind:
		return f.GetUintValue()
	case FloatKind:
		return f.GetFloatValue()
	}
	return nil
}

func IsNumber(f *Field) bool {
	if f == nil || f.Type == InvalidKind {
		return false
	}
	switch f.Type {
	case IntKind:
		return true
	case UintKind:
		return true
	case FloatKind:
		return true
	}
	return false
}

// 克隆对象
func Clone(f *Field) *Field {
	if f == nil {
		return nil
	}
	dst := Make(f.Name)
	dst.Index = f.Index
	dst.Parent = f.Parent
	return CloneInto(f, dst)
}

// 克隆类容,不改变层级
func CloneInto(src, dst *Field) *Field {
	if src == nil {
		return dst
	}
	dst.Type = src.Type
	dst.Flags = src.Flags
	dst.NullValue = src.NullValue
	if dst.NullValue {
		return dst
	}
	if v := src.IntValue; v != nil {
		v := *v
		dst.IntValue = &v
	}
	if v := src.UintValue; v != nil {
		v := *v
		dst.UintValue = &v
	}
	if v := src.FloatValue; v != nil {
		v := *v
		dst.FloatValue = &v
	}
	if v := src.IntValue; v != nil {
		v := *v
		dst.IntValue = &v
	}
	if v := src.StringValue; v != nil {
		v := *v
		dst.StringValue = &v
	}
	if v := src.BytesValue; len(v) != 0 {
		vCopy := make([]byte, len(v))
		copy(vCopy, v)
		dst.BytesValue = vCopy
	}

	if len(src.Items) == 0 {
		return dst
	}

	dst.Items = make([]*Field, 0, len(src.Items))
	dst.ItemsSchema = make([]*Schema, 0, len(src.Items))
	dst.ItemsValue = make([]*Value, 0, len(src.Items))

	for i, it := range src.Items {
		it := Clone(it)
		it.Index = i
		it.Parent = dst
		dst.Items = append(dst.Items, it)
		dst.ItemsSchema = append(dst.ItemsSchema, it.Schema)
		dst.ItemsValue = append(dst.ItemsValue, it.Value)
	}
	return dst
}

func (f *Field) GoString() string {
	return toString(f, 0)
}
