package jsonpath

import "fmt"

type Visitor interface {
	Error() error
	VisitBinaryExpr(e *BinaryExpr)
	VisitSymbol(e Symbol)
	VisitMemberExpr(e MemberExpr)
	VisitStringExpr(e StringExpr)
	VisitIndexExpr(e IndexExpr)
	VisitEmptyGroup(e *EmptyGroup)
	VisitIterExpr(e IterExpr)
	VisitRangeExpr(e RangeExpr)
	VisitFilterExpr(e *FilterExpr)
	VisitMatcherExpr(e MatcherExpr)
	VisitConstInt(e int)
	VisitConstFloat(e float64)
	VisitConstString(e string)
	VisitConstBool(e bool)
	VisitCallExpr(e *CallExpr)
}

func Visit(e Expr, v Visitor, setError func(err error)) {
	if v.Error() != nil {
		return
	}
	if e == nil {
		setError(fmt.Errorf("不能访问空表达式"))
		return
	}
	switch e := e.(type) {
	case *BinaryExpr:
		v.VisitBinaryExpr(e)
	case Symbol:
		v.VisitSymbol(e)
	case MemberExpr:
		v.VisitMemberExpr(e)
	case IndexExpr:
		v.VisitIndexExpr(e)
	case StringExpr:
		v.VisitStringExpr(e)
	case IterExpr:
		v.VisitIterExpr(e)
	case RangeExpr:
		v.VisitRangeExpr(e)
	case *EmptyGroup:
		v.VisitEmptyGroup(e)
	case *FilterExpr:
		v.VisitFilterExpr(e)
	case *CallExpr:
		v.VisitCallExpr(e)
	case MatcherExpr:
		v.VisitMatcherExpr(e)
	case int:
		v.VisitConstInt(e)
	case float64:
		v.VisitConstFloat(e)
	case bool:
		v.VisitConstBool(e)
	default:
		setError(fmt.Errorf("未定义的表达式类型, %T: %v", e, e))
	}
}

type Base struct {
	Err error
}

func (v *Base) Error() error {
	return v.Err
}
func (v *Base) SetError(err error) {
	v.Err = err
}

func (v *Base) VisitBinaryExpr(e *BinaryExpr) {
}

func (v *Base) VisitSymbol(e Symbol) {
}

func (v *Base) VisitMemberExpr(e MemberExpr) {
}
func (v *Base) VisitStringExpr(e StringExpr) {
}
func (v *Base) VisitIndexExpr(e IndexExpr) {
}
func (v *Base) VisitRangeExpr(e RangeExpr) {
}
func (v *Base) VisitIterExpr(e IterExpr) {
}
func (v *Base) VisitFilterExpr(e *FilterExpr) {
}
func (v *Base) VisitMatcherExpr(e MatcherExpr) {
}
func (v *Base) VisitConstInt(e int) {
}
func (v *Base) VisitConstFloat(e float64) {
}
func (v *Base) VisitConstString(e string) {
}
func (v *Base) VisitEmptyGroup(e *EmptyGroup) {}
func (v *Base) VisitCallExpr(e *CallExpr) {
}
