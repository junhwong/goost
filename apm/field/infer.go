package field

import (
	"net"
	"reflect"
	"time"

	"github.com/spf13/cast"
)

// 推导数值
func InferNumberValue(v any) (any, Type) {
	if v == nil {
		return v, InvalidKind
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
		if v == nil {
			return nil, IntKind
		}
		return int64(*v), IntKind
	case *int8:
		if v == nil {
			return nil, IntKind
		}
		return int64(*v), IntKind
	case *int16:
		if v == nil {
			return nil, IntKind
		}
		return int64(*v), IntKind
	case *int32:
		if v == nil {
			return nil, IntKind
		}
		return int64(*v), IntKind
	case *int64:
		if v == nil {
			return nil, IntKind
		}
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
		if v == nil {
			return nil, UintKind
		}
		return uint64(*v), UintKind
	case *uint8:
		if v == nil {
			return nil, UintKind
		}
		return uint64(*v), UintKind
	case *uint16:
		if v == nil {
			return nil, UintKind
		}
		return uint64(*v), UintKind
	case *uint32:
		if v == nil {
			return nil, UintKind
		}
		return uint64(*v), UintKind
	case *uint64:
		if v == nil {
			return nil, UintKind
		}
		return *v, UintKind

	case float32:
		return float64(v), FloatKind
	case float64:
		return v, FloatKind
	case *float32:
		if v == nil {
			return nil, FloatKind
		}
		return float64(*v), FloatKind
	case *float64:
		if v == nil {
			return nil, FloatKind
		}
		return *v, FloatKind
	}
	return v, InvalidKind
}

// 推导基本类型的值
func InferPrimitiveValue(v any) (any, Type) {
	if v == nil {
		return nil, InvalidKind
	}
	switch v := v.(type) {
	case bool:
		return v, BoolKind
	case *bool:
		if v == nil {
			return nil, BoolKind
		}
		return *v, BoolKind

	case string:
		return v, StringKind
	case *string:
		if v == nil {
			return nil, StringKind
		}
		return *v, StringKind

	case time.Time:
		return v, TimeKind
	case *time.Time:
		if v == nil {
			return nil, TimeKind
		}
		return *v, TimeKind

	case time.Duration:
		return v, DurationKind
	case *time.Duration:
		if v == nil {
			return nil, DurationKind
		}
		return *v, DurationKind
	case []byte:
		return v, BytesKind
	case net.IP:
		return v, IPKind
	case reflect.Value:
		return InferPrimitiveValueByReflect(v)
	}

	v, k := InferNumberValue(v)
	return v, k
	if k != InvalidKind {
		return v, k
	}

	rv := reflect.ValueOf(v)
	if v, k := InferPrimitiveValueByReflect(rv); k != InvalidKind {
		return v, k
	}

	rt := reflect.TypeOf(v)
	prt := false
	if rt.Kind() == reflect.Pointer {
		rt = rt.Elem()
		prt = true
	}
	if rt.Kind() == reflect.Invalid {
		return nil, InvalidKind
	}
	v, k = InferPrimitiveValueByReflect(reflect.Zero(rt)) // 创建默认值
	if prt {
		return nil, k
	}
	return v, k
}

// 反射获取重新定义基础类型的值
func InferPrimitiveValueByReflect(rv reflect.Value) (any, Type) {
	if rv.Kind() == reflect.Invalid {
		return rv, InvalidKind
	}
	rrv := rv
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
			n, _ := cast.ToTimeE(v)
			return n, TimeKind
		}
	}
	return rrv, InvalidKind
}
