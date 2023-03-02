package field

import (
	"fmt"
	"time"
)

// Field 表示一个标准字段。
type Field interface {
	Key() Key
	Kind() KeyKind
	Unwrap() (Key, interface{}) // 获取字段的键和值。返回 nil表示该字段无效
	GetObject() any
}

// Field 表示一个标准字段。
//
// 参考：
//
// [opentelemetry](https://opentelemetry.io/docs/reference/specification/logs/overview/)
//
// [ecs](https://github.com/elastic/ecs)
type structField struct {
	KeyField Key         `json:"key"`
	Value    interface{} `json:"value"`
	kind     KeyKind
	valid    bool // 防止自定义
	// sliceDataType KeyKind
	err error

	// ValueString *string
	// ValueInt    *int64
	// ValueUint   *uint64
	// ValueFloat  *float64
	// ValueBool   *bool

	// ValueSlice []*structField
}

func (f *structField) SetString(v string) *structField {
	if f.err != nil {
		return f
	}
	if len(v) == 0 {
		return f
	}
	f.kind = StringKind
	f.Value = v
	return f
}
func (f *structField) GetString() string {
	if f.err != nil || f.kind != StringKind || f.Value == nil {
		return ""
	}
	return f.Value.(string)
}
func (f *structField) SetBool(v bool) *structField {
	if f.err != nil {
		return f
	}
	f.kind = BoolKind
	f.Value = v
	return f
}
func (f *structField) GetBool() bool {
	if f.err != nil || f.kind != BoolKind || f.Value == nil {
		return false
	}
	return f.Value.(bool)
}
func (f *structField) SetInt(v int64) *structField {
	if f.err != nil {
		return f
	}
	f.kind = IntKind
	f.Value = v
	return f
}
func (f *structField) GetInt() int64 {
	if f.err != nil || f.kind != IntKind || f.Value == nil {
		return 0
	}
	return f.Value.(int64)
}

func (f *structField) SetUint(v uint64) *structField {
	if f.err != nil {
		return f
	}
	f.kind = UintKind
	f.Value = v
	return f
}
func (f *structField) GetUint() uint64 {
	if f.err != nil || f.kind != UintKind || f.Value == nil {
		return 0
	}
	return f.Value.(uint64)
}

func (f *structField) SetFloat(v float64) *structField {
	if f.err != nil {
		return f
	}
	f.kind = FloatKind
	f.Value = v
	return f
}
func (f *structField) GetFloat() float64 {
	if f.err != nil || f.kind != FloatKind || f.Value == nil {
		return 0
	}
	return f.Value.(float64)
}

func (f *structField) SetTime(v time.Time) *structField {
	if f.err != nil {
		return f
	}
	if v.IsZero() {
		return f
	}
	f.kind = TimeKind

	f.Value = v
	return f
}

func (f *structField) GetTime() time.Time {
	if f.err != nil || f.kind != TimeKind || f.Value == nil {
		return time.Time{}
	}
	return f.Value.(time.Time)
}

func (f *structField) SetDuration(v time.Duration) *structField {
	if f.err != nil {
		return f
	}
	f.kind = DurationKind
	f.Value = v
	return f
}

func (f *structField) GetDuration() time.Duration {
	if f.err != nil || f.kind != DurationKind || f.Value == nil {
		return 0
	}
	return f.Value.(time.Duration)
}
func (f *structField) GetSlice() []any {
	var arr []any
	f.Range(func(f Field) bool {
		arr = append(arr, f.GetObject())
		return true
	})
	return arr
}
func (f *structField) Range(iter func(Field) bool) {
	if f.err != nil || f.kind != SliceKind || f.Value == nil {
		return
	}
	for _, sf := range f.Value.([]*structField) {
		if !iter(sf) {
			return
		}
	}
}

func (f *structField) GetObject() any {
	switch f.kind {
	case IntKind:
		return f.GetInt()
	case UintKind:
		return f.GetUint()
	case StringKind:
		return f.GetString()
	case BoolKind:
		return f.GetBool()
	case TimeKind:
		return f.GetTime()
	case DurationKind:
		return f.GetDuration()
	case DynamicKind:
		return f.Value
	case SliceKind:
		return f.GetSlice()
	}
	return nil
}

func (f *structField) String() string {
	return fmt.Sprintf("structField{%v, Value=%v, kind=%v}", f.Key(), f.Value, f.kind)
}

func (f *structField) Valid() bool {
	return f != nil && f.Kind() != InvalidKind
}

// 获取字段数据类型
func (f *structField) Kind() (k KeyKind) {
	return f.kind
}
func (f structField) Key() Key {
	return f.KeyField
}

// 获取字段的键名
func (f *structField) Name() string {
	return f.KeyField.Name()
}

// 获取字段的键和值。返回 nil表示该字段无效
func (f *structField) Unwrap() (Key, interface{}) {
	if !f.Valid() {
		return f.KeyField, nil
	}
	return f.KeyField, f.GetObject()
}

// Fields 表示一个标签集合。
type Fields map[Key]Field

func (fs Fields) Copy() Fields {
	fieldsCopy := make(Fields, len(fs))
	for k, v := range fs {
		fieldsCopy[k] = v
	}
	return fieldsCopy
}
func (fs Fields) List() []Field {
	var arr []Field
	for _, v := range fs {
		arr = append(arr, v)
	}
	return arr
}
func (fs Fields) Set(f ...Field) {
	for _, it := range f {
		fs[it.Key()] = it
	}
}
func (fs Fields) Get(k Key, or ...interface{}) interface{} {
	// var v interface{}
	if k != nil {
		f := fs[k]
		if f != nil {
			_, v := f.Unwrap()
			return v
		}
	}

	l := len(or)
	if l == 0 {
		return nil
	}
	return or[l-1]
}
func (fs Fields) Del(k Key) interface{} {
	f, ok := fs[k]
	if !ok {
		return nil
	}
	delete(fs, k)
	if f != nil {
		_, v := f.Unwrap()
		return v
	}
	return nil
}

func (fs Fields) Keys() Keys {
	keys := make(Keys, 0)
	for _, f := range fs {
		k := f.Key()
		if k != nil {
			keys = append(keys, k)
		}
	}
	keys.Sort()
	return keys
}
