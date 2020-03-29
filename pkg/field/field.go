package field

import (
	"encoding/json"
	"fmt"
)

// Field 表示一个标签。
type Field struct {
	Key   Key         `json:"key"`
	Value interface{} `json:"value"`
	valid bool        // 防止自定义
}

func (f *Field) IsValid() bool {
	return f != nil && f.valid && f.Key != nil
}

func (f *Field) Kind() (k KeyKind) {
	if f.IsValid() {
		return f.Key.Kind()
	}
	return InvalidKind
}

func (f *Field) val(k KeyKind) (v interface{}) {
	if f.IsValid() && f.Key.Kind() == k {
		v = f.Value
	}
	return
}

// GetValue 返回该字段的 `StringKind` 值.
func (f *Field) GetValue() (v interface{}) {
	if !f.IsValid() {
		return
	}
	switch f.Key.Kind() {
	case StringKind:
		v = fmt.Sprintf("%q", f.Value)
	case SliceKind:
		s := f.Value.([]*Field)
		objs := make([]interface{}, len(s))
		for _, it := range s {
			objs = append(objs, it.GetValue())
		}
		v = objs
	default:
		v = f.Value
	}
	data, _ := json.Marshal(v)
	return data
}

// Fields 表示一个标签集合。
type Fields map[Key]*Field
