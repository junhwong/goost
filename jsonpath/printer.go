package jsonpath

import (
	"fmt"
	"io"
	"strconv"
)

type printer struct {
	Base
	Out   io.Writer
	Visit func(Expr)
}

func (v *printer) Write(p []byte) {
	if v.Err != nil {
		return
	}
	_, v.Err = v.Out.Write(p)
}
func (v *printer) VisitBinaryExpr(e *BinaryExpr) {
	v.Visit(e.Left)
	switch e.Op {
	case OPEN_BRACKET: // 索引
	case DOT:
		v.Write([]byte{'.'})
	case DOTDOT:
		v.Write([]byte{'.', '.'})
	case LT:
		v.Write([]byte{'<'})
	case LTE:
		v.Write([]byte{'<', '='})
	case GT:
		v.Write([]byte{'>'})
	case GTE:
		v.Write([]byte{'>', '='})
	case EQ:
		v.Write([]byte{'=', '='})
	case NEQ:
		v.Write([]byte{'!', '='})
	default:
		v.Err = fmt.Errorf("未定义的操作符:%q", yySymNames[yyXLAT[e.Op]])
	}
	v.Visit(e.Right)
}

func (v *printer) VisitSymbol(e Symbol) {
	switch e {
	case RootSymbol, CurrentSymbol, WildcardSymbol:
		v.Write([]byte(e))
	default:
		v.Err = fmt.Errorf("未定义的符合:%q", e)
	}
}

func (v *printer) VisitMemberExpr(e MemberExpr) {
	v.Write([]byte(e))
}
func (v *printer) VisitStringExpr(e StringExpr) {
	v.Write([]byte(e))
}
func (v *printer) VisitIndexExpr(e IndexExpr) {
	v.Write([]byte{'['})
	v.Write([]byte(strconv.Itoa(int(e))))
	v.Write([]byte{']'})
}
func (v *printer) VisitEmptyGroup(e *EmptyGroup) {
	// v.Visit(e.Owner)
	v.Write([]byte{'['})
	v.Write([]byte{']'})
}
func (v *printer) VisitRangeExpr(e RangeExpr) {
	v.Write([]byte{'['})
	if e[0] != 0 {
		v.Write([]byte(strconv.Itoa(int(e[0]))))
	}
	v.Write([]byte{':'})
	if e[1] != -1 {
		v.Write([]byte(strconv.Itoa(int(e[1]))))
	}
	v.Write([]byte{']'})
}
func (v *printer) VisitIterExpr(e IterExpr) {
	v.Write([]byte{'['})
	v.Write([]byte(strconv.Itoa(int(e[0]))))
	v.Write([]byte{':'})
	v.Write([]byte(strconv.Itoa(int(e[1]))))
	v.Write([]byte{':'})
	v.Write([]byte(strconv.Itoa(int(e[2]))))
	v.Write([]byte{']'})
}
func (v *printer) VisitMatcherExpr(e MatcherExpr) {
	v.Write([]byte{'['})
	for i, v2 := range e {
		if i != 0 {
			v.Write([]byte{','})
		}
		v.Visit(v2)
	}
	v.Write([]byte{']'})
}

func (v *printer) VisitFilterExpr(e *FilterExpr) {
	v.Write([]byte{'?', '('})
	v.Visit(e.Body)
	v.Write([]byte{')'})
}
func (v *printer) VisitConstInt(e int) {
	v.Write([]byte(strconv.Itoa(int(e))))
}
func (v *printer) VisitConstFloat(e float64) {
	v.Write([]byte(strconv.FormatFloat(e, 'f', -1, 64)))
}
func (v *printer) VisitConstString(e string) {
	v.Write([]byte(e))
}
