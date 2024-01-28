package field

import (
	"fmt"
	"net"
	"reflect"
	"strings"
	"time"

	"github.com/junhwong/goost/apm/field/loglevel"
)

// TODO 从池中获取或创建字段对象
func New(name string) *Field {
	return &Field{Schema: &Schema{Name: name}, Value: &Value{NullValue: true}}
}

func Release(f *Field) {
	if f == nil {
		return
	}
	f.Schema.Reset()
	f.Value.Reset()
}

func NewRoot() *Field {
	f := New("@")
	f.SetKind(GroupKind, false, false)
	f.SetNull(false)
	return f
}

// 构造一个动态字段
func SetPrimitiveValue(f *Field, v any, k Type) *Field {
	if v == nil || k == InvalidKind {
		f.resetValue()
		f.SetNull(true)
		f.Type = k
		return f
	}
	switch k {
	case StringKind:
		f.SetString(v.(string))
	case IntKind:
		f.SetInt(v.(int64))
	case UintKind:
		f.SetUint(v.(uint64))
	case FloatKind:
		f.SetFloat(v.(float64))
	case BoolKind:
		f.SetBool(v.(bool))
	case TimeKind:
		f.SetTime(v.(time.Time))
	case DurationKind:
		f.SetDuration(v.(time.Duration))
	case IPKind:
		f.SetIP(v.(net.IP))
	case BytesKind:
		f.SetBytes(v.([]byte))
	default:
		panic("todo")
	}
	return f
}

func Any(name string, v any, allows ...Type) *Field {
	iv, k := InferPrimitiveValue(v)
	if k == InvalidKind {
		iv, k = InferPrimitiveValueByReflect(reflect.ValueOf(v))
	}
	if k == InvalidKind {
		if err, _ := v.(error); err != nil {
			iv = err.Error()
			k = StringKind
		}
	}

	allow := k != InvalidKind
	if !allow {
		allows = nil
	}
	for _, t := range allows {
		if t == k {
			allow = true
			break
		}
	}
	f := New(name)
	if !allow {
		return f
	}
	if k == InvalidKind {
		return f
	}
	return SetPrimitiveValue(f, iv, k)
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
		return New(name).SetString(v)
	}
}

func Bool(name string) (Key, func(bool) *Field) {
	k := makeOrGetKey(name, BoolKind)
	return k, func(v bool) *Field {
		return New(name).SetBool(v)
	}
}

func Time(name string) (Key, func(time.Time) *Field) {
	k := makeOrGetKey(name, TimeKind)
	return k, func(v time.Time) *Field {
		return New(name).SetTime(v)
	}
}

func Int(name string) (Key, func(interface{}) *Field) {
	k := makeOrGetKey(name, IntKind)
	return k, func(v interface{}) *Field {
		return Any(name, v, k.Kind())
	}
}

func Uint(name string) (Key, func(interface{}) *Field) {
	k := makeOrGetKey(name, UintKind)
	return k, func(v interface{}) *Field {
		return Any(name, v, k.Kind())
	}
}

func Float(name string) (Key, func(interface{}) *Field) {
	k := makeOrGetKey(name, FloatKind)
	return k, func(v interface{}) *Field {
		return Any(name, v, k.Kind())
	}
}

// func Number(name string, allows ...Kind) func(interface{}) *Field {
// 	numTypes := []Kind{}
// 	if len(allows) == 0 {
// 		allows = numTypes
// 	} else {
// 		for _, t := range allows {
// 			p := false
// 			for _, t2 := range numTypes {
// 				if t2 == t {
// 					p = true
// 					break
// 				}
// 			}
// 			if !p {
// 				panic(fmt.Errorf("不是有效的数值类型: %v", t))
// 			}
// 		}
// 	}

// 	return func(v interface{}) *Field {
// 		v, k := InferNumberValue(v)
// 		if k == 0 {
// 		}
// 		return Any(name, v, allows...)
// 	}
// }

// var (
// 	kk = makeOrGetKey("name", DurationKind)
// 	mm = func(v any) *Field {
// 		return Any(kk.Name(), v, kk.Kind())
// 	}
// )

func Duration(name string) (Key, func(time.Duration) *Field) {
	k := makeOrGetKey(name, DurationKind)
	return k, func(v time.Duration) *Field {
		return New(name).SetDuration(v)
	}
}
func BuildLevel(name string) (Key, func(interface{}) *Field) {
	k := makeOrGetKey(name, LevelKind)
	return k, func(v interface{}) *Field {
		f := Any(name, v)
		var i int
		switch f.Type {
		case Type_UINT:
			i = int(f.GetUintValue())
		case Type_INT:
			i = int(f.GetIntValue())
		case LevelKind:
			return f
		default:
			f.Type = Type_UNKNOWN // panic
			return f
		}
		return f.SetLevel(loglevel.FromInt(i))
	}
}
