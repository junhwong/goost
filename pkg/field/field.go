package field

import "encoding/json"

// Field 表示一个标签。
type Field struct {
	Key           Key         `json:"key"`
	Value         interface{} `json:"value"`
	valid         bool        // 防止自定义
	sliceDataType KeyKind
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

// func (f *Field) val(k KeyKind) (v interface{}) {
// 	if f.IsValid() && f.Key.Kind() == k {
// 		v = f.Value
// 	}
// 	return
// }

// GetValue 返回该字段的 `StringKind` 值.
func (f *Field) GetValue() (v interface{}) {
	if !f.IsValid() {
		return
	}
	switch f.Key.Kind() {
	case SliceKind:
		s, _ := f.Value.([]interface{})
		objs := make([]interface{}, 0)
		for _, it := range s {
			switch f.sliceDataType {
			case IntKind:
				panic("未实现")
			default:
				var val string
				switch v := it.(type) {
				case string:
					val = v
				case *string:
					val = *v
				default:
					data, _ := json.Marshal(v)
					val = string(data)
				}
				objs = append(objs, val)
			}
		}
		// if s == nil {
		// 	s = make([]interface{}, 0)
		// }
		v = objs
	case StringsKind:
		s, _ := f.Value.([]interface{})
		objs := make([]interface{}, 0)
		for _, it := range s {
			var val string
			switch v := it.(type) {
			case string:
				val = v
			case *string:
				val = *v
			default:
				data, _ := json.Marshal(v)
				val = string(data)
			}
			objs = append(objs, val)
		}
		v = objs
	case MapKind:
		s := f.Value.([]*Field)
		objs := make(map[string]interface{}, len(s))
		for _, it := range s {
			objs[it.Key.Name()] = it.GetValue()
		}
		v = objs
	default:
		v = f.Value
	}
	return
	// data, _ := json.Marshal(v)
	// return data
}

// Fields 表示一个标签集合。
type Fields map[Key]*Field

func (fs Fields) Set(f *Field) {
	if f == nil || !f.IsValid() {
		return
	}
	fs[f.Key] = f
}
