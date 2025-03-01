package field

import (
	"fmt"

	"github.com/junhwong/goost/jsonpath"
)

func Find(p jsonpath.Expr, root *Field, opts ...func(*explorer)) ([]*Field, error) {
	v := &explorer{readonly: true, root: root, current: []*Field{root}, parent: []*Field{}, getCall: getCall}
	v.visit = func(e jsonpath.Expr) {
		jsonpath.Visit(e, v, v.setError)
	}
	for _, opt := range opts {
		opt(v)
	}
	v.visit(p)
	return v.current, v.Error()
}

func ParseNamePathWith(nameOrPath string, getCall CallFuncGetter) (jsonpath.Expr, error) {
	seg, parsed, err := jsonpath.Parse(nameOrPath)
	if err != nil {
		return nil, err
	}
	funcs := parsed.GetCallFuncNames()
	if len(funcs) != 0 && getCall == nil {
		return nil, fmt.Errorf("funcs %v not found", funcs)
	}
	for _, fn := range funcs {
		if _, err := getCall(fn); err != nil {
			return nil, err
		}
	}
	return seg, nil
}
func ParseNamePath(nameOrPath string) (jsonpath.Expr, error) {
	return ParseNamePathWith(nameOrPath, nil)
}

func WithGetCallerFunc(getCall CallFuncGetter) func(exp *explorer) {
	return func(exp *explorer) {
		exp.getCall = getCall
	}
}
func WithCurrent(current []*Field) func(exp *explorer) {
	return func(exp *explorer) {
		exp.current = current
	}
}
