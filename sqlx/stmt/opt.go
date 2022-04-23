package stmt

import (
	"strings"

	"github.com/junhwong/goost/sqlx/named"
)

type buildOpts struct {
	table             string
	stmtName          string
	namedBuildOptions []named.BuildOption
}

type buildApplyFn func(*buildOpts)

func WithTable(table string) buildApplyFn {
	if table == "" {
		panic("sqlx: table is required")
	}
	return func(opt *buildOpts) {
		opt.table = table
	}
}

////////////////////
type filterBuildOpts struct {
	filters []func(string) bool
}

type filterBuildApplyFn func(*filterBuildOpts)

func WithFilter(filters ...func(string) bool) filterBuildApplyFn {
	return func(opt *filterBuildOpts) {
		opt.filters = append(opt.filters, filters...)
	}
}

//////////////

type pkFieldsOption struct {
	pkFields []string
}
type pkFieldsOptionApplyFn func(*pkFieldsOption)

func WithPrimaryKeys(v ...string) pkFieldsOptionApplyFn {
	keys := []string{}
	for _, k := range v {
		k = strings.TrimSpace(k)
		if k != "" {
			keys = append(keys, k)
		}
	}
	return func(opt *pkFieldsOption) {
		opt.pkFields = append(opt.pkFields, keys...)
	}
}
