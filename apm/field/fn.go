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
	return &structField{KeyField: k, Value: v, valid: valid}
}

func infer(v any, trem KeyKind) (val any, kind KeyKind, b bool, err error) {
	if v == nil {
		kind = InvalidKind
		return
	}
	b = true
	switch v := v.(type) {
	case int:
		val = int64(v)
		kind = IntKind
	case int8:
		val = int64(v)
		kind = IntKind
	case int16:
		val = int64(v)
		kind = IntKind
	case int32:
		val = int64(v)
		kind = IntKind
	case int64:
		val = v
		kind = IntKind
	case *int:
		val = int64(*v)
		kind = IntKind
	case *int8:
		val = int64(*v)
		kind = IntKind
	case *int16:
		val = int64(*v)
		kind = IntKind
	case *int32:
		val = int64(*v)
		kind = IntKind
	case *int64:
		val = *v
		kind = IntKind

	case uint:
		val = uint64(v)
		kind = UintKind
	case uint8:
		val = uint64(v)
		kind = UintKind
	case uint16:
		val = uint64(v)
		kind = UintKind
	case uint32:
		val = uint64(v)
		kind = UintKind
	case uint64:
		val = v
		kind = UintKind
	case *uint:
		val = uint64(*v)
		kind = UintKind
	case *uint8:
		val = uint64(*v)
		kind = UintKind
	case *uint16:
		val = uint64(*v)
		kind = UintKind
	case *uint32:
		val = uint64(*v)
		kind = UintKind
	case *uint64:
		val = *v
		kind = UintKind

	case float32:
		val = float64(v)
		kind = FloatKind
	case float64:
		val = v
		kind = FloatKind
	case *float32:
		val = float64(*v)
		kind = FloatKind
	case *float64:
		val = *v
		kind = FloatKind

	case bool:
		val = v
		kind = BoolKind
	case *bool:
		val = *v
		kind = BoolKind

	case string:
		val = v
		kind = StringKind
		b = len(v) > 0
	case *string:
		val = *v
		kind = StringKind
		b = len(*v) > 0

	case time.Time:
		val = v
		kind = TimeKind
		b = !v.IsZero()
	case *time.Time:
		val = *v
		kind = TimeKind
		b = !(*v).IsZero()

	default:
		val = v
		kind = DynamicKind
	}
	if trem == DynamicKind || trem == kind {
		return
	}
	kind = trem
	switch trem {
	case StringKind:
		v, cerr := cast.ToStringE(val)
		if cerr != nil {
			kind = InvalidKind
			err = cerr
			return
		}
		val = v
		b = len(v) > 0
	case IntKind:
		v, cerr := cast.ToInt64E(val)
		if cerr != nil {
			kind = InvalidKind
			err = cerr
			return
		}
		val = v
	case UintKind:
		v, cerr := cast.ToUint64E(val)
		if cerr != nil {
			kind = InvalidKind
			err = cerr
			return
		}
		val = v
	case BoolKind:
		v, cerr := cast.ToBoolE(val)
		if cerr != nil {
			kind = InvalidKind
			err = cerr
			return
		}
		val = v
	case TimeKind:
		v, cerr := cast.ToTimeE(val)
		if cerr != nil {
			kind = InvalidKind
			err = cerr
			return
		}
		b = !v.IsZero()
	case DynamicKind:
		kind = DynamicKind
	default:
		kind = InvalidKind
	}
	return
}

// 构造一个动态字段
func Dynamic(name string, v any) Field {
	val, k, valid, _ := infer(v, DynamicKind)
	if !valid || k == InvalidKind {
		return &structField{KeyField: key{name: name, kind: k}, Value: v, valid: valid}
	}
	key := makeOrGetKey(name, k)
	return makeField(key, val, valid)
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

func Int(name string) (Key, func(interface{}) Field) {
	k := makeOrGetKey(name, IntKind)
	return k, func(v interface{}) Field {
		i, err := cast.ToInt64E(v)
		if err == nil {
			v = i
		}
		return makeField(k, v, err == nil, err)
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

// Slice 返回一个数组对象
func Slice(name string, kind KeyKind) (Key, func(...interface{}) Field) {
	k := &key{name: name, kind: kind}
	// dtype := DynamicKind
	// if len(kind) > 0 {
	// 	dtype = kind[len(kind)-1]
	// }
	return k, func(v ...interface{}) (r Field) {
		switch kind {
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
//
//		like json:
//		```json
//		{
//	 	key: {
//				subkey: 1,
//				subkey: "string"
//			}
//		}
//		```
func Map(name string) func(v ...Field) Field {
	k := makeOrGetKey(name, MapKind)
	return func(v ...Field) Field {
		return makeField(k, v, true)
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
