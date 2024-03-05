package jsonpath

import (
	"fmt"
	"strings"
	"text/scanner"
)

type parser struct {
	*lexer
	expr Expr
}

func Parse(s string) (Expr, error) {
	lex := &lexer{
		Scanner: *(&scanner.Scanner{}).Init(strings.NewReader(s)),
		builder: strings.Builder{},
	}
	parser := &parser{
		lexer: lex,
	}
	r := yyParse(parser)

	if r != 0 {
		return nil, fmt.Errorf("%v", parser.errs)
	}
	return parser.expr, nil

}
