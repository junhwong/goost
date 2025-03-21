package field

import (
	"fmt"
	"net"
	"reflect"
	"strings"
	"sync"
	"time"

	"github.com/junhwong/goost/apm/field/loglevel"
	"github.com/spf13/cast"
)

var fieldPool = &sync.Pool{
	New: func() interface{} {
		return &Field{Schema: &Schema{}, Value: &Value{}}
	},
}

// 注意: 应该遵循谁构造(或接收)谁负责回收.
func Make(name string) *Field {
	f := fieldPool.Get().(*Field)
	f.Name = name
	f.NullValue = true
	return f
}

// Field 释放对象以备复用
func Release(fs ...*Field) {
	var ready []*Field
	var add func(f *Field)
	add = func(f *Field) {
		if f == nil {
			return
		}

		items := f.Items
		f.Items = nil
		for _, item := range items {
			add(item)
		}
		for _, it := range ready {
			if it == f {
				return
			}
		}
		ready = append(ready, f)
	}

	for _, f := range fs {
		add(f)
	}

	for _, f := range ready {
		f.Schema.Reset()
		f.Value.Reset()
		f.Parent = nil
		f.Index = 0
		fieldPool.Put(f)
	}
}

func MakeRoot() *Field {
	f := Make("")
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
	dst := Make(name)
	if v == nil {
		return dst
	}

	allow := func(iv any, k Type) bool { //
		b := k != InvalidKind
		if !b || len(allows) == 0 {
			return b
		}
		for _, t := range allows {
			if t == k {
				return true
			}
		}
		return false
	}

	iv, k := InferNumberValue(v)
	if k != InvalidKind {
		if iv != nil && k == FloatKind && !hasDecimal(iv.(float64)) {
			k = IntKind
			iv = int64(iv.(float64))
		}
		if !allow(iv, k) {
			return dst
		}
		return SetPrimitiveValue(dst, iv, k)
	}

	iv, k = InferPrimitiveValueWithoutNumber(v)
	if k != InvalidKind {
		if !allow(iv, k) {
			return dst
		}
		return SetPrimitiveValue(dst, iv, k)
	}

	var rv reflect.Value
	switch v := v.(type) {
	case []any:
		var fs []*Field
		for _, it := range v {
			fs = append(fs, Any("", it))
		}
		return dst.SetArray(fs)
	case map[string]any:
		if !allow(iv, GroupKind) {
			return dst
		}
		fs := []*Field{}
		for kk, vv := range v {
			it := Any(kk, vv)
			if it.Type == InvalidKind {
				continue
			}
			fs = append(fs, it)
		}
		dst.SetGroup(fs, false)
		return dst
	case reflect.Value:
		rv = v
	default:
		rv = reflect.ValueOf(v)
	}

	rt := rv.Type()
	prt := false
	if rt.Kind() == reflect.Pointer {
		rt = rt.Elem()
		rv = rv.Elem()
		prt = true
	}
	if rt.Kind() == reflect.Invalid {
		return dst
	}

	switch rv.Kind() {
	case reflect.Array, reflect.Slice:
		fs := []*Field{}
		var t Type
		same := false
		for i := range rv.Len() {
			it := Any("", rv.Index(i).Interface())
			fs = append(fs, it)
			if len(fs) > 0 && t != it.Type {
				same = true
			}
			t = it.Type
		}
		dst.SetArray(fs, same)
		return dst
	case reflect.Map:
		if !allow(iv, GroupKind) {
			return dst
		}
		fs := []*Field{}
		iter := rv.MapRange()
		for iter.Next() {
			kk, kt := InferPrimitiveValue(iter.Key())
			if kk == nil || kt == InvalidKind {
				continue
			}
			if !(kt == StringKind || kt == IntKind) {
				continue
			}

			it := Any(cast.ToString(kk), iter.Value().Interface())
			if it.Type == InvalidKind {
				continue
			}
			fs = append(fs, it)
		}
		dst.SetGroup(fs, false)
		return dst
	case reflect.Struct:
		panic("todo")
		return dst
	case reflect.Func, reflect.Chan:
		return dst
	default:
		iv, k = InferPrimitiveValueByReflect(rv)
		if k != InvalidKind {
			if !allow(iv, k) {
				return dst
			}
			return SetPrimitiveValue(dst, iv, k)
		}
	}
	if prt { // 创建默认值,
		iv, k = InferPrimitiveValueByReflect(reflect.Zero(rt))
	}
	if k == InvalidKind {
		iv, k = InferPrimitiveValueByReflect(rv)
	}

	if !allow(iv, k) {
		return dst
	}
	return SetPrimitiveValue(dst, iv, k)
}

func String(name string) (Key, func(string, ...interface{}) *Field) {
	k := makeOrGetKey(name, StringKind)
	return k, func(s string, a ...interface{}) *Field {
		v := s
		if s != "" && len(a) > 0 {
			v = fmt.Sprintf(s, a...)
		} else if len(a) > 0 {
			aa := make([]any, 0, len(a))
			for _, it := range a {
				if it != nil {
					aa = append(aa, it)
				}
			}
			v = fmt.Sprint(aa...)
		}
		v = strings.TrimSpace(v)
		return Make(name).SetString(v)
	}
}

func Bool(name string) (Key, func(bool) *Field) {
	k := makeOrGetKey(name, BoolKind)
	return k, func(v bool) *Field {
		return Make(name).SetBool(v)
	}
}

func Time(name string) (Key, func(time.Time) *Field) {
	k := makeOrGetKey(name, TimeKind)
	return k, func(v time.Time) *Field {
		return Make(name).SetTime(v)
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
		return Make(name).SetDuration(v)
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
