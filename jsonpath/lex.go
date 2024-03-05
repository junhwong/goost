package jsonpath

import (
	"fmt"
	"strconv"
	"strings"
	"text/scanner"
	"unicode/utf8"
)

var tokens = map[string]int{
	"{": OPEN_BRACE,
	"}": CLOSE_BRACE,
	"(": OPEN_PARENTHESIS,
	")": CLOSE_PARENTHESIS,
	"[": OPEN_BRACKET,
	"]": CLOSE_BRACKET,
	",": COMMA,
	":": COLON,
	// ";":  SCOLON,
	".":  DOT,
	"..": DOTDOT,
	"==": EQ,
	"!=": NEQ,
	">":  GT,
	">=": GTE,
	"<":  LT,
	"<=": LTE,
	"=~": RE,
	"!~": NRE,
	// "|=": PIPE_EXACT,
	// "|~": PIPE_MATCH,
	// "|":  PIPE,

	// binops
	"||": OR,
	"&&": AND,

	// OpTypeUnless: UNLESS,
	"+": ADD,
	"-": SUB,
	"*": MUL,
	"/": DIV,
	"%": MOD,
	// "**": POW,
	"@": AT,
	"$": DOLLAR,
	"?": QUESTION,
}

var idTokens = map[string]int{
	"or":  OR2,
	"and": AND2,
	"in":  IN,
	"nin": NIN,
}

type lexer struct {
	scanner.Scanner
	errs    []error
	builder strings.Builder
}

// todo 这个方法是由问题的, 需要修复
func (l *lexer) Lex(lval *yySymType) int {
	r := l.Scan()

	switch r {
	case scanner.EOF:
		return 0
	case scanner.Int:
		s := l.TokenText()
		i, err := strconv.Atoi(s)
		if err != nil {
			panic(err)
		}
		lval.intVal = i
		return INT
	case scanner.Float:
		s := l.TokenText()
		lval.str = s
		return FLOAT
	case scanner.String, scanner.RawString:
		// var err error
		s := l.TokenText()
		if !utf8.ValidString(s) {
			l.Error("invalid UTF-8 rune")
			return 0
		}
		lval.str = s
		return STRING
	case scanner.Ident:
		tokenText := l.TokenText()

		if tok, ok := idTokens[tokenText]; ok {
			return tok
		}

		lval.str = tokenText
		return IDENTIFIER
	}

	tokenText := l.TokenText()
	if tok, ok := tokens[tokenText+string(l.Peek())]; ok {
		l.Next()
		return tok
	}
	if tok, ok := tokens[tokenText]; ok {
		return tok
	}
	l.Error("非预期的token: " + scanner.TokenString(r))
	return -1
}

func (l *lexer) Error(msg string) {
	l.errs = append(l.errs, fmt.Errorf("parser error: %v %v:%v", msg, l.Line, l.Column))
}

func GetOp(op int) string {
	return yySymNames[yyXLAT[op]]
}
