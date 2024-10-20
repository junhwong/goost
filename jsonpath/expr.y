%{
package jsonpath
%}


%union{
Expr                    Expr
FilterExpr              FilterExpr
MatcherExpr             MatcherExpr
op                      int
intVal                  int
bytes                   uint64
str                     string
falVal                  float64
duration                string
val                     string
exprType                int
strs                    []string
ValExpr                 any
EmptyGroup              any
BinaryExpr              *BinaryExpr
MemberExpr              MemberExpr
StringExpr              StringExpr
IndexExpr               IndexExpr
RootSymbol              Symbol
CurrentSymbol           Symbol
}

%start start
%type <MatcherExpr>  matchers
%type <Expr>         matcher
%type <Expr>         selector
%type <Expr>         valExpr
%type <Expr>         expr
%type <Expr>         segment
%type <Expr>         member
%type <Expr>         index
%type <op>           compairOp
%type <op>           includeOp
%type <Expr> range
%type <intVal> negint 
 
%token <str>     IDENTIFIER STRING 
%token <falVal>  FLOAT
%token <intVal>  INT
%left <op>       EQ NEQ GT GTE LT LTE RE NRE 
%left <op>       ADD SUB MUL DIV MOD POW 
%left <op>       AND OR OR2 AND2 IN NIN 
%left <op>       OPEN_BRACE CLOSE_BRACE OPEN_PARENTHESIS CLOSE_PARENTHESIS OPEN_BRACKET CLOSE_BRACKET 
%left <op>       DOT DOLLAR AT DOTDOT COMMA COLON SCOLON QUESTION

%%

start: 
  expr                          { yylex.(*parser).expr = $1 }
  ;

expr: 
  DOLLAR                                { $$ = RootSymbol }
  | DOLLAR index                        { $$ = &BinaryExpr{Left:RootSymbol, Right: $2, Op: OPEN_BRACKET} }
  | DOLLAR DOT segment                  { $$ = &BinaryExpr{Left:RootSymbol, Right: $3, Op: DOT} }
  | DOLLAR DOTDOT segment               { $$ = &BinaryExpr{Left:RootSymbol, Right: $3, Op: DOTDOT} }
  | AT                                  { $$ = CurrentSymbol }
  | AT index                            { $$ = &BinaryExpr{Left:CurrentSymbol, Right: $2, Op: OPEN_BRACKET} }
  | AT DOT segment                      { $$ = &BinaryExpr{Left:CurrentSymbol, Right: $3, Op: DOT} }
  | AT DOTDOT segment                   { $$ = &BinaryExpr{Left:CurrentSymbol, Right: $3, Op: DOTDOT} }
  | segment                             { $$ = $1 }
  ;

segment:
  member                        { $$ = $1 }
  | index                       { $$ = $1 }
  | selector                    { $$ = $1 }
  | MUL                         { $$ = WildcardSymbol }
  | segment index               { $$ = &BinaryExpr{Left:$1, Right: $2, Op: OPEN_BRACKET} }
  | segment selector            { $$ = &BinaryExpr{Left:$1, Right: $2, Op: OPEN_BRACKET} }
  | segment DOT segment         { $$ = &BinaryExpr{Left:$1, Right: $3, Op: DOT} }
  | segment DOTDOT segment      { $$ = &BinaryExpr{Left:$1, Right: $3, Op: DOTDOT} }
  // | segment DOT MUL             { $$ = &BinaryExpr{Left:$1, Right: WildcardSymbol, Op: DOT} }
  // | segment DOTDOT MUL          { $$ = &BinaryExpr{Left:$1, Right: WildcardSymbol, Op: DOTDOT} }
  ;

member:
  IDENTIFIER                    { $$ = MemberExpr($1) }
  | STRING                      { $$ = StringExpr($1) }
  // | DOT segment                 { $$ = &BinaryExpr{Left:$1, Right: $3, Op: DOT} }
  ;

index:
  OPEN_BRACKET negint CLOSE_BRACKET     { $$ = IndexExpr($2) }
  // 用于创建
  | OPEN_BRACKET CLOSE_BRACKET  { $$ = &EmptyGroup{Owner: nil} }
  ;

selector:
  OPEN_BRACKET range CLOSE_BRACKET      { $$ = $2 }
  | OPEN_BRACKET matchers CLOSE_BRACKET { $$ = $2 }
  ;

range:
  COLON                     { $$ = RangeExpr{0,-1} }
  | MUL                     { $$ = RangeExpr{0,-1} }
  | COLON negint            { $$ = RangeExpr{0,$2} }
  | negint COLON            { $$ = RangeExpr{$1,-1} }
  | negint COLON negint     { $$ = RangeExpr{$1,$3} }
  // | INT COLON INT COLON INT { $$ = IterExpr{$1,$3,$5} } // TODO: conflicts: 1 shift/reduce
  ;

matchers:
  matcher                   { $$ = MatcherExpr{ $1 } }
  | matchers COMMA matcher  { $$ = append($1, $3) }
  | matchers COMMA          { $$ = $1}
  ;

matcher:
  expr                                                                                             { $$ = $1}
  | QUESTION OPEN_PARENTHESIS expr compairOp valExpr CLOSE_PARENTHESIS                             { $$ = &FilterExpr{Body:&BinaryExpr{Left:$3, Right: $5, Op: $4}} } 
  | QUESTION OPEN_PARENTHESIS expr includeOp OPEN_BRACKET matchers CLOSE_BRACKET CLOSE_PARENTHESIS { $$ = &FilterExpr{Body:&BinaryExpr{Left:$3, Right: $6, Op: $4}} } 
  ;

valExpr:
  expr            { $$ = $1 }
  | FLOAT         { $$ = FloatValue($1) }
  | SUB FLOAT     { $$ = FloatValue(-$2) }
  | negint        { $$ = $1 }
  // | STRING        { $$ = $1 }
  ;

compairOp:
  EQ              { $$ = EQ }
  | NEQ           { $$ = NEQ }
  | GT            { $$ = GT }
  | GTE           { $$ = GTE }
  | LT            { $$ = LT }
  | LTE           { $$ = LTE }
  | RE            { $$ = RE }
  | NRE           { $$ = NRE }
  ;

includeOp:
  IN              { $$ = IN }
  | NIN           { $$ = NIN }
  ;
negint:
  INT             { $$ = $1 }
  | SUB INT       { $$ = -$2 }
  ;

%%