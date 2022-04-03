package field

// Field 表示一个标准字段。
//
// 参考：
//
// [opentelemetry](https://opentelemetry.io/docs/reference/specification/logs/overview/)
//
// [ecs](https://github.com/elastic/ecs)
type Field struct {
	Key           Key         `json:"key"`
	Value         interface{} `json:"value"`
	valid         bool        // 防止自定义
	sliceDataType KeyKind
}

func (f *Field) Valid() bool {
	return f != nil && f.valid && f.Key != nil
}

func (f *Field) Kind() (k KeyKind) {
	if f.Valid() {
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

// // GetValue 返回该字段的 `StringKind` 值.
// func (f *Field) GetValue() (v interface{}) {
// 	if !f.IsValid() {
// 		return
// 	}
// 	switch f.Key.Kind() {
// 	case StringKind:
// 	case IntKind:
// 	case UintKind:
// 	case FloatKind:
// 	case TimeKind:
// 	case SliceKind:
// 		s, _ := f.Value.([]interface{})
// 		objs := make([]interface{}, 0)
// 		for _, it := range s {
// 			switch f.sliceDataType {
// 			case IntKind:
// 				panic("未实现")
// 			default:
// 				var val string
// 				switch v := it.(type) {
// 				case string:
// 					val = v
// 				case *string:
// 					val = *v
// 				default:
// 					data, _ := json.Marshal(v)
// 					val = string(data)
// 				}
// 				objs = append(objs, val)
// 			}
// 		}

// 		v = objs
// 	// case StringsKind:
// 	// 	s, _ := f.Value.([]interface{})
// 	// 	objs := make([]interface{}, 0)
// 	// 	for _, it := range s {
// 	// 		var val string
// 	// 		switch v := it.(type) {
// 	// 		case string:
// 	// 			val = v
// 	// 		case *string:
// 	// 			val = *v
// 	// 		default:
// 	// 			data, _ := json.Marshal(v)
// 	// 			val = string(data)
// 	// 		}
// 	// 		objs = append(objs, val)
// 	// 	}
// 	// 	v = objs
// 	case MapKind:
// 		s := f.Value.([]*Field)
// 		objs := make(map[string]interface{}, len(s))
// 		for _, it := range s {
// 			objs[it.Key.Name()] = it.GetValue()
// 		}
// 		v = objs
// 	case DynamicKind:
// 		if u, ok := f.Value.(*url.URL); ok && u != nil {
// 			v = u.String()
// 		} else if u, ok := f.Value.(url.URL); ok {
// 			v = u.String()
// 		} else {
// 			v = f.Value
// 		}
// 	default:
// 		v = f.Value
// 	}
// 	return
// 	// data, _ := json.Marshal(v)
// 	// return data
// }

// Fields 表示一个标签集合。
type Fields map[Key]*Field

func (fs Fields) Set(f *Field) {
	if f == nil || !f.Valid() {
		return
	}
	fs[f.Key] = f
}
func (fs Fields) Get(k Key, or ...*Field) *Field {
	v := fs[k]
	if v != nil {
		return v
	}
	l := len(or)
	if l == 0 {
		return nil
	}
	return or[l-1]
}
func (fs Fields) Del(k Key) *Field {
	v, ok := fs[k]
	if !ok {
		return nil
	}
	delete(fs, k)
	return v
}
