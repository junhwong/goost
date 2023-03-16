package field

import (
	"fmt"
	"time"

	"github.com/junhwong/goost/apm/field/pb"
)

// Field 表示一个标准字段。
type Field = pb.Field

// 字段标志位. 最多7+1个, 2**7-1
type Flags int32

const (
	FlagTag Flags = 1 << iota // 是标记

	// RES                       // 是资源
	// ATTR                      // 是属性
	// BDG                       // 是传播
	// BDY                       // 是内容
	// SRC                       // 源
)

func IsTag(f *Field) bool { return (Flags(f.GetFlags()) & FlagTag) == FlagTag }

type wrapField struct {
	pb.Field
}

func SetString(f *pb.Field, v string) *pb.Field {
	if len(v) == 0 {
		return f
	}
	f.Type = StringKind
	f.StringValue = v
	return f
}
func (f *wrapField) SetString(v string) *wrapField {
	SetString(&f.Field, v)
	return f
}

//	func (f *structField) GetString() string {
//		if f.err != nil || f.kind != StringKind || f.Value == nil {
//			return ""
//		}
//		return f.Value.(string)
//	}

func SetBool(f *pb.Field, v bool) *pb.Field {
	f.Type = BoolKind
	f.BoolValue = v
	return f
}

func (f *wrapField) SetBool(v bool) *wrapField {
	SetBool(&f.Field, v)
	return f
}

//	func (f *structField) GetBool() bool {
//		if f.err != nil || f.kind != BoolKind || f.Value == nil {
//			return false
//		}
//		return f.Value.(bool)
//	}
func SetInt(f *pb.Field, v int64) *pb.Field {
	f.Type = IntKind
	f.IntValue = v
	return f
}
func (f *wrapField) SetInt(v int64) *wrapField {
	SetInt(&f.Field, v)
	return f
}

//	func (f *structField) GetInt() int64 {
//		if f.err != nil || f.kind != IntKind || f.Value == nil {
//			return 0
//		}
//		return f.Value.(int64)
//	}
func SetUint(f *pb.Field, v uint64) *pb.Field {
	f.Type = UintKind
	f.UintValue = v
	return f
}
func (f *wrapField) SetUint(v uint64) *wrapField {
	SetUint(&f.Field, v)
	return f
}

//	func (f *structField) GetUint() uint64 {
//		if f.err != nil || f.kind != UintKind || f.Value == nil {
//			return 0
//		}
//		return f.Value.(uint64)
//	}
func SetFloat(f *pb.Field, v float64) *pb.Field {
	f.Type = FloatKind
	f.FloatValue = v
	return f
}
func (f *wrapField) SetFloat(v float64) *wrapField {
	SetFloat(&f.Field, v)
	return f
}

//	func (f *structField) GetFloat() float64 {
//		if f.err != nil || f.kind != FloatKind || f.Value == nil {
//			return 0
//		}
//		return f.Value.(float64)
//	}
func SetTime(f *pb.Field, v time.Time) *pb.Field {
	if v.IsZero() {
		return f
	}
	f.Type = TimeKind

	f.IntValue = v.UnixNano()
	return f
}
func (f *wrapField) SetTime(v time.Time) *wrapField {
	SetTime(&f.Field, v)
	return f
}

func (f wrapField) GetTimeValue() time.Time {
	return GetTimeValue(&f.Field)
}
func SetDuration(f *pb.Field, v time.Duration) *pb.Field {
	f.Type = DurationKind

	f.IntValue = int64(v)
	return f
}
func (f *wrapField) SetDuration(v time.Duration) *wrapField {
	SetDuration(&f.Field, v)
	return f
}

func (f wrapField) GetDurationValue() time.Duration {
	return GetDurationValue(&f.Field)
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
	if f == nil || f.Type != TimeKind {
		return time.Time{}
	}
	return time.Unix(0, f.GetIntValue())
}
func GetDurationValue(f *pb.Field) time.Duration {
	if f == nil || f.Type != DurationKind {
		return 0
	}
	return time.Duration(f.GetIntValue())
}

func GetObject(f *pb.Field) any {
	switch f.Type {
	case IntKind:
		return f.GetIntValue()
	case UintKind:
		return f.GetUintValue()
	case StringKind:
		return f.GetStringValue()
	case BoolKind:
		return f.GetBoolValue()
	case TimeKind:
		return GetTimeValue(f)
	case DurationKind:
		return GetDurationValue(f)
	}
	return nil
}

func (f wrapField) String() string {
	return fmt.Sprintf("structField{%v, Value=%v, kind=%v}", f.Key, GetObject(&f.Field), f.Type)
}

func (f *wrapField) Valid() bool {
	return f != nil && f.Type != InvalidKind
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
// func (f *wrapField) Unwrap() (string, interface{}) {
// 	if !f.Valid() {
// 		return f.Key, nil
// 	}
// 	return f.Key, f.GetObject()
// }

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
