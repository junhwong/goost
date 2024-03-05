package jsonpath

import (
	"fmt"
	"strings"
	"testing"
	"text/scanner"
)

func TestId(t *testing.T) {
	testCases := []struct {
		desc string
		err  bool
	}{
		{
			desc: "abc",
		},
		{
			desc: "abc_def",
		},
		{
			desc: "gh-jk",
		},
		{
			desc: "中文",
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			lex := &lexer{
				Scanner: *(&scanner.Scanner{}).Init(strings.NewReader(tC.desc)),
				builder: strings.Builder{},
			}
			smt := &yySymType{}
			fmt.Printf("lex.Lex(smt): %v\n", lex.Lex(smt))
			fmt.Printf("smt.str: %v\n", smt.str)

		})
	}
}
