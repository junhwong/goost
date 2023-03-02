package field

import (
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/spf13/cast"
)

func makeField(k Key, v interface{}, valid bool, err ...error) *structField {
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
func makeField2(k Key, err ...error) *structField {
	var ex error
	for _, it := range err {
		ex = it
	}
	if ex != nil {
	}
	return &structField{KeyField: k}
}
func infer(v any) *structField {
	f := &structField{kind: InvalidKind}
	iv, k := InferPrimitiveValue(v)
	if k != InvalidKind {
		f.Value = iv
		f.kind = k
		return f
	}
	rv := reflect.ValueOf(v)
	if rv.Kind() == reflect.Pointer {
		rv = rv.Elem()
	}
	iv, k = InferPrimitiveValueByReflect(rv)
	if k != InvalidKind {
		f.Value = iv
		f.kind = k
		return f
	}
	switch rv.Kind() {
	case reflect.Slice, reflect.Array:
		var slice []*structField
		i := rv.Len()
		for i > -1 {
			i--
			iv := rv.Index(i)
			if iv.Kind() != reflect.Invalid {
				sf := infer(iv.Interface())
				if sf.kind != InvalidKind {
					slice = append(slice, sf)
				}
			}
		}
		f.Value = slice
		f.kind = SliceKind
		return f
	case reflect.Struct:
		f.Value = v
		f.kind = DynamicKind
		return f
	}

	return f
	// f := &structField{kind: InvalidKind}
	// if v == nil {
	// 	return f
	// }
	// switch v := v.(type) {
	// case int:
	// 	f.SetInt(int64(v))
	// case int8:
	// 	f.SetInt(int64(v))
	// case int16:
	// 	f.SetInt(int64(v))
	// case int32:
	// 	f.SetInt(int64(v))
	// case int64:
	// 	f.SetInt(v)
	// case *int:
	// 	f.SetInt(int64(*v))
	// case *int8:
	// 	f.SetInt(int64(*v))
	// case *int16:
	// 	f.SetInt(int64(*v))
	// case *int32:
	// 	f.SetInt(int64(*v))
	// case *int64:
	// 	f.SetInt(int64(*v))
	// case uint:
	// 	f.SetUint(uint64(v))
	// case uint8:
	// 	f.SetUint(uint64(v))
	// case uint16:
	// 	f.SetUint(uint64(v))
	// case uint32:
	// 	f.SetUint(uint64(v))
	// case uint64:
	// 	f.SetUint(v)
	// case *uint:
	// 	f.SetUint(uint64(*v))
	// case *uint8:
	// 	f.SetUint(uint64(*v))
	// case *uint16:
	// 	f.SetUint(uint64(*v))
	// case *uint32:
	// 	f.SetUint(uint64(*v))
	// case *uint64:
	// 	f.SetUint(*v)
	// case float32:
	// 	f.SetFloat(float64(v))
	// case float64:
	// 	f.SetFloat(v)
	// case *float32:
	// 	f.SetFloat(float64(*v))
	// case *float64:
	// 	f.SetFloat(*v)
	// case bool:
	// 	f.SetBool(v)
	// case *bool:
	// 	f.SetBool(*v)
	// case string:
	// 	f.SetString(v)
	// case *string:
	// 	f.SetString(*v)
	// case time.Time:
	// 	f.SetTime(v)
	// case *time.Time:
	// 	f.SetTime(*v)
	// case time.Duration:
	// 	f.SetDuration(v)
	// case *time.Duration:
	// 	f.SetDuration(*v)
	// }
	// if f.kind != InvalidKind {
	// 	return f
	// }
	// // kind = trem
	// rv := reflect.ValueOf(v)
	// if rv.Kind() == reflect.Pointer {
	// 	rv = rv.Elem()
	// }
	// switch rv.Kind() {
	// case reflect.Bool:
	// case reflect.String:
	// 	v = rv.Interface()
	// 	if rv.Type().String() != "string" {
	// 		v = fmt.Sprint(v)
	// 	}
	// 	return infer(v, trem)
	// case reflect.Slice, reflect.Array:
	// 	i := rv.Len()
	// 	for i > -1 {
	// 		i--
	// 		iv := rv.Index(i)
	// 		if iv.Kind() != reflect.Invalid {
	// 			sf := infer(iv.Interface(), trem)
	// 			if sf.kind != InvalidKind {
	// 				f.kind = SliceKind
	// 				f.ValueSlice = append(f.ValueSlice, sf)
	// 			}
	// 		}
	// 	}
	// }

	// return f
}

func InferPrimitiveValue(v any) (any, KeyKind) {
	if v == nil {
		return nil, InvalidKind
	}
	switch v := v.(type) {
	case int:
		return int64(v), IntKind
	case int8:
		return int64(v), IntKind
	case int16:
		return int64(v), IntKind
	case int32:
		return int64(v), IntKind
	case int64:
		return v, IntKind
	case *int:
		return int64(*v), IntKind
	case *int8:
		return int64(*v), IntKind
	case *int16:
		return int64(*v), IntKind
	case *int32:
		return int64(*v), IntKind
	case *int64:
		return *v, IntKind

	case uint:
		return uint64(v), UintKind
	case uint8:
		return uint64(v), UintKind
	case uint16:
		return uint64(v), UintKind
	case uint32:
		return uint64(v), UintKind
	case uint64:
		return v, UintKind
	case *uint:
		return uint64(*v), UintKind
	case *uint8:
		return uint64(*v), UintKind
	case *uint16:
		return uint64(*v), UintKind
	case *uint32:
		return uint64(*v), UintKind
	case *uint64:
		return *v, UintKind

	case float32:
		return float64(v), FloatKind
	case float64:
		return v, FloatKind
	case *float32:
		return float64(*v), FloatKind
	case *float64:
		return *v, FloatKind

	case bool:
		return v, BoolKind
	case *bool:
		return *v, BoolKind

	case string:
		return v, StringKind
	case *string:
		return *v, StringKind

	case time.Time:
		return v, TimeKind
	case *time.Time:
		return *v, TimeKind

	case time.Duration:
		return v, DurationKind
	case *time.Duration:
		return *v, DurationKind
	}
	return nil, InvalidKind
}

// 反射获取重新定义基础类型的值
func InferPrimitiveValueByReflect(rv reflect.Value) (any, KeyKind) {
	if rv.Kind() == reflect.Pointer {
		rv = rv.Elem()
	}
	switch rv.Kind() {
	case reflect.Bool:
		return rv.Bool(), BoolKind
	case reflect.Int64:
		v := rv.Int()
		if rv.Type().String() == "time.Duration" {
			return v, DurationKind
		}
		return v, IntKind
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32:
		return rv.Int(), IntKind
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return rv.Uint(), UintKind
	case reflect.Float32, reflect.Float64:
		return rv.Float(), FloatKind
	case reflect.Complex64, reflect.Complex128:
		panic("TODO Complex")
	case reflect.String:
		return rv.String(), StringKind
	case reflect.Struct:
		v := rv.Interface()
		if rv.Type().String() == "time.Time" {
			if n, err := cast.ToTimeE(v); err == nil {
				return n, TimeKind
			}
		}
	}
	fmt.Printf("rv.Kind(): %v\n", rv.Kind() == reflect.Int64)
	return nil, InvalidKind
}

// 构造一个动态字段
func Dynamic(name string, v any) Field {
	f := infer(v)
	f.KeyField = makeOrGetKey(name, f.kind)
	return f
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
		return makeField2(k).SetString(strings.TrimSpace(v))
	}
}

