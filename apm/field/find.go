package field

import "github.com/junhwong/goost/jsonpath"

func Find(root *Field, nameOrPath string) ([]*Field, error) {
	seg, err := jsonpath.Parse(nameOrPath)
	if err != nil {
		return nil, err
	}
	return FindWith(seg, root)
}

func FindWith(p jsonpath.Expr, root *Field) ([]*Field, error) {
	v := &explorer{readonly: true, root: root, current: []*Field{root}, parent: []*Field{}}
	v.visit = func(e jsonpath.Expr) {
		jsonpath.Visit(e, v, v.setError)
	}
	v.visit(p)
	return v.current, v.Error()
}

func FindWithCurrent(p jsonpath.Expr, root *Field, current []*Field) ([]*Field, error) {
	v := &explorer{readonly: true, root: root, current: current, parent: []*Field{}}
	v.visit = func(e jsonpath.Expr) {
		jsonpath.Visit(e, v, v.setError)
	}
	v.visit(p)
	return v.current, v.Error()
}

func ParseNamePath(nameOrPath string) (jsonpath.Expr, error) {
	return jsonpath.Parse(nameOrPath)
}
