package field

import (
	"fmt"
	"net"
	"reflect"
	"strings"
	"time"
)

// TODO 从池中获取或创建字段对象
func New(name string) *Field {
	return &Field{Key: name}
}
func Release(f *Field) {
	if f == nil {
		return
	}
	f.Key = ""
	f.Reset()
}

// 构造一个动态字段
func Any(name string, v any) *Field {
	f := New(name)
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
		s := iv.(string)
		if f.Key == "level" { // todo 更多可能的名称
			l := ParseLevel(s)
			if l != LevelUnset {
				return SetLevel(f, l)
			}
		}
		SetString(f, s)
	case IntKind:
		SetInt(f, iv.(int64))
	case UintKind:
		SetUint(f, iv.(uint64))
	case FloatKind:
		SetFloat(f, iv.(float64))
	case BoolKind:
		SetBool(f, iv.(bool))
	case TimeKind:
		SetTime(f, iv.(time.Time))
	case DurationKind:
		SetDuration(f, iv.(time.Duration))
	case IPKind:
		SetIP(f, iv.(net.IP))
	case BytesKind:
		SetBytes(f, iv.([]byte))
	}
	return f
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
		return SetString(New(name), v)
	}
}

func Bool(name string) (Key, func(bool) *Field) {
	k := makeOrGetKey(name, BoolKind)
	return k, func(v bool) *Field {
		return SetBool(New(name), v)
	}
}

func Time(name string) (Key, func(time.Time) *Field) {
	k := makeOrGetKey(name, TimeKind)
	return k, func(v time.Time) *Field {
		return SetTime(New(name), v)
	}
}

func Int(name string) (Key, func(interface{}) *Field) {
	k := makeOrGetKey(name, IntKind)
	return k, func(v interface{}) *Field {
		return Any(name, v)
	}
}

func Uint(name string) (Key, func(interface{}) *Field) {
	k := makeOrGetKey(name, UintKind)
	return k, func(v interface{}) *Field {
		return Any(name, v)
	}
}

func Float(name string) (Key, func(interface{}) *Field) {
	k := makeOrGetKey(name, FloatKind)
	return k, func(v interface{}) *Field {
		return Any(name, v)
	}
}

func Duration(name string) (Key, func(time.Duration) *Field) {
	k := makeOrGetKey(name, IntKind)
	return k, func(v time.Duration) *Field {
		return SetDuration(New(name), v)
	}
}
func BuildLevel(name string) (Key, func(interface{}) *Field) {
	k := makeOrGetKey(name, Type_LEVEL)
	return k, func(v interface{}) *Field {
		f := Any(name, v)
		var i int
		switch f.Type {
		case Type_UINT:
			i = int(f.GetUintValue())
		case Type_INT:
			i = int(f.GetIntValue())
		case Type_LEVEL:
			return f
		default:
			f.Type = Type_UNKNOWN // panic
			return f
		}
		return SetLevel(f, LevelFromInt(i))
	}
}
