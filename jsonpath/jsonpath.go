package jsonpath

import (
	"fmt"
	"strconv"
)

type Segment interface {
	Key() string
	Type() Kind
}

type Kind int

const (
	IndexSegment Kind = iota + 1
	KeySegment
	QuoteSegment
	SymbolSegment
	RangeSegment
	PathSegment
	MulSegment
)

var kindMap = map[Kind]string{
	IndexSegment:  "i",
	KeySegment:    "f",
	QuoteSegment:  "q",
	SymbolSegment: "s",
	RangeSegment:  "r",
	MulSegment:    "m",
}

type Index int

func (s Index) Type() Kind  { return IndexSegment }
func (i Index) Key() string { return strconv.Itoa(int(i)) }
func (i Index) Index() int  { return int(i) }

type Key string

func (s Key) Type() Kind  { return KeySegment }
func (s Key) Key() string { return string(s) }

type Quote string

func (s Quote) Type() Kind  { return QuoteSegment }
func (s Quote) Key() string { return string(s)[1 : len(s)-1] }

type Symbol string

func (s Symbol) Type() Kind  { return SymbolSegment }
func (s Symbol) Key() string { return string(s) }

type Multiple []Segment

func (s Multiple) Type() Kind  { return MulSegment }
func (s Multiple) Key() string { return "" }

type Path []Segment

func (s Path) Type() Kind  { return PathSegment }
func (s Path) Key() string { return "" }

type Range [3]int

func (s Range) Type() Kind     { return RangeSegment }
func (s Range) Key() string    { return "" }
func (s Range) ToSlice() []int { return s[:] }

func String(i Segment) string {
	if i == nil {
		return ""
	}
	if i.Type() == MulSegment {
		var s string
		for i, v := range i.(Multiple) {
			if i != 0 {
				s += ","
			}
			s += String(v)
		}
		return fmt.Sprintf("%s:(%v)", kindMap[i.Type()], s)
	}
	return fmt.Sprintf("%s:%v", kindMap[i.Type()], i)
}

func Parse(in string) (Segment, error) {
	s := []rune(in)
	var p Segment
	var r Path
	for len(s) != 0 {
		n, i, err := doParse(s)
		if err != nil {
			return nil, err
		}
		if p == nil {
			p = n
			s = s[i:]
			switch n.Type() {
			case SymbolSegment:
				if n.Key() == "." {
					continue
				}
			}
			r = append(r, n)
			continue
		}

		switch p.Type() {
		case SymbolSegment:
			switch n.Type() {
			case SymbolSegment:
				return nil, fmt.Errorf("非法的段:%q", s)
			}
			r = append(r, n)
		case QuoteSegment:
			switch n.Type() {
			case SymbolSegment:
				if n.Key() != "." {
					r = append(r, n)
				}
			default:
				return nil, fmt.Errorf("非法的段:%q", s)
			}
		default:
			switch n.Type() {
			case QuoteSegment:
				return nil, fmt.Errorf("非法的段2:%q %v", s, n.Type())
			case SymbolSegment:
				if n.Key() != "." {
					r = append(r, n)
				}
			default:
				r = append(r, n)
				// return nil, fmt.Errorf("非法的段2:%q %v", s, n.Type())
			}
		}
		p = n
		s = s[i:]
	}
	if len(r) == 1 {
		return r[0], nil
	}
	return r, nil

}

func indexOf(s []rune, c rune) int {
	for i, v := range s {
		if v == c {
			return i
		}
	}
	return -1
}

func doParse(s []rune) (Segment, int, error) {
	if len(s) == 0 {
		return nil, 0, nil
	}
	c := s[0]
	switch c {
	case '\'', '"':

		j := indexOf(s[1:], c)
		if j < 0 {
			return nil, 1, fmt.Errorf("字符串未结束, 期待: %q", c)
		}
		j += 2
		return Quote(s[:j]), j, nil
	case '.', '@', '$', '*':
		for j := 0; j < len(s); j++ {
			if s[j] != c {
				return Symbol(s[:j]), j, nil
			}
		}
	case '#':
		n, i, err := doParse(s[1:])
		if err != nil {
			return n, i, err
		}
		seg, ok := n.(Key)
		if !ok {
			return nil, 1, fmt.Errorf("不是有效的段, 期待: 数字")
		}
		x, err := strconv.ParseInt(seg.Key(), 10, 64)
		if err != nil {
			return nil, -1, err
		}
		return Index(x), i + 1, nil
	case '[':
		i := 1
		s = s[i:]
		isrange := false
		ismul := false
		var tmp Multiple
		for len(s) != 0 {
			switch s[0] {
			case ':':
				if c := len(tmp); c == 0 || c == 3 {
					return nil, i, fmt.Errorf("非法段,range语法错误: %v", s)
				}
				if ismul {
					return nil, i, fmt.Errorf("非法段,多参数不支持range: %v", s)
				}
				isrange = true
				i++
				s = s[1:]
			case ',':
				i++
				s = s[1:]
				if isrange {
					return nil, i, fmt.Errorf("非法段,range不支持多参数: %v", s)
				}
				ismul = true
			case ']':
				i++
				if len(tmp) == 1 {
					n := tmp[0]
					seg, ok := n.(Key)
					if !ok { // todo 新类型?
						return n, i, nil
					}
					x, err := strconv.ParseInt(seg.Key(), 10, 64)
					if err != nil {
						return n, i, nil
					}
					return Index(x), i, nil
				}
				if isrange {
					r := Range{}
					for j, n := range tmp {
						switch n.Type() {
						case IndexSegment:
							r[j] = int(n.(Index))
						case KeySegment:
							x, err := strconv.ParseInt(n.Key(), 10, 64)
							if err != nil {
								return nil, -1, err
							}
							r[j] = int(x)
						default:
							return nil, i, fmt.Errorf("非法段,range不支持项: %v", n)
						}
					}
					return r, i, nil
				}
				return tmp, i, nil
			}
			n, ni, err := doParse(s)
			if err != nil {
				return n, i + ni, err
			}
			switch n.Type() {
			case KeySegment:
				tmp = append(tmp, n)
			default:
				return nil, i, fmt.Errorf("非法段: %v", n)
			}
			s = s[ni:]
			i += ni
		}
		return nil, 1, fmt.Errorf("不是有效的段, 期待: ]")
	case '?', '(', ')':
		panic("todo 解析表达式:" + string(s))
	default:
		end := -1 //len(s)
	F:
		for i := 0; i < len(s); i++ {
			c := s[i]
			switch { // https://symbl.cc/cn/unicode/table/#linear-b-syllabary
			case c >= 0 && c <= 44:
				break F
			case c >= 46 && c <= 47:
				break F
			case c >= 58 && c <= 64:
				break F
			case c >= 91 && c <= 94:
				break F
			case c == 96:
				break F
			case c >= 123 && c <= 127:
				break F
			default:
				end = i
			}
		}

		if end != -1 {
			end++
			return Key(s[:end]), end, nil
		}
	}
	return nil, len(s), fmt.Errorf("非法的段: %q", s)
}
