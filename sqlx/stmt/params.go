package stmt

import (
	"reflect"

	"github.com/junhwong/goost/apm"
	"github.com/junhwong/goost/runtime"
)

type ParamterFilter func(string, interface{}) (interface{}, error)

type structedParams struct {
	names   map[string]int
	val     reflect.Value
	filters []ParamterFilter
}

var (
	ParameterInvalidErr, newParameterInvalidErr = apm.NewErrorf("sqlx_param_invalid", 500,
		"Parameter is missing or invalid")
)

func NewStructedParams(obj interface{}, names map[string]int,
	filters ...ParamterFilter) (*structedParams, error) {
	v := reflect.ValueOf(obj)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	if v.Kind() != reflect.Struct {
		return nil, newParameterInvalidErr("Must be a struct instace")
	}
	return &structedParams{val: v, names: names, filters: filters}, nil
}
func (params *structedParams) Get(key string) (val interface{}, err error) {
	index, ok := params.names[key]
	if !ok {
		err = newParameterInvalidErr("%q undefined", key)
		return
	}
	defer runtime.HandleCrash(func(ex error) {

		err = ex
	})
	val = params.val.Field(index).Interface()
	for _, filter := range params.filters {
		if filter != nil {
			val, err = filter(key, val)
			if err != nil {
				return
			}
		}
	}

	return
}
