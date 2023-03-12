package field

import (
	"fmt"
	"time"

	"github.com/junhwong/goost/apm/field/pb"
)

// Field 表示一个标准字段。
type Field = *pb.Field

// type Field interface {
// 	// Key() Key
// 	// Kind() KeyKind
// 	// Unwrap() (Key, interface{}) // 获取字段的键和值。返回 nil表示该字段无效
// 	// GetObject() any

// 	GetKey() string
// 	GetKind() KeyKind
// 	GetStringValue() string
// 	GetBoolValue() bool
// 	GetIntValue() int64
// 	GetUintValue() uint64
// 	GetFloatValue() float64

// 	// GetTimeValue() time.Time
// 	// GetDurationValue() time.Duration

// 	// GetSlice() []any
// }

// Field 表示一个标准字段。
//
// 参考：
//
// [opentelemetry](https://opentelemetry.io/docs/reference/specification/logs/overview/)
//
// [ecs](https://github.com/elastic/ecs)
// type structField struct {
// 	pb.Label
// 	// KeyField Key         `json:"key"`
// 	// Value    interface{} `json:"value"`
// 	valid bool // 防止自定义
// 	// sliceDataType KeyKind
// 	err error

// 	// ValueString *string
// 	// ValueInt    *int64
// 	// ValueUint   *uint64
// 	// ValueFloat  *float64
// 	// ValueBool   *bool

// 	// ValueSlice []*structField
// }

type structField pb.Field

func (f *structField) SetString(v string) *structField {
	if len(v) == 0 {
		return f
	}
	f.Kind = StringKind
	f.StringValue = &v
	return f
}

//	func (f *structField) GetString() string {
//		if f.err != nil || f.kind != StringKind || f.Value == nil {
//			return ""
//		}
//		return f.Value.(string)
//	}
func (f *structField) SetBool(v bool) *structField {
	f.Kind = BoolKind
	f.BoolValue = &v
	return f
}

//	func (f *structField) GetBool() bool {
//		if f.err != nil || f.kind != BoolKind || f.Value == nil {
//			return false
//		}
//		return f.Value.(bool)
//	}
func (f *structField) SetInt(v int64) *structField {
	f.Kind = IntKind
	f.IntValue = &v
	return f
}

// func (f *structField) GetInt() int64 {
// 	if f.err != nil || f.kind != IntKind || f.Value == nil {
// 		return 0
// 	}
// 	return f.Value.(int64)
// }

func (f *structField) SetUint(v uint64) *structField {
	f.Kind = UintKind
	f.UintValue = &v
	return f
}

// func (f *structField) GetUint() uint64 {
// 	if f.err != nil || f.kind != UintKind || f.Value == nil {
// 		return 0
// 	}
// 	return f.Value.(uint64)
// }

func (f *structField) SetFloat(v float64) *structField {
	f.Kind = FloatKind
	f.FloatValue = &v
	return f
}

// func (f *structField) GetFloat() float64 {
// 	if f.err != nil || f.kind != FloatKind || f.Value == nil {
// 		return 0
// 	}
// 	return f.Value.(float64)
// }

func (f *structField) SetTime(v time.Time) *structField {
	if v.IsZero() {
		return f
	}
	f.Kind = TimeKind
	u := v.UnixNano()
	f.IntValue = &u
	return f
}

func (f structField) GetTimeValue() time.Time {
	if f.Kind != TimeKind || f.IntValue == nil {
		return time.Time{}
	}
	ff := pb.Field(f)
	return time.Unix(0, ff.GetIntValue())
}

func (f *structField) SetDuration(v time.Duration) *structField {

	f.Kind = DurationKind
	d := int64(v)
	f.IntValue = &d
	return f
}

func (f structField) GetDurationValue() time.Duration {
	if f.Kind != DurationKind || f.IntValue == nil {
		return 0
	}
	ff := pb.Field(f)
	return time.Duration(ff.GetIntValue())
}

