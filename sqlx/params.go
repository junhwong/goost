package sqlx

import "fmt"

type ParameterGetter interface {
	Get(key string) (interface{}, error)
}

type ParameterHolder func(ParameterGetter) (interface{}, error)

type ParameterHolders []ParameterHolder

func (ph ParameterHolders) Values(getter ParameterGetter) ([]interface{}, error) {
	values := make([]interface{}, 0)
	for _, holder := range ph {
		val, err := holder(getter)
		if err != nil {
			return nil, err
		}
		values = append(values, val)
	}
	return values, nil
}

type MapParameters map[string]interface{}

func (params MapParameters) Get(key string) (interface{}, error) {
	val, ok := params[key]
	if !ok {
		return nil, fmt.Errorf("Parameter %q not found", key)
	}
	return val, nil
}
