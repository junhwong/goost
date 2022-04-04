package field

// Field 表示一个标准字段。
type Field interface {
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
	Key           Key         `json:"key"`
	Value         interface{} `json:"value"`
	valid         bool        // 防止自定义
	sliceDataType KeyKind
}

func (f *structField) Valid() bool {
	return f != nil && f.valid && f.Key != nil
}

// 获取字段数据类型
func (f *structField) Kind() (k KeyKind) {
	if f.Valid() {
		return f.Key.Kind()
	}
	return InvalidKind
}

// 获取字段的键名
func (f *structField) Name() string {
	return f.Key.Name()
}

// 获取字段的键和值。返回 nil表示该字段无效
func (f *structField) Unwrap() (Key, interface{}) {
	if f == nil || !f.valid {
		return nil, nil
	}
	return f.Key, f.Value
}

// Fields 表示一个标签集合。
type Fields map[Key]interface{}

func (fs Fields) Set(f Field) {
	if f == nil {
		return
	}
	k, v := f.Unwrap()
	if k == nil || v == nil {
		return
	}
	fs[k] = v
}
func (fs Fields) Get(k Key, or ...interface{}) interface{} {
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
func (fs Fields) Del(k Key) interface{} {
	v, ok := fs[k]
	if !ok {
		return nil
	}
	delete(fs, k)
	return v
}
