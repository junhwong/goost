package conv

import "reflect"

type Value interface {
	Unwarp() interface{}
}

type Object struct {
	val  interface{}
	ref  *reflect.Value
	kind reflect.Kind
}

func Warp(v interface{}) Object {
	if v == nil {
		return Object{}
	}
	switch v := v.(type) {
	case reflect.Value:
		return Object{ref: &v}
	case Object:
		return v
	case *reflect.Value:
		return Object{ref: v}
	case *Object:
		return *v
	default:
		r := reflect.ValueOf(v)
		return Object{ref: &r, val: v}
	}
}

func WarpString(v string) Object {
	return Object{val: v, kind: reflect.String}
}
func WarpStringPointer(v *string) Object {
	return Object{val: *v, kind: reflect.String}
}
func (v Object) IsNil() bool {
	return v.ref == nil || !v.ref.IsValid() || v.ref.IsNil()
}

func (v Object) Unwarp() interface{} {
	if v.IsNil() {
		return nil
	}
	return v.ref.Interface()
}
