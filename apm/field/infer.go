package field

import (
	"net"
	"reflect"
	"time"

	"github.com/spf13/cast"
)

func InferPrimitiveValue(v any) (any, Kind) {
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
	case []byte:
		return v, BytesKind
	case net.IP:
		return v, IPKind
	}
	return nil, InvalidKind
}

// 反射获取重新定义基础类型的值
func InferPrimitiveValueByReflect(rv reflect.Value) (any, Kind) {
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
	return nil, InvalidKind
}
