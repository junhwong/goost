package field

import (
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/junhwong/goost/apm/field/pb"
	"github.com/spf13/cast"
)

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
	f := &pb.Field{Key: name}
	iv, k := InferPrimitiveValue(v)
	if k == InvalidKind {
		rv := reflect.ValueOf(v)
		if rv.Kind() == reflect.Pointer {
			rv = rv.Elem()
		}
		iv, k = InferPrimitiveValueByReflect(rv)
	}
	if k == InvalidKind {
		if err, _ := v.(error); err != nil {
			iv = err.Error()
			k = StringKind
		}
	}
	if k == InvalidKind {
		return f
	}

	switch k {
	case StringKind:
		v := iv.(string)
		f.Kind = StringKind
		f.StringValue = &v
	case IntKind:
		v := iv.(int64)
		f.Kind = IntKind
		f.IntValue = &v
	case UintKind:
		v := iv.(uint64)
		f.Kind = UintKind
		f.UintValue = &v
	case FloatKind:
		v := iv.(float64)
		f.Kind = FloatKind
		f.FloatValue = &v
	case BoolKind:
		v := iv.(bool)
		f.Kind = BoolKind
		f.BoolValue = &v
	case TimeKind:
		v := iv.(time.Time).UnixNano()
		f.Kind = TimeKind
		f.IntValue = &v
	case DurationKind:
		v := int64(iv.(time.Duration))
		f.Kind = DurationKind
		f.IntValue = &v
	}
	// f.KeyField = makeOrGetKey(name, f.kind)
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
		v = strings.TrimSpace(v)
		if len(v) == 0 {
			return &pb.Field{Key: name}
		}
		return &pb.Field{Key: name, Kind: StringKind, StringValue: &v}
	}
}

func Bool(name string) (Key, func(bool) Field) {
	k := makeOrGetKey(name, BoolKind)
	return k, func(v bool) Field {
		return &pb.Field{Key: name, Kind: BoolKind, BoolValue: &v}
	}
}

func Time(name string) (Key, func(time.Time) Field) {
	k := makeOrGetKey(name, TimeKind)
	return k, func(t time.Time) Field {
		v := t.UnixNano()
		return &pb.Field{Key: name, Kind: TimeKind, IntValue: &v}
	}
}

func Int(name string) (Key, func(interface{}) Field) {
	k := makeOrGetKey(name, IntKind)
	return k, func(v interface{}) Field {
		return Dynamic(name, v)
	}
}

func Uint(name string) (Key, func(interface{}) Field) {
	k := makeOrGetKey(name, UintKind)
	return k, func(v interface{}) Field {
		return Dynamic(name, v)
	}
}

func Float(name string) (Key, func(interface{}) Field) {
	k := makeOrGetKey(name, FloatKind)
	return k, func(v interface{}) Field {
		return Dynamic(name, v)
	}
}

// 时延。纳秒?
func Duration(name string) (Key, func(time.Duration) Field) {
	k := makeOrGetKey(name, IntKind)
	return k, func(v time.Duration) Field {
		d := int64(v)
		return &pb.Field{Key: name, Kind: DurationKind, IntValue: &d}
	}
}

// // Slice 返回一个数组对象
// func Slice(name string, kind KeyKind) (Key, func(...interface{}) Field) {
// 	k := &key{name: name, kind: kind}
// 	// dtype := DynamicKind
// 	// if len(kind) > 0 {
// 	// 	dtype = kind[len(kind)-1]
// 	// }
// 	return k, func(v ...interface{}) (r Field) {

// 		for _, v2 := range v {
// 			infer(v2)
// 		}

// 		switch kind {
// 		case StringKind:
// 			val, err := cast.ToStringSliceE(v)
// 			r = makeField(k, val, len(val) > 0, err)
// 		case IntKind:
// 			val, err := ToInt64SliceE(v)
// 			r = makeField(k, val, len(val) > 0, err)
// 		case UintKind:
// 			val, err := ToUint64SliceE(v)
// 			r = makeField(k, val, len(val) > 0, err)
// 		case FloatKind:
// 			val, err := ToFloat64SliceE(v)
// 			r = makeField(k, val, len(val) > 0, err)
// 		case BoolKind:
// 			val, err := cast.ToBoolSliceE(v)
// 			r = makeField(k, val, len(val) > 0, err)
// 		case TimeKind:
// 			// TODO: 时区未解决, 目前是UTC
// 			val := []time.Time{}
// 			var err error
// 			for _, it := range v {
// 				if t, ex := cast.ToTimeE(it); ex == nil {
// 					val = append(val, t)
// 				} else {
// 					err = ex
// 					break
// 				}
// 			}
// 			r = makeField(k, val, len(val) > 0, err)
// 		default:
// 			r = makeField(k, v, len(v) > 0, nil)
// 		}
// 		return
// 	}
// }

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
