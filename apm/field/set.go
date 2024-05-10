package field

import (
	"github.com/junhwong/goost/jsonpath"
)

func Find(root *Field, nameOrPath string) ([]*Field, error) {
	seg, err := jsonpath.Parse(nameOrPath)
	if err != nil {
		return nil, err
	}
	return FindWith(seg, root)
}

func FindWith(p jsonpath.Expr, root *Field) ([]*Field, error) {
	v := &explorer{readonly: true, root: root, current: root.Items, parent: []*Field{root}}
	v.Visit = func(e jsonpath.Expr) {
		jsonpath.Visit(e, v, v.SetError)
	}
	v.Visit(p)
	return v.current, v.Error()
}

// 获取最后一个名称匹配的项
func GetLast(fs []*Field, name string) *Field {
	for i := len(fs) - 1; i >= 0; i-- {
		if fs[i].Name == name {
			return fs[i]
		}
	}
	return nil
}

func Get(fs []*Field, name string) (r []*Field) {
	for _, it := range fs {
		if it.Name == name {
			r = append(r, it)
		}
	}
	return
}

// 剔除与名称相符的项
func RidOf(fs []*Field, name string) (nf, r []*Field) {
	for _, it := range fs {
		if it.Name == name {
			r = append(r, it)
		} else {
			nf = append(nf, it)
		}
	}
	return
}

func Apply(expr jsonpath.Expr, root *Field, apply func(*Field)) error {
	v := &explorer{root: root, current: root.Items, parent: []*Field{root}}
	v.Visit = func(e jsonpath.Expr) {
		jsonpath.Visit(e, v, v.SetError)
	}
	v.Visit(expr)
	if err := v.Error(); err != nil {
		return err
	}
	for _, f := range v.current {
		apply(f)
	}
	return nil
}
