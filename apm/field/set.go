package field

import (
	"github.com/junhwong/goost/jsonpath"
)

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
func RidOf(fs []*Field, name string) (others, taget []*Field) {
	for _, it := range fs {
		if it.Name == name {
			taget = append(taget, it)
		} else {
			others = append(others, it)
		}
	}
	return
}

func Apply(expr jsonpath.Expr, root *Field, apply func(*Field)) error {
	p := []*Field{root}
	// r := root.Items
	// if !root.IsCollection() { // todo 确认当前
	// 	r = p
	// }
	v := &explorer{root: root, current: p, parent: nil}
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
func ApplyWithCurrent(expr jsonpath.Expr, root *Field, current []*Field, apply func(*Field)) error {
	// r := root.Items
	// if !root.IsCollection() { // todo 确认当前
	// 	r = p
	// }
	v := &explorer{root: root, current: current, parent: nil}
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