func Bool(name string) (Key, func(bool) Field) {
	k := makeOrGetKey(name, BoolKind)
	return k, func(v bool) Field {
		return makeField2(k).SetBool(v)
	}
}

func Time(name string) (Key, func(time.Time) Field) {
	k := makeOrGetKey(name, TimeKind)
	return k, func(t time.Time) Field {
		return makeField2(k).SetTime(t)
	}
}

func Int(name string) (Key, func(interface{}) Field) {
	k := makeOrGetKey(name, IntKind)
	return k, func(v interface{}) Field {
		i, err := cast.ToInt64E(v)
		if err == nil {
			v = i
		}
		return makeField2(k, err).SetInt(i)
	}
}

func Uint(name string) (Key, func(interface{}) Field) {
	k := makeOrGetKey(name, UintKind)
	return k, func(v interface{}) Field {
		i, err := cast.ToUint64E(v)
		return makeField2(k, err).SetUint(i)
	}
}

func Float(name string) (Key, func(interface{}) Field) {
	k := makeOrGetKey(name, FloatKind)
	return k, func(v interface{}) Field {
		n, err := cast.ToFloat64E(v)
		return makeField2(k, err).SetFloat(n)
	}
}

// 时延。纳秒?
func Duration(name string) (Key, func(time.Duration) Field) {
	k := makeOrGetKey(name, IntKind)
	return k, func(v time.Duration) Field {
		return makeField2(k).SetDuration(v)
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

		for _, v2 := range v {
			infer(v2)
		}

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
// func Map(name string) func(v ...Field) Field {
// 	k := makeOrGetKey(name, MapKind)
// 	return func(v ...Field) Field {
// 		return makeField(k, v, true)
// 	}
// }

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
