package field

import (
	"fmt"

	"github.com/junhwong/goost/jsonpath"
)

var _ jsonpath.Visitor = (*explorer)(nil)

type explorer struct {
	// jsonpath.Base
	visit    func(jsonpath.Expr)
	root     *Field
	current  []*Field
	parent   []*Field
	readonly bool
	err      error
	getCall  CallFuncGetter
}

func (v *explorer) Error() error {
	return v.err
}
func (v *explorer) setError(err error) {
	if err == nil {
		panic("nil error")
	}
	v.err = err
}
func (v *explorer) VisitBinaryExpr(e *jsonpath.BinaryExpr) {
	v.visit(e.Left)
	switch e.Op {
	case jsonpath.NEXT_SELECT:
		v.visit(e.Right)
	case jsonpath.NEXT_CALL:
		v.visit(e.Right)
	case jsonpath.DOT:
		// 父级是group
		currentCopy := v.current
		if len(currentCopy) > 1 {
			currentCopy = make([]*Field, len(v.current))
			copy(currentCopy, v.current)
		}
		v.parent = currentCopy

		var tmp []*Field
		for _, f := range currentCopy {
			v.current = f.Items
			v.visit(e.Right)
			if v.Error() != nil {
				return
			}
			tmp = append(tmp, v.current...)
		}
		v.current = tmp
	case jsonpath.DOTDOT:
		if !v.readonly {
			v.setError(fmt.Errorf("..操作符不能写入"))
		}

		var tmp []*Field
		for _, it := range v.current {
			vsub := &explorer{readonly: true, root: v.root, current: readItems(it)}
			vsub.visit = func(e jsonpath.Expr) {
				jsonpath.Visit(e, vsub, vsub.setError)
			}
			vsub.visit(e.Right)
			if err := vsub.Error(); err != nil {
				v.setError(err)
				return
			}
			tmp = append(tmp, vsub.current...)
		}
		v.current = tmp
		return
	default:
		v.setError(fmt.Errorf("未定义的操作符:%q", jsonpath.GetOp(e.Op)))
	}

}

func readItems(f *Field) (r []*Field) {
	r = append(r, f)
	for _, it := range f.Items {
		r = append(r, readItems(it)...)
	}
	return
}

func (v *explorer) VisitSymbol(e jsonpath.Symbol) {
	switch e {
	case jsonpath.CurrentSymbol:
	case jsonpath.RootSymbol:
		v.current = []*Field{v.root}
	case jsonpath.ParentSymbol:
		var tmp []*Field
		for _, f := range v.current {
			if f.Parent != nil {
				tmp = append(tmp, f.Parent)
			}
		}
		v.current = tmp
	case jsonpath.WildcardSymbol:
	default:
		v.setError(fmt.Errorf("未定义的符合:%q", e))
	}
}

func (v *explorer) VisitMemberExpr(e jsonpath.MemberExpr) {
	k := string(e)
	var tmp []*Field
	for _, f := range v.current {
		if f.Name == k {
			tmp = append(tmp, f)
		}
	}
	if len(tmp) != 0 || v.readonly {
		v.current = tmp
		return
	}

	p := v.parent
	if len(p) == 0 {
		v.setError(fmt.Errorf("member创建成员必须有parent"))
		return
	}

	for _, g := range p {
		if g.Type == InvalidKind { // 新创建
			g.SetKind(GroupKind, false, false)
		}
		n := Make(k)
		g.Set(n)
		tmp = append(tmp, n)
	}
	v.current = tmp
}

func (v *explorer) VisitStringExpr(e jsonpath.StringExpr) {
	k := string(e)
	if len(k) <= 2 {
		v.current = nil
		return
	}
	k = k[1 : len(k)-1]
	v.visit(jsonpath.MemberExpr(k))
}

// 改变类型为数组
func (v *explorer) changeToArray(it *Field) bool {
	if it.IsArray() {
		return true
	}
	if v.readonly {
		return false
	}

	it.SetArray(nil) // todo 强制设置为数组,搞个开关
	return true
}

