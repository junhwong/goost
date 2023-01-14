package field

import (
	"fmt"
)

// Field 表示一个标准字段。
type Field interface {
	Key() Key
	Unwrap() (Key, interface{}) // 获取字段的键和值。返回 nil表示该字段无效
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
	valid    bool        // 防止自定义

	// sliceDataType KeyKind
	// Int   int64
	// Uint  uint64
	// Float float64
	// Time  time.Time
	// Bool  bool
}

func (f *structField) String() string {
	return fmt.Sprintf("structField{%v, Value=%v, valid=%v}", f.Key(), f.Value, f.Valid())
}

func (f *structField) Valid() bool {
	return f != nil && f.valid && f.KeyField.Kind() != InvalidKind
}

// 获取字段数据类型
func (f *structField) Kind() (k KeyKind) {
	return f.KeyField.Kind()
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
	return f.KeyField, f.Value
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
