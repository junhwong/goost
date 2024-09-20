package field

import (
	"fmt"

	"github.com/junhwong/goost/jsonpath"
)

type explorer struct {
	jsonpath.Base
	Visit    func(jsonpath.Expr)
	root     *Field
	current  []*Field
	parent   []*Field
	readonly bool
}

func (v *explorer) VisitBinaryExpr(e *jsonpath.BinaryExpr) {
	v.Visit(e.Left)
	switch e.Op {
	case jsonpath.OPEN_BRACKET: // 跳过符合
		v.Visit(e.Right)
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
			v.Visit(e.Right)
			if v.Error() != nil {
				return
			}
			tmp = append(tmp, v.current...)
		}
		v.current = tmp
	case jsonpath.DOTDOT:
		if !v.readonly {
			v.SetError(fmt.Errorf("..操作符不能写入"))
		}

		var tmp []*Field
		for _, it := range v.current {
			vsub := &explorer{readonly: true, root: v.root, current: readItems(it)}
			vsub.Visit = func(e jsonpath.Expr) {
				jsonpath.Visit(e, vsub, vsub.SetError)
			}
			vsub.Visit(e.Right)
			if err := vsub.Error(); err != nil {
				v.SetError(err)
				return
			}
			tmp = append(tmp, vsub.current...)
		}
		v.current = tmp
		return
	default:
		v.Err = fmt.Errorf("未定义的操作符:%q", jsonpath.GetOp(e.Op))
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
	case jsonpath.RootSymbol:
		v.current = []*Field{v.root} //v.root.Items
	case jsonpath.CurrentSymbol:
	case jsonpath.WildcardSymbol:
	default:
		v.Err = fmt.Errorf("未定义的符合:%q", e)
	}
}

func (v *explorer) VisitMemberExpr(e jsonpath.MemberExpr) {
	k := string(e)
	var tmp []*Field
	for _, f := range v.current {
		// if f.Parent != nil && !f.Parent.IsGroup() {
		// 	v.SetError(fmt.Errorf("member访问必须是group"))
		// 	return
		// }
		if f.Name == k {
			tmp = append(tmp, f)
		}
	}

	if len(tmp) == 0 {
		if v.readonly {
			v.current = nil
			return
		}

		p := v.parent
		if len(p) == 0 {
			v.SetError(fmt.Errorf("member创建成员必须有parent"))
			return
		}

		for _, f := range p {
			if f.Type == InvalidKind { // 新创建
				f.SetKind(GroupKind, false, false)
			}
			n := Make(k)
			f.Set(n)
			tmp = append(tmp, n)
		}
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
	v.Visit(jsonpath.MemberExpr(k))
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
		// todo 超出索引是否创建
		// if i != len(it.Items) || v.readonly {
		// 	continue
		// }

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
		v.SetError(fmt.Errorf("EmptyArray just working a write mode"))
		return
	}

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
	var tmp []*Field
	for _, it := range v.current {
		if !it.IsArray() {
			continue
		}
		i := e[0]
		j := e[1]
		if i < 0 {
			i += len(it.Items)
		}
		if j < 0 {
			j += len(it.Items) + 1
		}
		if i < 0 || i >= len(it.Items) {
			continue
		}
		if j < 0 || j > len(it.Items) {
			continue
		}
		if i > j {
			continue
		}
		tmp = append(tmp, it.Items[i:j]...)
	}

	v.current = tmp
}
func (v *explorer) VisitIterExpr(e jsonpath.IterExpr) {
	panic("todo")
}
func (v *explorer) VisitMatcherExpr(e jsonpath.MatcherExpr) {
	if !v.readonly {
		v.SetError(fmt.Errorf("匹配操作不能写入"))
		return
	}
	var tmp []*Field
	currentCopy := make([]*Field, len(v.current))
	copy(currentCopy, v.current)

	for _, it := range currentCopy {
		for _, exp := range e {
			v.current = []*Field{it}
			v.Visit(exp)
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
		v.Visit(e.Body)
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
