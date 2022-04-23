package named

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/junhwong/goost/sqlx"
)

var (
	parameterExtract  = regexp.MustCompile(`\:\w+|\{\:\w+\s*(\=\s*\w+)?(\s*\|\s*[a-zA-Z_]\w+)*\}|\{\@\w+\s*(\|\s*[a-zA-Z_]\w+)*\}`)
	removeBlankString = regexp.MustCompile(`(\r|\n|\s|\t)+`).ReplaceAllString
)

type BuildOption interface {
	apply(*buildOptions)
}
type buildOptions struct {
	holder func(i int) string
}
type holderFnOption func(i int) string

func (fn holderFnOption) apply(opt *buildOptions) {
	opt.holder = fn
}
func WithPostgreSQLPlaceholder() holderFnOption {
	return func(i int) string {
		return fmt.Sprintf("$%d", i)
	}
}

type buildError []error

func (err buildError) Error() string {
	arr := []string{}
	for _, e := range err {
		arr = append(arr, e.Error())
	}
	return fmt.Sprintf("sqlx: BuildNamedQueryError %s", strings.Join(arr, ";"))
}

type (
	PipeFn func(in interface{}) (out interface{}, err error)
)

var pipes = map[string]PipeFn{}

// TODO: 不同的驱动使用不同的参数占位符
func BuildNamedQuery(tpl string, options ...BuildOption) (sql string, holders sqlx.ParameterHolders, err error) {
	opts := buildOptions{
		holder: func(i int) string {
			return "?"
		},
	}
	for _, opt := range options {
		if opt != nil {
			opt.apply(&opts)
		}
	}
	holders = sqlx.ParameterHolders{}
	index := 0
	errs := buildError{}
	replacer := func(s string) string {
		s = strings.TrimPrefix(s, ":")
		iscall := false
		if strings.HasPrefix(s, "{") {
			s = s[1 : len(s)-1]
			// TODO 处理管道符
			if strings.HasPrefix(s, "@") { // 内置函数
				s = s[1:]
				iscall = true
			}
		}
		arr := strings.Split(s, "|")
		ps := []PipeFn{}
		for _, s := range arr[1:] {
			fn, ok := pipes[s]
			if !ok {
				errs = append(errs, fmt.Errorf("sqlx: pipe not defined: %s", s))
				continue
			}
			ps = append(ps, fn)
		}

		if iscall {
			f, ok := functions[arr[0]]
			if !ok {
				errs = append(errs, fmt.Errorf("sqlx: function not defined: %s", s))
			}
			holders = append(holders, func(getter sqlx.ParameterGetter) (out interface{}, err error) {
				out, err = f()
				for _, p := range ps {
					if err != nil {
						break
					}
					out, err = p(out)
				}
				return
			})
			index++
			return opts.holder(index)
		}
		s = arr[0]
		name := s
		if i := strings.Index(s, "="); i > 0 {
			name = s[:i]
			s = s[i+1:] // TODO: 默认值
		}
		holders = append(holders, func(getter sqlx.ParameterGetter) (out interface{}, err error) {
			out, err = getter.Get(name)
			// TODO: 默认值
			for _, p := range ps {
				if err != nil {
					break
				}
				out, err = p(out)
			}
			return
		})
		index++
		return opts.holder(index)
	}
	sql = parameterExtract.ReplaceAllStringFunc(tpl, replacer)
	if len(errs) != 0 {
		err = errs
	}
	sql = strings.TrimSpace(removeBlankString(removeBlankString(sql, " "), " "))
	return
}
