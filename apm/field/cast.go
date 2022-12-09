package field

import (
	"fmt"
	"reflect"

	"github.com/spf13/cast"
)

// ToInt64SliceE casts an interface to a []int64 type.
func ToInt64SliceE(i interface{}) ([]int64, error) {
	if i == nil {
		return nil, fmt.Errorf("unable to cast %#v of type %T to []int64", i, i)
	}

	switch v := i.(type) {
	case []int64:
		return v, nil
	}

	kind := reflect.TypeOf(i).Kind()
	switch kind {
	case reflect.Slice, reflect.Array:
		s := reflect.ValueOf(i)
		a := make([]int64, s.Len())
		for j := 0; j < s.Len(); j++ {
			val, err := cast.ToInt64E(s.Index(j).Interface())
			if err != nil {
				return nil, fmt.Errorf("unable to cast %#v of type %T to []int64", i, i)
			}
			a[j] = val
		}
		return a, nil
	default:
		return nil, fmt.Errorf("unable to cast %#v of type %T to []int64", i, i)
	}
}

// ToUint64SliceE casts an interface to a []int64 type.
func ToUint64SliceE(i interface{}) ([]uint64, error) {
	if i == nil {
		return nil, fmt.Errorf("unable to cast %#v of type %T to []uint64", i, i)
	}

	switch v := i.(type) {
	case []uint64:
		return v, nil
	}

	kind := reflect.TypeOf(i).Kind()
	switch kind {
	case reflect.Slice, reflect.Array:
		s := reflect.ValueOf(i)
		a := make([]uint64, s.Len())
		for j := 0; j < s.Len(); j++ {
			val, err := cast.ToUint64E(s.Index(j).Interface())
			if err != nil {
				return nil, fmt.Errorf("unable to cast %#v of type %T to []uint64", i, i)
			}
			a[j] = val
		}
		return a, nil
	default:
		return nil, fmt.Errorf("unable to cast %#v of type %T to []uint64", i, i)
	}
}

// ToFloat64SliceE casts an interface to a []float64 type.
func ToFloat64SliceE(i interface{}) ([]float64, error) {
	if i == nil {
		return nil, fmt.Errorf("unable to cast %#v of type %T to []float64", i, i)
	}

	switch v := i.(type) {
	case []float64:
		return v, nil
	}

	kind := reflect.TypeOf(i).Kind()
	switch kind {
	case reflect.Slice, reflect.Array:
		s := reflect.ValueOf(i)
		a := make([]float64, s.Len())
		for j := 0; j < s.Len(); j++ {
			val, err := cast.ToFloat64E(s.Index(j).Interface())
			if err != nil {
				return nil, fmt.Errorf("unable to cast %#v of type %T to []float64", i, i)
			}
			a[j] = val
		}
		return a, nil
	default:
		return nil, fmt.Errorf("unable to cast %#v of type %T to []float64", i, i)
	}
}