//	func (f *structField) GetSlice() []any {
//		var arr []any
//		f.Range(func(f Field) bool {
//			arr = append(arr, f.GetObject())
//			return true
//		})
//		return arr
//	}
//
//	func (f *structField) Range(iter func(Field) bool) {
//		if f.err != nil || f.kind != SliceKind || f.Value == nil {
//			return
//		}
//		for _, sf := range f.Value.([]*structField) {
//			if !iter(sf) {
//				return
//			}
//		}
//	}
func GetTimeValue(f *pb.Field) time.Time {
	if f == nil || f.Kind != TimeKind || f.IntValue == nil {
		return time.Time{}
	}
	return time.Unix(0, f.GetIntValue())
}
func GetDurationValue(f *pb.Field) time.Duration {
	if f == nil || f.Kind != DurationKind || f.IntValue == nil {
		return 0
	}
	return time.Duration(f.GetIntValue())
}
func (f structField) GetObject() any {
	ff := pb.Field(f)
	switch f.Kind {
	case IntKind:
		return ff.GetIntValue()
	case UintKind:
		return ff.GetUintValue()
	case StringKind:
		return ff.GetStringValue()
	case BoolKind:
		return ff.GetBoolValue()
	case TimeKind:
		return f.GetTimeValue()
	case DurationKind:
		return f.GetDurationValue()
		// case DynamicKind:
		// 	return f.Value
		// case SliceKind:
		// 	return f.GetSlice()
	}
	return nil
}
func GetObject(ff *pb.Field) any {
	switch ff.Kind {
	case IntKind:
		return ff.GetIntValue()
	case UintKind:
		return ff.GetUintValue()
	case StringKind:
		return ff.GetStringValue()
	case BoolKind:
		return ff.GetBoolValue()
	case TimeKind:
		return GetTimeValue(ff)
	case DurationKind:
		return GetDurationValue(ff)
		// case DynamicKind:
		// 	return f.Value
		// case SliceKind:
		// 	return f.GetSlice()
	}
	return nil
}

func (f *structField) String() string {
	return fmt.Sprintf("structField{%v, Value=%v, kind=%v}", f.Key, f.GetObject(), f.Kind)
}

func (f *structField) Valid() bool {
	return f != nil && f.Kind != InvalidKind
}

// // 获取字段数据类型
//
//	func (f *structField) Kind() (k KeyKind) {
//		return f.kind
//	}
// func (f structField) Key() Key {
// 	return f.KeyField
// }

// 获取字段的键名
// func (f *structField) Name() string {
// 	return f.KeyField.Name()
// }

// 获取字段的键和值。返回 nil表示该字段无效
func (f *structField) Unwrap() (string, interface{}) {
	if !f.Valid() {
		return f.Key, nil
	}
	return f.Key, f.GetObject()
}

// // Fields 表示一个标签集合。
// type Fields map[string]Field

// func (fs Fields) Copy() Fields {
// 	fieldsCopy := make(Fields, len(fs))
// 	for k, v := range fs {
// 		fieldsCopy[k] = v
// 	}
// 	return fieldsCopy
// }
// func (fs Fields) List() []Field {
// 	var arr []Field
// 	for _, v := range fs {
// 		arr = append(arr, v)
// 	}
// 	return arr
// }
// func (fs Fields) Set(f ...Field) {
// 	for _, it := range f {
// 		fs[it.GetKey()] = it
// 	}
// }
// func (fs Fields) Get(k string, or ...interface{}) interface{} {
// 	// var v interface{}
// 	// if k != "" {
// 	// 	f := fs[k]
// 	// 	if f != nil {
// 	// 		_, v := f.Unwrap()
// 	// 		return v
// 	// 	}
// 	// }

// 	// l := len(or)
// 	// if l == 0 {
// 	// 	return nil
// 	// }
// 	// return or[l-1]
// }
// func (fs Fields) Del(k string) interface{} {
// 	// f, ok := fs[k]
// 	// if !ok {
// 	// 	return nil
// 	// }
// 	// delete(fs, k)
// 	// if f != nil {
// 	// 	_, v := f.Unwrap()
// 	// 	return v
// 	// }
// 	return nil
// }

// // func (fs Fields) Keys() Keys {
// // 	keys := make(Keys, 0)
// // 	for _, f := range fs {
// // 		k := f.Key()
// // 		if k != nil {
// // 			keys = append(keys, k)
// // 		}
// // 	}
// // 	keys.Sort()
// // 	return keys
// // }
