package field

import (
	"fmt"
	"strings"
	"time"
)

func makeField(k Key, v interface{}, b ...bool) *Field {
	valid := true
	for _, it := range b {
		valid = it
	}
	return &Field{Key: k, Value: v, valid: valid}
}

func String(name string) func(v string) *Field {
	k := makeOrGetKey(name, StringKind)
	return func(v string) *Field {
		v = strings.TrimSpace(v)
		return makeField(k, v, v != "")
	}
}
func Stringf(name string) func(s string, a ...interface{}) *Field {
	k := makeOrGetKey(name, StringKind)
	return func(s string, a ...interface{}) *Field {
		v := fmt.Sprintf(s, a...)
		v = strings.TrimSpace(v)
		return makeField(k, v, true)
	}
}
func Strings(name string) func(a ...interface{}) *Field {
	k := makeOrGetKey(name, StringsKind)
	return func(a ...interface{}) *Field {
		return makeField(k, a, true)
	}
}

// Slice 返回一个数组对象
//	like json:
//	```json
//	{
//		key: [1, "string", true]
//	}
//	```
func Slice(name string, dataType ...KeyKind) func(v ...interface{}) *Field {
	dt := StringKind
	for _, t := range dataType {
		dt = t
	}
	k := makeOrGetKey(name, SliceKind)
	return func(v ...interface{}) *Field {
		f := makeField(k, v)
		f.sliceDataType = dt
		return f
	}
}

// Map 返回一个嵌套对象
//	like json:
//	```json
//	{
//  	key: {
//			subkey: 1,
//			subkey: "string"
//		}
//	}
//	```
func Map(name string) func(v ...*Field) *Field {
	k := makeOrGetKey(name, MapKind)
	return func(v ...*Field) *Field {
		return makeField(k, v)
	}
}

func Int(name string) func(v int) *Field {
	k := makeOrGetKey(name, IntKind)
	return func(v int) *Field {
		return makeField(k, int64(v))
	}
}
func Int8(name string) func(v int8) *Field {
	k := makeOrGetKey(name, IntKind)
	return func(v int8) *Field {
		return makeField(k, int64(v))
	}
}
func Int16(name string) func(v int16) *Field {
	k := makeOrGetKey(name, IntKind)
	return func(v int16) *Field {
		return makeField(k, int64(v))
	}
}
func Int32(name string) func(v int32) *Field {
	k := makeOrGetKey(name, IntKind)
	return func(v int32) *Field {
		return makeField(k, int64(v))
	}
}
func Int64(name string) func(v int64) *Field {
	k := makeOrGetKey(name, IntKind)
	return func(v int64) *Field {
		return makeField(k, v)
	}
}

func Uint(name string) func(v uint) *Field {
	k := makeOrGetKey(name, IntKind)
	return func(v uint) *Field {
		return makeField(k, int64(v))
	}
}
func Uint8(name string) func(v uint8) *Field {
	k := makeOrGetKey(name, IntKind)
	return func(v uint8) *Field {
		return makeField(k, int64(v))
	}
}
func Uint16(name string) func(v uint16) *Field {
	k := makeOrGetKey(name, IntKind)
	return func(v uint16) *Field {
		return makeField(k, int64(v))
	}
}
func Uint32(name string) func(v uint32) *Field {
	k := makeOrGetKey(name, IntKind)
	return func(v uint32) *Field {
		return makeField(k, int64(v))
	}
}
func Uint64(name string) func(v uint64) *Field {
	k := makeOrGetKey(name, IntKind)
	return func(v uint64) *Field {
		return makeField(k, int64(v))
	}
}
func Duration(name string) func(v time.Duration) *Field {
	k := makeOrGetKey(name, IntKind)
	return func(v time.Duration) *Field {
		return makeField(k, v.Nanoseconds(), v >= 0)
	}
}

func Byte(name string) func(v byte) *Field {
	k := makeOrGetKey(name, IntKind)
	return func(v byte) *Field {
		return makeField(k, int64(v))
	}
}

func Float32(name string) func(v float32) *Field {
	k := makeOrGetKey(name, FloatKind)
	return func(v float32) *Field {
		return makeField(k, float64(v))
	}
}
func Float(name string) func(v float64) *Field {
	k := makeOrGetKey(name, FloatKind)
	return func(v float64) *Field {
		return makeField(k, float64(v))
	}
}

// func Any(name string, value interface{}) *Field {
// 	switch v := value.(type) {
// 	case int:
// 		return newField(name, int64(v), FKInteger)
// 	case uint:
// 		return newField(name, int64(v), FKInteger)
// 	case int16:
// 		return newField(name, int64(v), FKInteger)
// 	case uint16:
// 		return newField(name, int64(v), FKInteger)
// 	case int32:
// 		return newField(name, int64(v), FKInteger)
// 	case uint32:
// 		return newField(name, int64(v), FKInteger)
// 	case int64:
// 		return newField(name, v, FKInteger)
// 	case uint64:
// 		return newField(name, int64(v), FKInteger)
// 	case uint8:
// 		return newField(name, int64(v), FKInteger)
// 	case uintptr:
// 		return newField(name, int64(v), FKInteger)
// 	case float32:
// 		f := float64(v)
// 		if !strings.Contains(strconv.FormatFloat(f, 'f', -1, 64), ".") {
// 			return newField(name, int64(v), FKInteger)
// 		} else {
// 			return newField(name, f, FKFloat)
// 		}
// 	case float64:
// 		if !strings.Contains(strconv.FormatFloat(v, 'f', -1, 64), ".") {
// 			return newField(name, int64(v), FKInteger)
// 		} else {
// 			return newField(name, v, FKFloat)
// 		}
// 	case bool:
// 		return newField(name, v, FKBool)
// 	case string:
// 		t := FKString
// 		if v == "" {
// 			t = FKInvalid
// 		}
// 		return newField(name, v, t)
// 	case time.Time:
// 		return newField(name, v, FKTime)
// 	case time.Duration:
// 		return newField(name, v, FKDuration)
// 	default:
// 		return newField(name, v, FKAny)
// 	}
// }
