package jsonpath

type Expr interface {
}

type BinaryExpr struct {
	Left  Expr
	Right Expr
	Op    int
}
type FilterExpr struct {
	Body Expr
}
type MatcherExpr []Expr
type EmptyGroup struct {
	Owner Expr
}
type MemberExpr string
type StringExpr string
type IndexExpr int
type IterExpr [3]int
type RangeExpr [2]int
type FloatValue float64
type Symbol string

const (
	RootSymbol     Symbol = "$"
	CurrentSymbol  Symbol = "@"
	WildcardSymbol Symbol = "*"
)

// goyacc -l -o expr.y.go expr.y
// type Segment interface {
// 	Key() string
// 	Type() Kind
// }

// type Kind int

// const (
// 	IndexSegment Kind = iota + 1
// 	KeySegment
// 	QuoteSegment
// 	SymbolSegment
// 	RangeSegment
// 	PathSegment
// 	MulSegment
// )

// var kindMap = map[Kind]string{
// 	IndexSegment:  "i",
// 	KeySegment:    "f",
// 	QuoteSegment:  "q",
// 	SymbolSegment: "s",
// 	RangeSegment:  "r",
// 	MulSegment:    "m",
// }

// type Index int

// func (s Index) Type() Kind  { return IndexSegment }
// func (i Index) Key() string { return strconv.Itoa(int(i)) }
// func (i Index) Index() int  { return int(i) }

// type Key string

// func (s Key) Type() Kind  { return KeySegment }
// func (s Key) Key() string { return string(s) }

// type Quote string

// func (s Quote) Type() Kind  { return QuoteSegment }
// func (s Quote) Key() string { return string(s)[1 : len(s)-1] }

// type Symbol string

// func (s Symbol) Type() Kind  { return SymbolSegment }
// func (s Symbol) Key() string { return string(s) }

// const (
// 	RootSymbol     Symbol = "$"  // 根元素
// 	CurrentSymbol  Symbol = "@"  // 当前元素
// 	RDSymbol       Symbol = ".." // 递归下级, 表示任意子元素筛选,不能是结束。$..book[0,1]	前两本书
// 	WildcardSymbol Symbol = "*"  // 通配符
// 	childSymbol    Symbol = "."  // 成员
// 	colonSymbol    Symbol = ":"  // 范围
// )

// type Multiple []Segment

// func (s Multiple) Type() Kind  { return MulSegment }
// func (s Multiple) Key() string { return "" }

// type Path []Segment

// func (s Path) Type() Kind  { return PathSegment }
// func (s Path) Key() string { return "" }

// type Range [3]int

// func (s Range) Type() Kind     { return RangeSegment }
// func (s Range) Key() string    { return "" }
// func (s Range) ToSlice() []int { return s[:] }

// func String(i Segment) string {
// 	if i == nil {
// 		return ""
// 	}
// 	if i.Type() == MulSegment {
// 		var s string
// 		for i, v := range i.(Multiple) {
// 			if i != 0 {
// 				s += ","
// 			}
// 			s += String(v)
// 		}
// 		return fmt.Sprintf("%s:(%v)", kindMap[i.Type()], s)
// 	}
// 	return fmt.Sprintf("%s:%v", kindMap[i.Type()], i)
// }

// // 解析
// func Parse1(in string) (Segment, error) {
// 	s := []rune(in)
// 	var p Segment
// 	var r Path
// 	for len(s) != 0 {
// 		n, i, err := doParse(s)
// 		if err != nil {
// 			return nil, err
// 		}
// 		if p == nil {
// 			p = n
// 			s = s[i:]
// 			switch n.Type() {
// 			case SymbolSegment:
// 				if n.Key() == "." {
// 					continue
// 				}
// 			}
// 			r = append(r, n)
// 			continue
// 		}

