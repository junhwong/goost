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

		var r []*Field
		for _, f := range currentCopy {
			v.current = f.Items
			v.Visit(e.Right)
			if v.Error() != nil {
				return
			}
			r = append(r, v.current...)
		}
		v.current = r
	case jsonpath.DOTDOT:
		panic("todo")
	default:
		v.Err = fmt.Errorf("未定义的操作符:%q", jsonpath.GetOp(e.Op))
	}

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
		if f.Name == k { // todo 判断当前是否是素组元素
			tmp = append(tmp, f)
		}
	}

	if len(tmp) == 0 {
		if v.readonly {
			v.current = nil
			return
		}

		p := v.parent // todo
		if len(p) == 0 {
			p = v.parent
		}
		for _, f := range p {
			if f.Type == InvalidKind { // 新创建
				f.SetKind(GroupKind, false, false)
			}
			n := New(k)
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

func (v *explorer) VisitIndexExpr(e jsonpath.IndexExpr) {
	var tmp []*Field
	for _, it := range v.current {
		if !v.readonly && it.Type == InvalidKind {
			it.SetKind(ArrayKind, false, false)
		}
		if !it.IsArray() {
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
		if i != len(it.Items) || v.readonly {
			continue
		}
		f := New("")
		if err := it.Append(f); err != nil {
			v.SetError(err)
			return
		}
	}

	v.current = tmp
}

func (v *explorer) VisitEmptyGroup(e *jsonpath.EmptyGroup) {
	if v.readonly {
		v.current = nil
		return
	}

	v.Visit(e.Owner)
	if v.Error() != nil {
		return
	}
	var tmp []*Field
	for _, it := range v.current {
		if it.Type == InvalidKind {
			it.SetKind(ArrayKind, false, false)
			tmp = append(tmp, it)
		} else if it.IsArray() {
			tmp = append(tmp, it)
		}
	}
	var tmp2 []*Field
	for _, p := range tmp {
		if v.Error() != nil {
			return
		}
		f := New("")
		if err := p.Append(f); err != nil {
			v.SetError(err)
		}
		tmp2 = append(tmp2, f)
	}
	v.current = tmp2
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
			j += len(it.Items)
		}
		if i < 0 || i > len(it.Items) {
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
	panic("todo")
}
func (v *explorer) VisitFilterExpr(e *jsonpath.FilterExpr) {
	panic("todo")
	// v.inFilter = true
	// for i, v2 := range e {
	// 	if i != 0 {
	// 		v.Write([]byte{','})
	// 	}
	// 	v.Visit(v2)
	// }
	// v.Write([]byte{']'})
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