// 访问数组索引
func (v *explorer) VisitIndexExpr(e jsonpath.IndexExpr) {
	var tmp []*Field
	for _, it := range v.current {
		if !v.changeToArray(it) {
			continue
		}
		i := int(e)
		if i < 0 {
			i += len(it.Items)
		}
		if i >= 0 && i < len(it.Items) {
			tmp = append(tmp, it.Items[i])
			continue
		}
		if i == len(it.Items) && !v.readonly {
			// 刚好在数组尾部，则新增
			f := Make("")
			if !it.IsNull() {
				f.Type = it.Items[0].Type
			}
			it.Append(f)
			f.Type = InvalidKind
			tmp = append(tmp, f)
		}
	}

	v.current = tmp
}

func (v *explorer) VisitEmptyGroup(e *jsonpath.EmptyGroup) {
	// if e.Owner != nil {
	// 	v.Visit(e.Owner)
	// 	if v.Error() != nil {
	// 		return
	// 	}
	// }

	if v.readonly {
		v.setError(fmt.Errorf("EmptyArray just working a write mode"))
		return
	}
	// todo 写入模式
	panic("not implemented")
	var tmp []*Field
	for _, it := range v.current {
		if !it.IsArray() {
			it.SetArray(nil) // todo 强制设置为数组,搞个开关
			// tmp = append(tmp, it)
			// continue
		}
		f := Make("")
		if !it.IsNull() { // 临时处理 hack column类型验证
			f.Type = it.Items[0].Type
		}
		it.Append(f)
		f.Type = InvalidKind
		tmp = append(tmp, f)
	}
	v.current = tmp
}

func (v *explorer) VisitRangeExpr(e jsonpath.RangeExpr) {
	i := e[0]
	j := e[1]
	if i < 0 {
		i += len(v.current)
	}
	if j < 0 {
		j += len(v.current) + 1
	}
	if i < 0 || i >= len(v.current) {
		v.current = nil
		return
	}
	if j < 0 || j > len(v.current) {
		v.current = nil
		return
	}
	if i > j {
		v.current = nil
		return
	}
	v.current = v.current[i:j]
}

func (v *explorer) VisitIterExpr(e jsonpath.IterExpr) {
	panic("todo")
}
func (v *explorer) VisitMatcherExpr(e jsonpath.MatcherExpr) {
	if !v.readonly {
		v.setError(fmt.Errorf("匹配操作不能写入"))
		return
	}
	var tmp []*Field
	currentCopy := make([]*Field, len(v.current))
	copy(currentCopy, v.current)

	for _, it := range currentCopy {
		for _, exp := range e {
			v.current = []*Field{it}
			v.visit(exp)
			tmp = append(tmp, v.current...)
		}
	}
	v.current = tmp
}
func (v *explorer) VisitFilterExpr(e *jsonpath.FilterExpr) {
	var tmp []*Field
	currentCopy := make([]*Field, len(v.current))
	copy(currentCopy, v.current)

	for _, it := range currentCopy {
		v.current = []*Field{it}
		v.visit(e.Body)
		tmp = append(tmp, v.current...)
	}
	v.current = tmp
}
func (v *explorer) VisitConstInt(e int) {
	panic("todo")
}
func (v *explorer) VisitConstFloat(e float64) {
	panic("todo")
}
func (v *explorer) VisitConstString(e string) {
	panic("todo")
}
func (v *explorer) VisitConstBool(e bool) {
	panic("todo")
}
func (v *explorer) VisitCallExpr(e *jsonpath.CallExpr) {
	if !v.readonly {
		v.setError(fmt.Errorf("写入模式: 不能执行函数"))
		return
	}
	if v.getCall == nil {
		v.setError(fmt.Errorf("未设置函数调用器"))
		return
	}
	call, err := v.getCall(e.Func)
	if err != nil {
		v.setError(err)
		return
	}

	host := v.current
	var args []*Field
	for _, arg := range e.Args {
		v.visit(arg)
		args = append(args, v.current...)
	}
	if v.Error() != nil {
		return
	}
	v.current, err = call(host, args)
	if err != nil {
		v.setError(err)
		return
	}
}

type CallFunc func(host []*Field, args []*Field) ([]*Field, error)
type CallFuncGetter func(name string) (CallFunc, error)
