package jsonpath

import (
	"fmt"
	"strings"
	"text/scanner"
)

type parser struct {
	*lexer
	expr  Expr
	funcs map[string]struct{}
}

func (l *parser) NewCallExpr(fn string, args []Expr) *CallExpr {
	if l.funcs == nil {
		l.funcs = make(map[string]struct{})
	}
	l.funcs[fn] = struct{}{}
	return &CallExpr{Func: fn, Args: args}
}
func (l *parser) GetCallFuncNames() []string {
	var a []string
	for k := range l.funcs {
		a = append(a, k)
	}
	return a
}

func Parse(s string) (Expr, Parsed, error) {
	lex := &lexer{
		Scanner: *(&scanner.Scanner{}).Init(strings.NewReader(s)),
		builder: strings.Builder{},
	}
	parser := &parser{
		lexer: lex,
	}
	r := yyParse(parser)

	if r != 0 {
		return nil, nil, fmt.Errorf("SyntaxError: %v", parser.errs)
	}
	return parser.expr, parser, nil

}

type Parsed interface {
	GetCallFuncNames() []string
}