// 		switch p.Type() {
// 		case SymbolSegment:
// 			switch n.Type() {
// 			case SymbolSegment:
// 				return nil, fmt.Errorf("非法的段:%q", s)
// 			}
// 			r = append(r, n)
// 		case QuoteSegment:
// 			switch n.Type() {
// 			case SymbolSegment:
// 				if n.Key() != "." {
// 					r = append(r, n)
// 				}
// 			default:
// 				return nil, fmt.Errorf("非法的段:%q", s)
// 			}
// 		default:
// 			switch n.Type() {
// 			case QuoteSegment:
// 				return nil, fmt.Errorf("非法的段2:%q %v", s, n.Type())
// 			case SymbolSegment:
// 				if n.Key() != "." {
// 					r = append(r, n)
// 				}
// 			default:
// 				r = append(r, n)
// 				// return nil, fmt.Errorf("非法的段2:%q %v", s, n.Type())
// 			}
// 		}
// 		p = n
// 		s = s[i:]
// 	}
// 	if len(r) == 1 {
// 		return r[0], nil
// 	}
// 	return r, nil

// }

// var (
// 	allSlice = Range{0, -1, 0}
// )

// // 解析 https://jsonpath.com/
// func Parse(in string) (Segment, error) {
// 	s := []rune(in)
// 	var r Path
// 	for len(s) != 0 {
// 		n, i, err := doParse(s)
// 		if err != nil {
// 			return nil, err
// 		}
// 		s = s[i:]
// 		r = append(r, n)
// 	}

// 	dive := true
// 	var r2 Path
// 	for _, s := range r {
// 		switch s.Type() {
// 		case SymbolSegment:
// 			switch s {
// 			case RootSymbol:
// 			}
// 		case MulSegment:
// 			s := s.(Multiple)
// 			parseItem := func(s Segment) (int, error) {
// 				if s.Type() == KeySegment {
// 					if x, err := strconv.ParseInt(s.Key(), 10, 64); err != nil { // 兼容处理, 非数字转为成员
// 						return int(x), nil
// 					}
// 				}
// 				return -1, fmt.Errorf("非法的段, 期待数字: %v", s)
// 			}

// 			if len(s) == 0 {
// 				panic("todo: 空数组")
// 			}

// 			var rc = 0
// 			for _, v := range s {
// 				if v == colonSymbol {
// 					rc++
// 				}
// 			}
// 			if rc > 0 && rc < 3 {
// 				switch len(s) {
// 				case 1: // [:]
// 					r2 = append(r2, Range{0, -1, -1})
// 				case 2: // [1:],[:1]
// 					var it Range
// 					var other Segment
// 					var idx int
// 					if s[0] == colonSymbol {
// 						it[0] = -1
// 						other = s[1]
// 						idx = 1
// 					} else if s[1] == colonSymbol {
// 						it[1] = -1
// 						other = s[0]
// 						idx = 0
// 					}
// 					i, err := parseItem(other)
// 					if err != nil {
// 						return nil, err
// 					}
// 					it[idx] = i
// 					r2 = append(r2, it)
// 				case 3: // [1:2]
// 					if s[1] != colonSymbol {
// 						return nil, fmt.Errorf("非法的段,语法错误: %v", s)
// 					}
// 					var it Range
// 					{
// 						i, err := parseItem(s[0])
// 						if err != nil {
// 							return nil, err
// 						}
// 						it[0] = i
// 					}
// 					{
// 						i, err := parseItem(s[2])
// 						if err != nil {
// 							return nil, err
// 						}
// 						it[2] = i
// 					}
// 					r2 = append(r2, it)
// 				case 5: // [1:2:3]
// 					if !(s[1] == colonSymbol && s[3] == colonSymbol) {
// 						return nil, fmt.Errorf("非法的段,语法错误: %v", s)
// 					}
// 					var it Range
// 					{
// 						i, err := parseItem(s[0])
// 						if err != nil {
// 							return nil, err
// 						}
// 						it[0] = i
// 					}
// 					{
// 						i, err := parseItem(s[2])
// 						if err != nil {
// 							return nil, err
// 						}
// 						it[1] = i
// 					}
// 					{
// 						i, err := parseItem(s[4])
// 						if err != nil {
// 							return nil, err
// 						}
// 						it[2] = i
// 					}
// 					r2 = append(r2, it)
// 				default:
// 					return nil, fmt.Errorf("非法的段,语法错误: %v", s)
// 				}
// 				continue
// 			} else if rc > 0 {
// 				return nil, fmt.Errorf("非法的段,语法错误: %v", s)
// 			}
// 			return s, nil // todo 过滤表达式
// 		}
// 		if !dive {

