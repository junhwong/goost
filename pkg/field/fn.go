package field

import (
	"fmt"
	"strings"
	"time"

	"github.com/spf13/cast"
)

func makeField(k Key, v interface{}, valid bool, err ...error) Field {
	var ex error
	for _, it := range err {
		ex = it
	}
	if ex != nil {
		valid = false
	}
	if !valid {
		v = nil
	}
	// return fnField(func() (Key, interface{}) {
	// 	return k, v
	// })
	return &structField{Key: k, Value: v, valid: valid}
}

type fnField func() (Key, interface{})

func (fn fnField) Unwrap() (Key, interface{}) {
	return fn()
}

// 构造一个动态字段
func Dynamic(name string, checkKey ...bool) (Key, func(v interface{}) Field) {
	check := false
	for _, b := range checkKey {
		check = b
	}
	var k Key
	if check {
		k = makeOrGetKey(name, DynamicKind)
	} else {
		if !IsValidKeyName(name) {
			panic(fmt.Errorf("field: Invalid key name: %s", name))
		}
		k = &key{name: name, kind: DynamicKind}
	}

	return k, func(v interface{}) Field {
		return makeField(k, v, v != nil && v != "")
	}
}

func String(name string) (Key, func(string, ...interface{}) Field) {
	k := makeOrGetKey(name, StringKind)
	return k, func(s string, a ...interface{}) Field {
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

// func Strings(name string) func(a ...interface{}) Field {
// 	k := makeOrGetKey(name, StringsKind)
// 	return func(a ...interface{}) Field {
// 		return makeField(k, a, true)
// 	}
// }

func Bool(name string) (Key, func(bool) Field) {
	k := makeOrGetKey(name, BoolKind)
	return k, func(v bool) Field {
		return makeField(k, v, true)
	}
}

func Time(name string) (Key, func(time.Time) Field) {
	k := makeOrGetKey(name, TimeKind)
	return k, func(t time.Time) Field {
		return makeField(k, t, !t.IsZero())
	}
}

// Slice 返回一个数组对象
func Slice(name string, kind ...KeyKind) (Key, func(...interface{}) Field) {
	k := &key{name: name}
	dtype := DynamicKind
	if len(kind) > 0 {
		dtype = kind[len(kind)-1]
	}
	return k, func(v ...interface{}) (r Field) {
		switch dtype {
		case StringKind:
			val, err := cast.ToStringSliceE(v)
			r = makeField(k, val, len(val) > 0, err)
		case IntKind:
			val, err := ToInt64SliceE(v)
			r = makeField(k, val, len(val) > 0, err)
		case UintKind:
			val, err := ToUint64SliceE(v)
			r = makeField(k, val, len(val) > 0, err)
		case FloatKind:
			val, err := ToFloat64SliceE(v)
			r = makeField(k, val, len(val) > 0, err)
		case BoolKind:
			val, err := cast.ToBoolSliceE(v)
			r = makeField(k, val, len(val) > 0, err)
		case TimeKind:
			// TODO: 时区未解决, 目前是UTC
			val := []time.Time{}
			var err error
			for _, it := range v {
				if t, ex := cast.ToTimeE(it); ex == nil {
					val = append(val, t)
				} else {
					err = ex
					break
				}
			}
			r = makeField(k, val, len(val) > 0, err)
		default:
			r = makeField(k, v, len(v) > 0, nil)
		}
		return
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
func Map(name string) func(v ...Field) Field {
	k := makeOrGetKey(name, MapKind)
	return func(v ...Field) Field {
		return makeField(k, v, true)
	}
}

func Int(name string) (Key, func(interface{}) Field) {
	k := makeOrGetKey(name, IntKind)
	return k, func(v interface{}) Field {
		i, err := cast.ToInt64E(v)
		return makeField(k, i, err == nil, err)
	}
}

func Uint(name string) (Key, func(interface{}) Field) {
	k := makeOrGetKey(name, UintKind)
	return k, func(v interface{}) Field {
		i, err := cast.ToUint64E(v)
		return makeField(k, i, err == nil, err)
	}
}

func Float(name string) (Key, func(interface{}) Field) {
	k := makeOrGetKey(name, FloatKind)
	return k, func(v interface{}) Field {
		i, err := cast.ToFloat64E(v)
		return makeField(k, i, err == nil, err)
	}
}

// 时延。微秒
func Duration(name string) (Key, func(time.Duration) Field) {
	k := makeOrGetKey(name, IntKind)
	return k, func(v time.Duration) Field {
		return makeField(k, v.Microseconds(), v >= 0)
	}
}

// func Byte(name string) (Key,func(byte) Field) {
// 	k := makeOrGetKey(name, IntKind)
// 	return k,func(v byte) Field {
// 		return makeField(k, int64(v))
// 	}
// }

// func Any(name string, value interface{}) Field {
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
