package field

import (
	"fmt"
	"strings"
	"time"

	"github.com/spf13/cast"
)

func makeField(k Key, v interface{}, valid bool, err ...error) *Field {
	var ex error
	for _, it := range err {
		ex = it
	}
	if ex != nil {
		valid = false
	}
	return &Field{Key: k, Value: v, valid: valid}
}

// 构造一个动态字段
func Dynamic(name string) func(v interface{}) *Field {
	return func(v interface{}) *Field {
		return &Field{Key: &key{
			name: name,
			kind: DynamicKind,
		}, Value: v, valid: v != nil && v != ""}
	}
}

func String(name string) (Key, func(string, ...interface{}) *Field) {
	k := makeOrGetKey(name, StringKind)
	return k, func(s string, a ...interface{}) *Field {
		v := s
		if s != "" && len(a) > 0 {
			v = fmt.Sprintf(s, a...)
		} else if len(a) > 0 {
			v = fmt.Sprint(a...)
		}
		v = strings.TrimSpace(v)
		return makeField(k, v, v != "")
	}
}

// func Strings(name string) func(a ...interface{}) *Field {
// 	k := makeOrGetKey(name, StringsKind)
// 	return func(a ...interface{}) *Field {
// 		return makeField(k, a, true)
// 	}
// }

func Bool(name string) (Key, func(bool) *Field) {
	k := makeOrGetKey(name, BoolKind)
	return k, func(v bool) *Field {
		return makeField(k, v, true)
	}
}

func Time(name string) (Key, func(time.Time) *Field) {
	k := makeOrGetKey(name, TimeKind)
	return k, func(t time.Time) *Field {
		return makeField(k, t, !t.IsZero())
	}
}

// Slice 返回一个数组对象
//	like json:
//	```json
//	{
//		key: [1, "string", true]
//	}
//	```
func Slice(name string, dataType ...KeyKind) (Key, func(...interface{}) *Field) {
	dt := StringKind
	for _, t := range dataType {
		dt = t
	}
	k := makeOrGetKey(name, SliceKind)
	return k, func(v ...interface{}) *Field {
		f := makeField(k, v, true)
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
		return makeField(k, v, true)
	}
}

func Int(name string) (Key, func(interface{}) *Field) {
	k := makeOrGetKey(name, IntKind)
	return k, func(v interface{}) *Field {
		i, err := cast.ToInt64E(v)
		return makeField(k, i, err == nil, err)
	}
}

func Uint(name string) (Key, func(interface{}) *Field) {
	k := makeOrGetKey(name, UintKind)
	return k, func(v interface{}) *Field {
		i, err := cast.ToUint64E(v)
		return makeField(k, i, err == nil, err)
	}
}

func Float(name string) (Key, func(interface{}) *Field) {
	k := makeOrGetKey(name, FloatKind)
	return k, func(v interface{}) *Field {
		i, err := cast.ToFloat64E(v)
		return makeField(k, i, err == nil, err)
	}
}

// 时延。微秒
func Duration(name string) (Key, func(time.Duration) *Field) {
	k := makeOrGetKey(name, IntKind)
	return k, func(v time.Duration) *Field {
		return makeField(k, v.Microseconds(), v >= 0)
	}
}

// func Byte(name string) (Key,func(byte) *Field) {
// 	k := makeOrGetKey(name, IntKind)
// 	return k,func(v byte) *Field {
// 		return makeField(k, int64(v))
// 	}
// }

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
