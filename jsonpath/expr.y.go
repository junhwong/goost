// Code generated by goyacc - DO NOT EDIT.

package jsonpath

import __yyfmt__ "fmt"

type yySymType struct {
	yys           int
	Expr          Expr
	FilterExpr    FilterExpr
	MatcherExpr   MatcherExpr
	op            int
	intVal        int
	bytes         uint64
	str           string
	falVal        float64
	duration      string
	val           string
	exprType      int
	strs          []string
	ValExpr       any
	EmptyGroup    any
	BinaryExpr    *BinaryExpr
	MemberExpr    MemberExpr
	StringExpr    StringExpr
	IndexExpr     IndexExpr
	RootSymbol    Symbol
	CurrentSymbol Symbol
}

type yyXError struct {
	state, xsym int
}

const (
	yyDefault         = 57384
	yyEofCode         = 57344
	ADD               = 57358
	AND               = 57364
	AND2              = 57367
	AT                = 57378
	CLOSE_BRACE       = 57371
	CLOSE_BRACKET     = 57375
	CLOSE_PARENTHESIS = 57373
	COLON             = 57381
	COMMA             = 57380
	DIV               = 57361
	DOLLAR            = 57377
	DOT               = 57376
	DOTDOT            = 57379
	EQ                = 57350
	FLOAT             = 57348
	GT                = 57352
	GTE               = 57353
	IDENTIFIER        = 57346
	IN                = 57368
	INT               = 57349
	LT                = 57354
	LTE               = 57355
	MOD               = 57362
	MUL               = 57360
	NEQ               = 57351
	NIN               = 57369
	NRE               = 57357
	OPEN_BRACE        = 57370
	OPEN_BRACKET      = 57374
	OPEN_PARENTHESIS  = 57372
	OR                = 57365
	OR2               = 57366
	POW               = 57363
	QUESTION          = 57383
	RE                = 57356
	SCOLON            = 57382
	STRING            = 57347
	SUB               = 57359
	yyErrCode         = 57345

	yyMaxDepth = 200
	yyTabOfs   = -48
)

