%{
package jsonpath
%}


%union{
Expr                    Expr
Exprs                   []Expr
FilterExpr              FilterExpr
MatcherExpr             MatcherExpr
op                      int
intVal                  int
bytes                   uint64
str                     string
falVal                  float64
boolVal                 bool
duration                string
val                     string
exprType                int
strs                    []string
ValExpr                 any
EmptyGroup              any
BinaryExpr              *BinaryExpr
CallExpr                *CallExpr
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
%type <CallExpr>     callExpr
%type <Exprs>        callArgsExpr
%type <Expr>         expr
%type <Expr>         parentExpr
%type <Expr>         segment
%type <Expr>         member
%type <Expr>         index
%type <op>           compairOp
%type <op>           includeOp
%type <Expr>         range
%type <intVal>       intExpr 
 
%token <str>     IDENTIFIER STRING 
%token <falVal>  FLOAT
%token <intVal>  INT
%token <boolVal> TRUE FALSE
%left <op>       EQ NEQ GT GTE LT LTE RE NRE 
%left <op>       ADD SUB MUL DIV MOD POW 
%left <op>       AND OR OR2 AND2 IN NIN 
%left <op>       OPEN_BRACE CLOSE_BRACE OPEN_PARENTHESIS CLOSE_PARENTHESIS OPEN_BRACKET CLOSE_BRACKET 
%left <op>       DOT DOLLAR AT DOTDOT COMMA COLON SCOLON QUESTION
%left <op>       NEXT_CALL NEXT_SELECT,NEXT_DOT_NOP

%%

start: 
  expr                          { yylex.(*parser).expr = $1 }
  ;

expr: 
  DOLLAR                                { $$ = RootSymbol }
  | DOLLAR index                        { $$ = &BinaryExpr{Left:RootSymbol, Right: $2, Op: NEXT_SELECT} }
  | DOLLAR DOT segment                  { $$ = &BinaryExpr{Left:RootSymbol, Right: $3, Op: DOT} }
  | DOLLAR DOTDOT segment               { $$ = &BinaryExpr{Left:RootSymbol, Right: $3, Op: DOTDOT} }
  | AT                                  { $$ = CurrentSymbol }
  | AT index                            { $$ = &BinaryExpr{Left:CurrentSymbol, Right: $2, Op: NEXT_SELECT} }
  | AT DOT segment                      { $$ = &BinaryExpr{Left:CurrentSymbol, Right: $3, Op: DOT} }
  | AT DOTDOT segment                   { $$ = &BinaryExpr{Left:CurrentSymbol, Right: $3, Op: DOTDOT} }
  | parentExpr                          { $$ = $1 }
  | index                               { $$ = $1 }
  | selector                            { $$ = $1 }
  | AT DOT callExpr                     { $$ = &BinaryExpr{Left:CurrentSymbol, Right: $3, Op: NEXT_CALL} }
  ;

parentExpr: 
  AT AT                               { $$ = ParentSymbol }
  | AT AT index                       { $$ = &BinaryExpr{Left:ParentSymbol, Right: $3, Op: NEXT_SELECT} }
  | AT AT DOT segment                 { $$ = &BinaryExpr{Left:ParentSymbol, Right: $4, Op: DOT} }
  | AT AT DOTDOT segment              { $$ = &BinaryExpr{Left:ParentSymbol, Right: $4, Op: DOTDOT} }
  | AT AT DOT parentExpr              { $$ = &BinaryExpr{Left:ParentSymbol, Right: $4, Op: NEXT_DOT_NOP} }
  | AT AT DOTDOT parentExpr           { $$ = &BinaryExpr{Left:ParentSymbol, Right: $4, Op: DOTDOT} }
  ;

segment:
  member                        { $$ = $1 }
  | index                       { $$ = $1 }
  | selector                    { $$ = $1 }
  | MUL                         { $$ = WildcardSymbol }
  | segment index               { $$ = &BinaryExpr{Left:$1, Right: $2, Op: NEXT_SELECT} }
  | segment selector            { $$ = &BinaryExpr{Left:$1, Right: $2, Op: NEXT_SELECT} }
  | segment DOT segment         { $$ = &BinaryExpr{Left:$1, Right: $3, Op: DOT} }
  | segment DOTDOT segment      { $$ = &BinaryExpr{Left:$1, Right: $3, Op: DOTDOT} }
  | segment DOT callExpr        { $$ = &BinaryExpr{Left:$1, Right: $3, Op: NEXT_CALL} }
  // | segment DOT MUL             { $$ = &BinaryExpr{Left:$1, Right: WildcardSymbol, Op: DOT} }
  // | segment DOTDOT MUL          { $$ = &BinaryExpr{Left:$1, Right: WildcardSymbol, Op: DOTDOT} }
  ;

member:
  IDENTIFIER                    { $$ = MemberExpr($1) }
  | STRING                      { $$ = StringExpr($1) }
  // | DOT segment                 { $$ = &BinaryExpr{Left:$1, Right: $3, Op: DOT} }
  ;

callExpr:
  IDENTIFIER OPEN_PARENTHESIS callArgsExpr CLOSE_PARENTHESIS { $$ = yylex.(*parser).NewCallExpr($1,  $3) }
  | IDENTIFIER OPEN_PARENTHESIS  CLOSE_PARENTHESIS           { $$ = yylex.(*parser).NewCallExpr($1,  nil) }
  ;

callArgsExpr:
  valExpr                         { $$ = []Expr {$1} }
  | callArgsExpr COMMA valExpr    { $$ = append($1, $3) }
  | callArgsExpr COMMA STRING     { $$ = append($1, $3) }
  | callArgsExpr COMMA TRUE       { $$ = append($1, $3) }
  | callArgsExpr COMMA FALSE      { $$ = append($1, $3) }
  ;

index:
  OPEN_BRACKET intExpr CLOSE_BRACKET     { $$ = IndexExpr($2) }
  // 用于创建新元素 []
  | OPEN_BRACKET CLOSE_BRACKET  { $$ = &EmptyGroup{Owner: nil} }
  ;

selector:
  OPEN_BRACKET range CLOSE_BRACKET      { $$ = $2 }
  | OPEN_BRACKET matchers CLOSE_BRACKET { $$ = $2 }
  ;

range:
  COLON                     { $$ = RangeExpr{0,-1} }
  | MUL                     { $$ = RangeExpr{0,-1} }
  | COLON intExpr            { $$ = RangeExpr{0,$2} }
  | intExpr COLON            { $$ = RangeExpr{$1,-1} }
  | intExpr COLON intExpr     { $$ = RangeExpr{$1,$3} }
  // | INT COLON INT COLON INT { $$ = IterExpr{$1,$3,$5} } // TODO: conflicts: 1 shift/reduce
  ;

matchers:
  matcher                   { $$ = MatcherExpr{ $1 } }
  | matchers COMMA matcher  { $$ = append($1, $3) }
  | matchers COMMA          { $$ = $1}
  ;

matcher:
  expr                                                                                             { $$ = $1}
  | member                                                                                         { $$ = $1 } // 适配 @[a,b]
  | QUESTION OPEN_PARENTHESIS expr compairOp valExpr CLOSE_PARENTHESIS                             { $$ = &FilterExpr{Body:&BinaryExpr{Left:$3, Right: $5, Op: $4}} } 
  | QUESTION OPEN_PARENTHESIS expr includeOp OPEN_BRACKET matchers CLOSE_BRACKET CLOSE_PARENTHESIS { $$ = &FilterExpr{Body:&BinaryExpr{Left:$3, Right: $6, Op: $4}} } 
  ;

valExpr:
  expr            { $$ = $1 }
  | FLOAT         { $$ = $1 }
  | SUB FLOAT     { $$ = float64(-$2) }
  | intExpr        { $$ = $1 }
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

// 整数表达式
intExpr:
  INT             { $$ = $1 }
  | SUB INT       { $$ = -$2 }
  ;

%%