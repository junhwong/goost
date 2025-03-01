package field

import "fmt"

var funcs = map[string]CallFunc{
	"len": callLen,
}

func getCall(fn string) (CallFunc, error) {
	if fn == "" {
		return nil, fmt.Errorf("func name is empty")
	}
	if f, ok := funcs[fn]; ok {
		return f, nil
	}
	return nil, fmt.Errorf("func %s not found", fn)
}

func callLen(host []*Field, args []*Field) ([]*Field, error) {
	return []*Field{Make("").SetInt(int64(len(host)))}, nil
}

func datetime_part(host []*Field, args []*Field) ([]*Field, error) {
	return []*Field{Make("").SetInt(int64(len(host)))}, nil
}