var (
	yyPrec = map[int]int{
		EQ:                0,
		NEQ:               0,
		GT:                0,
		GTE:               0,
		LT:                0,
		LTE:               0,
		RE:                0,
		NRE:               0,
		ADD:               1,
		SUB:               1,
		MUL:               1,
		DIV:               1,
		MOD:               1,
		POW:               1,
		AND:               2,
		OR:                2,
		OR2:               2,
		AND2:              2,
		IN:                2,
		NIN:               2,
		OPEN_BRACE:        3,
		CLOSE_BRACE:       3,
		OPEN_PARENTHESIS:  3,
		CLOSE_PARENTHESIS: 3,
		OPEN_BRACKET:      3,
		CLOSE_BRACKET:     3,
		DOT:               4,
		DOLLAR:            4,
		AT:                4,
		DOTDOT:            4,
		COMMA:             4,
		COLON:             4,
		SCOLON:            4,
		QUESTION:          4,
	}

	yyXLAT = map[int]int{
		57374: 0,  // OPEN_BRACKET (42x)
		57375: 1,  // CLOSE_BRACKET (37x)
		57373: 2,  // CLOSE_PARENTHESIS (28x)
		57380: 3,  // COMMA (28x)
		57344: 4,  // $end (22x)
		57346: 5,  // IDENTIFIER (22x)
		57347: 6,  // STRING (22x)
		57350: 7,  // EQ (21x)
		57352: 8,  // GT (21x)
		57353: 9,  // GTE (21x)
		57368: 10, // IN (21x)
		57388: 11, // index (21x)
		57354: 12, // LT (21x)
		57355: 13, // LTE (21x)
		57351: 14, // NEQ (21x)
		57369: 15, // NIN (21x)
		57357: 16, // NRE (21x)
		57356: 17, // RE (21x)
		57395: 18, // selector (21x)
		57376: 19, // DOT (19x)
		57378: 20, // AT (15x)
		57377: 21, // DOLLAR (15x)
		57349: 22, // INT (15x)
		57379: 23, // DOTDOT (14x)
		57391: 24, // member (14x)
		57394: 25, // segment (14x)
		57359: 26, // SUB (13x)
		57348: 27, // FLOAT (10x)
		57386: 28, // expr (7x)
		57381: 29, // COLON (5x)
		57392: 30, // negint (5x)
		57389: 31, // matcher (4x)
		57383: 32, // QUESTION (4x)
		57390: 33, // matchers (3x)
		57360: 34, // MUL (2x)
		57393: 35, // range (2x)
		57385: 36, // compairOp (1x)
		57387: 37, // includeOp (1x)
		57372: 38, // OPEN_PARENTHESIS (1x)
		57396: 39, // start (1x)
		57397: 40, // valExpr (1x)
		57384: 41, // $default (0x)
		57358: 42, // ADD (0x)
		57364: 43, // AND (0x)
		57367: 44, // AND2 (0x)
		57371: 45, // CLOSE_BRACE (0x)
		57361: 46, // DIV (0x)
		57345: 47, // error (0x)
		57362: 48, // MOD (0x)
		57370: 49, // OPEN_BRACE (0x)
		57365: 50, // OR (0x)
		57366: 51, // OR2 (0x)
		57363: 52, // POW (0x)
		57382: 53, // SCOLON (0x)
	}

	yySymNames = []string{
		"OPEN_BRACKET",
		"CLOSE_BRACKET",
		"CLOSE_PARENTHESIS",
		"COMMA",
		"$end",
		"IDENTIFIER",
		"STRING",
		"EQ",
		"GT",
		"GTE",
		"IN",
		"index",
		"LT",
		"LTE",
		"NEQ",
		"NIN",
		"NRE",
		"RE",
		"selector",
		"DOT",
		"AT",
		"DOLLAR",
		"INT",
		"DOTDOT",
		"member",
		"segment",
		"SUB",
		"FLOAT",
		"expr",
		"COLON",
		"negint",
		"matcher",
		"QUESTION",
		"matchers",
		"MUL",
		"range",
		"compairOp",
		"includeOp",
		"OPEN_PARENTHESIS",
		"start",
		"valExpr",
		"$default",
		"ADD",
		"AND",
		"AND2",
		"CLOSE_BRACE",
		"DIV",
		"error",
		"MOD",
		"OPEN_BRACE",
		"OR",
		"OR2",
		"POW",
		"SCOLON",
	}

	yyTokenLiteralStrings = map[int]string{}

	yyReductions = map[int]struct{ xsym, components int }{
		0:  {0, 1},
		1:  {39, 1},
		2:  {28, 1},
		3:  {28, 3},
		4:  {28, 3},
		5:  {28, 3},
		6:  {28, 3},
		7:  {28, 3},
		8:  {28, 3},
		9:  {28, 3},
		10: {28, 3},
		11: {25, 1},
		12: {25, 1},
		13: {25, 1},
		14: {25, 2},
		15: {25, 2},
		16: {25, 3},
		17: {24, 1},
		18: {24, 1},
		19: {11, 3},
		20: {18, 3},
		21: {18, 3},
		22: {35, 1},
		23: {35, 2},
		24: {35, 2},
		25: {35, 3},
		26: {33, 1},
		27: {33, 3},
		28: {33, 2},
		29: {31, 1},
		30: {31, 6},
		31: {31, 8},
		32: {40, 1},
		33: {40, 1},
		34: {40, 2},
		35: {40, 1},
		36: {36, 1},
		37: {36, 1},
		38: {36, 1},
		39: {36, 1},
		40: {36, 1},
		41: {36, 1},
		42: {36, 1},
		43: {36, 1},
		44: {37, 1},
		45: {37, 1},
		46: {30, 1},
		47: {30, 2},
	}

	yyXErrors = map[yyXError]string{}

	yyParseTab = [74][]uint8{
		// 0
		{59, 5: 57, 58, 11: 55, 18: 56, 20: 53, 52, 24: 54, 51, 28: 50, 39: 49},
		{4: 48},
		{4: 47},
		{117, 46, 46, 46, 46, 7: 46, 46, 46, 46, 106, 46, 46, 46, 46, 46, 46, 107, 116, 23: 115},
		{19: 111, 23: 112},
		// 5
		{19: 103, 23: 104},
		{37, 37, 37, 37, 37, 7: 37, 37, 37, 37, 12: 37, 37, 37, 37, 37, 37, 19: 37, 23: 37},
		{36, 36, 36, 36, 36, 7: 36, 36, 36, 36, 12: 36, 36, 36, 36, 36, 36, 19: 36, 23: 36},
		{35, 35, 35, 35, 35, 7: 35, 35, 35, 35, 12: 35, 35, 35, 35, 35, 35, 19: 35, 23: 35},
		{31, 31, 31, 31, 31, 7: 31, 31, 31, 31, 12: 31, 31, 31, 31, 31, 31, 19: 31, 23: 31},
		// 10
		{30, 30, 30, 30, 30, 7: 30, 30, 30, 30, 12: 30, 30, 30, 30, 30, 30, 19: 30, 23: 30},
		{59, 5: 57, 58, 11: 55, 18: 56, 20: 53, 52, 67, 24: 54, 51, 68, 28: 65, 63, 60, 64, 66, 62, 35: 61},
		{1: 100, 29: 101},
		{1: 99},
		{1: 98, 3: 86},
		// 15
		{1: 26, 22: 67, 26: 68, 30: 97},
		{1: 22, 3: 22},
		{1: 19, 3: 19},
		{38: 70},
		{1: 2, 2, 29: 2},
		// 20
		{22: 69},
		{1: 1, 1, 29: 1},
		{59, 5: 57, 58, 11: 55, 18: 56, 20: 53, 52, 24: 54, 51, 28: 71},
		{7: 74, 76, 77, 82, 12: 78, 79, 75, 83, 81, 80, 36: 72, 73},
		{59, 5: 57, 58, 11: 55, 18: 56, 20: 53, 52, 67, 24: 54, 51, 93, 92, 91, 30: 94, 40: 90},
		// 25
		{84},
		{12, 5: 12, 12, 20: 12, 12, 12, 26: 12, 12},
		{11, 5: 11, 11, 20: 11, 11, 11, 26: 11, 11},
		{10, 5: 10, 10, 20: 10, 10, 10, 26: 10, 10},
		{9, 5: 9, 9, 20: 9, 9, 9, 26: 9, 9},
		// 30
		{8, 5: 8, 8, 20: 8, 8, 8, 26: 8, 8},
		{7, 5: 7, 7, 20: 7, 7, 7, 26: 7, 7},
		{6, 5: 6, 6, 20: 6, 6, 6, 26: 6, 6},
		{5, 5: 5, 5, 20: 5, 5, 5, 26: 5, 5},
		{4},
		// 35
		{3},
		{59, 5: 57, 58, 11: 55, 18: 56, 20: 53, 52, 24: 54, 51, 28: 65, 31: 64, 66, 85},
		{1: 87, 3: 86},
		{59, 20, 3: 20, 5: 57, 58, 11: 55, 18: 56, 20: 53, 52, 24: 54, 51, 28: 65, 31: 89, 66},
		{2: 88},
		// 40
		{1: 17, 3: 17},
		{1: 21, 3: 21},
		{2: 96},
		{2: 16},
		{2: 15},
		// 45
		{22: 69, 27: 95},
		{2: 13},
		{2: 14},
		{1: 18, 3: 18},
		{1: 25},
		// 50
		{27, 27, 27, 27, 27, 7: 27, 27, 27, 27, 12: 27, 27, 27, 27, 27, 27, 19: 27, 23: 27},
		{28, 28, 28, 28, 28, 7: 28, 28, 28, 28, 12: 28, 28, 28, 28, 28, 28, 19: 28, 23: 28},
		{29, 29, 29, 29, 29, 7: 29, 29, 29, 29, 12: 29, 29, 29, 29, 29, 29, 19: 29, 23: 29},
		{1: 24, 22: 67, 26: 68, 30: 102},
		{1: 23},
		// 55
		{59, 5: 57, 58, 11: 55, 18: 56, 24: 54, 110},
		{59, 5: 57, 58, 11: 55, 18: 56, 24: 54, 105},
		{59, 42, 42, 42, 42, 7: 42, 42, 42, 42, 106, 42, 42, 42, 42, 42, 42, 107, 108},
		{34, 34, 34, 34, 34, 7: 34, 34, 34, 34, 12: 34, 34, 34, 34, 34, 34, 19: 34, 23: 34},
		{33, 33, 33, 33, 33, 7: 33, 33, 33, 33, 12: 33, 33, 33, 33, 33, 33, 19: 33, 23: 33},
		// 60
		{59, 5: 57, 58, 11: 55, 18: 56, 24: 54, 109},
		{32, 32, 32, 32, 32, 7: 32, 32, 32, 32, 106, 32, 32, 32, 32, 32, 32, 107, 32, 23: 32},
		{59, 44, 44, 44, 44, 7: 44, 44, 44, 44, 106, 44, 44, 44, 44, 44, 44, 107, 108},
		{59, 5: 57, 58, 11: 55, 18: 56, 24: 54, 114},
		{59, 5: 57, 58, 11: 55, 18: 56, 24: 54, 113},
		// 65
		{59, 43, 43, 43, 43, 7: 43, 43, 43, 43, 106, 43, 43, 43, 43, 43, 43, 107, 108},
		{59, 45, 45, 45, 45, 7: 45, 45, 45, 45, 106, 45, 45, 45, 45, 45, 45, 107, 108},
		{59, 5: 57, 58, 11: 55, 18: 56, 24: 54, 120, 34: 121},
		{59, 5: 57, 58, 11: 55, 18: 56, 24: 54, 109, 34: 119},
		{59, 118, 5: 57, 58, 11: 55, 18: 56, 20: 53, 52, 67, 24: 54, 51, 68, 28: 65, 63, 60, 64, 66, 62, 35: 61},
		// 70
		{1: 38, 38, 38, 38, 7: 38, 38, 38, 38, 12: 38, 38, 38, 38, 38, 38},
		{1: 40, 40, 40, 40, 7: 40, 40, 40, 40, 12: 40, 40, 40, 40, 40, 40},
		{59, 41, 41, 41, 41, 7: 41, 41, 41, 41, 106, 41, 41, 41, 41, 41, 41, 107, 108},
		{1: 39, 39, 39, 39, 7: 39, 39, 39, 39, 12: 39, 39, 39, 39, 39, 39},
	}
)