// 		}

// 	}

// 	if len(r) == 1 {
// 		return r[0], nil
// 	}
// 	return r, nil

// }

// func indexOf(s []rune, c rune) int {
// 	for i, v := range s {
// 		if v == c {
// 			return i
// 		}
// 	}
// 	return -1
// }

// func doParse(s []rune) (Segment, int, error) {
// 	if len(s) == 0 {
// 		return nil, 0, nil
// 	}
// 	c := s[0]
// 	switch c {
// 	case '\'', '"':
// 		j := indexOf(s[1:], c)
// 		if j < 0 {
// 			return nil, 1, fmt.Errorf("字符串未结束, 期待: %q", c)
// 		}
// 		j += 2
// 		return Quote(s[:j]), j, nil
// 	case '.', '@', '$', '*':
// 		for j := 0; j < len(s); j++ {
// 			if s[j] != c {
// 				switch string(s[:j]) {
// 				case RootSymbol.Key():
// 					return RootSymbol, j, nil
// 				case CurrentSymbol.Key():
// 					return CurrentSymbol, j, nil
// 				case WildcardSymbol.Key():
// 					return WildcardSymbol, j, nil
// 				case childSymbol.Key():
// 					return childSymbol, j, nil
// 				case RDSymbol.Key():
// 					return RDSymbol, j, nil
// 				}
// 				return nil, j, fmt.Errorf("未定义的符合: %q", s[:j])
// 			}
// 		}
// 	case '#': // 是[]的变体,只接受数字参数
// 		n, i, err := doParse(s[1:])
// 		if err != nil {
// 			return n, i, err
// 		}
// 		seg, ok := n.(Key)
// 		if !ok {
// 			return nil, 1, fmt.Errorf("不是有效的段, 期待: 数字")
// 		}
// 		x, err := strconv.ParseInt(seg.Key(), 10, 64)
// 		if err != nil {
// 			return nil, -1, err
// 		}
// 		return Index(x), i + 1, nil
// 	case '[':
// 		i := 1
// 		s = s[i:]
// 		isrange := false
// 		ismul := false
// 		var tmp Multiple

// 		for len(s) > 0 {
// 			switch s[0] {
// 			case ':':
// 				if ismul {
// 					return nil, i, fmt.Errorf("非法段,多参数不支持range: %v", s)
// 				}
// 				tmp = append(tmp, colonSymbol)
// 				i++
// 				s = s[1:]
// 				isrange = true
// 			case ',':
// 				if isrange {
// 					return nil, i, fmt.Errorf("非法段,range不支持多参数: %v", s)
// 				}
// 				i++
// 				s = s[1:]
// 				ismul = true
// 			case ']':
// 				i++
// 				return tmp, i, nil
// 			default:
// 				n, ni, err := doParse(s)
// 				if err != nil {
// 					return n, i + ni, err
// 				}
// 				tmp = append(tmp, n)
// 				s = s[ni:]
// 				i += ni
// 			}
// 		}
// 		return nil, -1, fmt.Errorf("不是有效的段, 期待: ]")

// 	case '?', '(', ')':
// 		return nil, -1, fmt.Errorf("暂时不支持表达式语义: %q", string(s))
// 	default:
// 		end := -1 //len(s)
// 	F:
// 		for i := 0; i < len(s); i++ {
// 			c := s[i]
// 			if !unicode.IsPrint(c) {
// 				break F
// 			}
// 			switch { // https://symbl.cc/cn/unicode/table/#linear-b-syllabary
// 			case c >= 0 && c <= 44 && c <= 0x002c:
// 				break F
// 			case c >= 46 && c <= 47:
// 				break F
// 			case c >= 58 && c <= 64:
// 				break F
// 			case c >= 91 && c <= 94:
// 				break F
// 			case c == 96:
// 				break F
// 			case c >= 123 && c <= 127:
// 				break F
// 			default:
// 				end = i
// 			}
// 		}

// 		if end != -1 {
// 			end++
// 			return Key(s[:end]), end, nil
// 		}
// 	}
// 	return nil, len(s), fmt.Errorf("非法的段: %q", s)
// }
