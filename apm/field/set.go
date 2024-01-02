package field

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
)

// 字段集合
type FieldSet []*Field

func (x FieldSet) Len() int           { return len(x) }
func (x FieldSet) Less(i, j int) bool { return x[i].GetKey() < x[j].GetKey() } // 字典序
func (x FieldSet) Swap(i, j int)      { x[i], x[j] = x[j], x[i] }
func (x FieldSet) Sort() {
	if len(x) == 0 {
		return
	}
	sort.Sort(x)
}

func (fs *FieldSet) Set(f *Field) *Field {
	f, _ = fs.Put(f)
	return f
}
func (fs *FieldSet) Put(f *Field) (crt, old *Field) {
	crt = f
	i := fs.At(f.GetKey())
	if i < 0 {
		*fs = append(*fs, f)
		return
	}
	tmp := *fs
	old = tmp[i]
	tmp[i] = f
	return
}

func (fs FieldSet) Get(k string) *Field {
	for _, v := range fs {
		if v.GetKey() == k {
			return v
		}
	}
	return nil
}
func (fs FieldSet) At(k string) int {
	for i, v := range fs {
		if v.GetKey() == k {
			return i
		}
	}
	return -1
}

func (fs *FieldSet) Remove(k string) *Field {
	i := fs.At(k)
	if i < 0 {
		return nil
	}

	tmp := *fs
	f := tmp[i]
	n := len(tmp) - 1
	for j := i; j < n; j++ {
		tmp[j] = tmp[j+1]
	}
	*fs = tmp[:n]
	return f
}

// 清除重复
func (fs *FieldSet) Unique() FieldSet {
	if fs == nil {
		return nil
	}
	tmp := FieldSet{}
	for _, f := range *fs {
		tmp.Set(f)
	}
	*fs = tmp
	return tmp
}

func (fs FieldSet) GetDive(k string) *Field {

	for _, v := range fs {
		if v.GetKey() == k {
			return v
		}
	}
	return nil
}
func (fs FieldSet) getDive(k string) *Field {
	f := fs.Get(k)
	if f != nil {
		return f
	}

	for _, v := range fs {
		if v.GetKey() == k {
			return v
		}
	}
	return nil
}

const ()

func SplitPath(s string) (r []PathSegment, err error) {

	r, _, err = doSplitPath(s, false)

	return

}

func doSplitPath(s string, inp bool) (r []PathSegment, end int, err error) {

	start := 0
	for i := 0; i < len(s); i++ {
		switch s[i] {
		case '\'', '"':
			if i-start != 0 {
				return nil, -1, fmt.Errorf("字符串索引必须是独立段")
			}
			ss := s[i+1:]
			j := strings.IndexAny(ss, string([]byte{s[i]}))
			if j < 0 {
				return nil, -1, fmt.Errorf("字符串未结束")
			}
			k := ss[:j]
			r = append(r, PathSegment{S: k, I: -1, Q: true})
			i += j + 1
			start = i + 1
			if len(s) > start && s[start] == '.' {
				i++
			}
		case '.':
			r = append(r, PathSegment{S: s[start:i], I: -1})
			start = i + 1
		case '#':
			ss := s[i+1:]
			j := strings.IndexAny(ss, ".[#")
			if j < 0 {
				j = len(s) - i - 1
			}
			k := ss[:j]
			n, err := strconv.ParseUint(k, 10, 64)
			if err != nil {
				return nil, -1, err
			}
			r = append(r, PathSegment{S: k, I: int(n)})
			i += j
			start = i
		case '[':
			if k := s[start:i]; len(k) > 0 {
				r = append(r, PathSegment{S: k, I: -1})
			}
			cr, ci, err := doSplitPath(s[i+1:], true)
			if err != nil {
				return nil, -1, err
			}
			if ci <= 0 {
				return nil, -1, fmt.Errorf("未结束")
			}
			if len(cr) == 1 {
				r = append(r, cr[0])
			} else {
				r = append(r, PathSegment{C: cr, I: -1, P: true})
			}
			i += ci
		case ',':
			if !inp {
				return nil, -1, fmt.Errorf("只能出现在括号中")
			}
		case ']':
			if !inp {
				return nil, -1, fmt.Errorf("只能出现在括号中")
			}
			end = i + 1
			return
		}
	}

	if k := s[start:]; len(k) > 0 {
		r = append(r, PathSegment{S: k, I: -1})
	}

	return

}

func readToEnd(s string, end ...byte) {

	for i := 0; i < len(s); i++ {

	}
	strings.IndexAny(s, "")

}

type PathSegment struct {
	S string // 键
	I int    // 数字索引
	Q bool   // 引号
	P bool   // 是否解析过了数字
	C []PathSegment
}