var yyDebug = 0

type yyLexer interface {
	Lex(lval *yySymType) int
	Error(s string)
}

type yyLexerEx interface {
	yyLexer
	Reduced(rule, state int, lval *yySymType) bool
}

func yySymName(c int) (s string) {
	x, ok := yyXLAT[c]
	if ok {
		return yySymNames[x]
	}

	if c < 0x7f {
		return __yyfmt__.Sprintf("%q", c)
	}

	return __yyfmt__.Sprintf("%d", c)
}

func yylex1(yylex yyLexer, lval *yySymType) (n int) {
	n = yylex.Lex(lval)
	if n <= 0 {
		n = yyEofCode
	}
	if yyDebug >= 3 {
		__yyfmt__.Printf("\nlex %s(%#x %d), lval: %+v\n", yySymName(n), n, n, lval)
	}
	return n
}

func yyParse(yylex yyLexer) int {
	const yyError = 47

	yyEx, _ := yylex.(yyLexerEx)
	var yyn int
	var yylval yySymType
	var yyVAL yySymType
	yyS := make([]yySymType, 200)

	Nerrs := 0   /* number of errors */
	Errflag := 0 /* error recovery flag */
	yyerrok := func() {
		if yyDebug >= 2 {
			__yyfmt__.Printf("yyerrok()\n")
		}
		Errflag = 0
	}
	_ = yyerrok
	yystate := 0
	yychar := -1
	var yyxchar int
	var yyshift int
	yyp := -1
	goto yystack

ret0:
	return 0

ret1:
	return 1

yystack:
	/* put a state and value onto the stack */
	yyp++
	if yyp >= len(yyS) {
		nyys := make([]yySymType, len(yyS)*2)
		copy(nyys, yyS)
		yyS = nyys
	}
	yyS[yyp] = yyVAL
	yyS[yyp].yys = yystate

yynewstate:
	if yychar < 0 {
		yylval.yys = yystate
		yychar = yylex1(yylex, &yylval)
		var ok bool
		if yyxchar, ok = yyXLAT[yychar]; !ok {
			yyxchar = len(yySymNames) // > tab width
		}
	}
	if yyDebug >= 4 {
		var a []int
		for _, v := range yyS[:yyp+1] {
			a = append(a, v.yys)
		}
		__yyfmt__.Printf("state stack %v\n", a)
	}
	row := yyParseTab[yystate]
	yyn = 0
	if yyxchar < len(row) {
		if yyn = int(row[yyxchar]); yyn != 0 {
			yyn += yyTabOfs
		}
	}
	switch {
	case yyn > 0: // shift
		yychar = -1
		yyVAL = yylval
		yystate = yyn
		yyshift = yyn
		if yyDebug >= 2 {
			__yyfmt__.Printf("shift, and goto state %d\n", yystate)
		}
		if Errflag > 0 {
			Errflag--
		}
		goto yystack
	case yyn < 0: // reduce
	case yystate == 1: // accept
		if yyDebug >= 2 {
			__yyfmt__.Println("accept")
		}
		goto ret0
	}

	if yyn == 0 {
		/* error ... attempt to resume parsing */
		switch Errflag {
		case 0: /* brand new error */
			if yyDebug >= 1 {
				__yyfmt__.Printf("no action for %s in state %d\n", yySymName(yychar), yystate)
			}
			msg, ok := yyXErrors[yyXError{yystate, yyxchar}]
			if !ok {
				msg, ok = yyXErrors[yyXError{yystate, -1}]
			}
			if !ok && yyshift != 0 {
				msg, ok = yyXErrors[yyXError{yyshift, yyxchar}]
			}
			if !ok {
				msg, ok = yyXErrors[yyXError{yyshift, -1}]
			}
			if yychar > 0 {
				ls := yyTokenLiteralStrings[yychar]
				if ls == "" {
					ls = yySymName(yychar)
				}
				if ls != "" {
					switch {
					case msg == "":
						msg = __yyfmt__.Sprintf("unexpected %s", ls)
					default:
						msg = __yyfmt__.Sprintf("unexpected %s, %s", ls, msg)
					}
				}
			}
			if msg == "" {
				msg = "syntax error"
			}
			yylex.Error(msg)
			Nerrs++
			fallthrough

		case 1, 2: /* incompletely recovered error ... try again */
			Errflag = 3

			/* find a state where "error" is a legal shift action */
			for yyp >= 0 {
				row := yyParseTab[yyS[yyp].yys]
				if yyError < len(row) {
					yyn = int(row[yyError]) + yyTabOfs
					if yyn > 0 { // hit
						if yyDebug >= 2 {
							__yyfmt__.Printf("error recovery found error shift in state %d\n", yyS[yyp].yys)
						}
						yystate = yyn /* simulate a shift of "error" */
						goto yystack
					}
				}

				/* the current p has no shift on "error", pop stack */
				if yyDebug >= 2 {
					__yyfmt__.Printf("error recovery pops state %d\n", yyS[yyp].yys)
				}
				yyp--
			}
			/* there is no state on the stack with an error shift ... abort */
			if yyDebug >= 2 {
				__yyfmt__.Printf("error recovery failed\n")
			}
			goto ret1

		case 3: /* no shift yet; clobber input char */
			if yyDebug >= 2 {
				__yyfmt__.Printf("error recovery discards %s\n", yySymName(yychar))
			}
			if yychar == yyEofCode {
				goto ret1
			}

			yychar = -1
			goto yynewstate /* try again in the same state */
		}
	}

	r := -yyn
	x0 := yyReductions[r]
	x, n := x0.xsym, x0.components
	yypt := yyp
	_ = yypt // guard against "declared and not used"

	yyp -= n
	if yyp+1 >= len(yyS) {
		nyys := make([]yySymType, len(yyS)*2)
		copy(nyys, yyS)
		yyS = nyys
	}
	yyVAL = yyS[yyp+1]

	/* consult goto table to find next state */
	exState := yystate
	yystate = int(yyParseTab[yyS[yyp].yys][x]) + yyTabOfs
	/* reduction by production r */
	if yyDebug >= 2 {
		__yyfmt__.Printf("reduce using rule %v (%s), and goto state %d\n", r, yySymNames[x], yystate)
	}

	switch r {
	case 1:
		{
			yylex.(*parser).expr = yyS[yypt-0].Expr
		}
	case 2:
		{
			yyVAL.Expr = yyS[yypt-0].Expr
		}
	case 3:
		{
			yyVAL.Expr = &BinaryExpr{Left: RootSymbol, Right: yyS[yypt-0].Expr, Op: DOT}
		}
	case 4:
		{
			yyVAL.Expr = &BinaryExpr{Left: CurrentSymbol, Right: yyS[yypt-0].Expr, Op: DOT}
		}
	case 5:
		{
			yyVAL.Expr = &BinaryExpr{Left: RootSymbol, Right: yyS[yypt-0].Expr, Op: DOTDOT}
		}
	case 6:
		{
			yyVAL.Expr = &BinaryExpr{Left: CurrentSymbol, Right: yyS[yypt-0].Expr, Op: DOTDOT}
		}
	case 7:
		{
			yyVAL.Expr = &BinaryExpr{Left: yyS[yypt-2].Expr, Right: yyS[yypt-0].Expr, Op: DOTDOT}
		}
	case 8:
		{
			yyVAL.Expr = &BinaryExpr{Left: yyS[yypt-2].Expr, Right: WildcardSymbol, Op: DOT}
		}
	case 9:
		{
			yyVAL.Expr = &BinaryExpr{Left: yyS[yypt-2].Expr, Right: WildcardSymbol, Op: DOTDOT}
		}
	case 10:
		{
			yyVAL.Expr = &EmptyGroup{Owner: yyS[yypt-2].Expr}
		}
	case 11:
		{
			yyVAL.Expr = yyS[yypt-0].Expr
		}
	case 12:
		{
			yyVAL.Expr = yyS[yypt-0].Expr
		}
	case 13:
		{
			yyVAL.Expr = yyS[yypt-0].Expr
		}
	case 14:
		{
			yyVAL.Expr = &BinaryExpr{Left: yyS[yypt-1].Expr, Right: yyS[yypt-0].Expr, Op: OPEN_BRACKET}
		}
	case 15:
		{
			yyVAL.Expr = &BinaryExpr{Left: yyS[yypt-1].Expr, Right: yyS[yypt-0].Expr, Op: OPEN_BRACKET}
		}
	case 16:
		{
			yyVAL.Expr = &BinaryExpr{Left: yyS[yypt-2].Expr, Right: yyS[yypt-0].Expr, Op: DOT}
		}
	case 17:
		{
			yyVAL.Expr = MemberExpr(yyS[yypt-0].str)
		}
	case 18:
		{
			yyVAL.Expr = StringExpr(yyS[yypt-0].str)
		}
	case 19:
		{
			yyVAL.Expr = IndexExpr(yyS[yypt-1].intVal)
		}
	case 20:
		{
			yyVAL.Expr = yyS[yypt-1].Expr
		}
	case 21:
		{
			yyVAL.Expr = yyS[yypt-1].MatcherExpr
		}
	case 22:
		{
			yyVAL.Expr = RangeExpr{0, -1}
		}
	case 23:
		{
			yyVAL.Expr = RangeExpr{0, yyS[yypt-0].intVal}
		}
	case 24:
		{
			yyVAL.Expr = RangeExpr{yyS[yypt-1].intVal, -1}
		}
	case 25:
		{
			yyVAL.Expr = RangeExpr{yyS[yypt-2].intVal, yyS[yypt-0].intVal}
		}
	case 26:
		{
			yyVAL.MatcherExpr = MatcherExpr{yyS[yypt-0].Expr}
		}
	case 27:
		{
			yyVAL.MatcherExpr = append(yyS[yypt-2].MatcherExpr, yyS[yypt-0].Expr)
		}
	case 28:
		{
			yyVAL.MatcherExpr = yyS[yypt-1].MatcherExpr
		}
	case 29:
		{
			yyVAL.Expr = yyS[yypt-0].Expr
		}
	case 30:
		{
			yyVAL.Expr = &FilterExpr{Body: &BinaryExpr{Left: yyS[yypt-3].Expr, Right: yyS[yypt-1].Expr, Op: yyS[yypt-2].op}}
		}
	case 31:
		{
			yyVAL.Expr = &FilterExpr{Body: &BinaryExpr{Left: yyS[yypt-5].Expr, Right: yyS[yypt-2].MatcherExpr, Op: yyS[yypt-4].op}}
		}
	case 32:
		{
			yyVAL.Expr = yyS[yypt-0].Expr
		}
	case 33:
		{
			yyVAL.Expr = FloatValue(yyS[yypt-0].falVal)
		}
	case 34:
		{
			yyVAL.Expr = FloatValue(-yyS[yypt-0].falVal)
		}
	case 35:
		{
			yyVAL.Expr = yyS[yypt-0].intVal
		}
	case 36:
		{
			yyVAL.op = EQ
		}
	case 37:
		{
			yyVAL.op = NEQ
		}
	case 38:
		{
			yyVAL.op = GT
		}
	case 39:
		{
			yyVAL.op = GTE
		}
	case 40:
		{
			yyVAL.op = LT
		}
	case 41:
		{
			yyVAL.op = LTE
		}
	case 42:
		{
			yyVAL.op = RE
		}
	case 43:
		{
			yyVAL.op = NRE
		}
	case 44:
		{
			yyVAL.op = IN
		}
	case 45:
		{
			yyVAL.op = NIN
		}
	case 46:
		{
			yyVAL.intVal = yyS[yypt-0].intVal
		}
	case 47:
		{
			yyVAL.intVal = -yyS[yypt-0].intVal
		}

	}

	if yyEx != nil && yyEx.Reduced(r, exState, &yyVAL) {
		return -1
	}
	goto yystack /* stack new state and value */
}
