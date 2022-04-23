package stmt

import (
	"fmt"
	"strings"

	"github.com/junhwong/goost/sqlx/named"
)

type StructedInertBuildOption interface {
	applyInertBuilde(*insertBuilder)
}
type insertBuilder struct {
	buildOpts
	filterBuildOpts
}

func (f filterBuildApplyFn) applyInertBuilde(opt *insertBuilder) { f(&opt.filterBuildOpts) }
func (f buildApplyFn) applyInertBuilde(opt *insertBuilder)       { f(&opt.buildOpts) }

////////////////////

func BuildInsertSQL(obj interface{}, options ...StructedInertBuildOption) (*structedStmt, error) {
	opts := insertBuilder{}
	for _, opt := range options {
		if opt != nil {
			opt.applyInertBuilde(&opts)
		}
	}
	names, _, err := getNames(obj, opts.filters)
	if err != nil {
		return nil, err
	}
	if len(names) == 0 {
		return nil, fmt.Errorf("无可用字段")
	}

	opts.table, err = getTable(obj, opts.table)
	if err != nil {
		return nil, err
	}

	sets := []string{}
	vals := []string{}
	for name := range names {
		sets = append(sets, fmt.Sprintf("%q", name))
		vals = append(vals, fmt.Sprintf(":%s", name))
	}

	tpl := fmt.Sprintf("insert into %s (%s) VALUES (%s)", opts.table,
		strings.Join(sets, ","),
		strings.Join(vals, ","))
	sql, p, err := named.BuildNamedQuery(tpl, opts.namedBuildOptions...) //
	if err != nil {
		return nil, err
	}
	return &structedStmt{
		Statement: New(opts.stmtName, sql, p),
		names:     names,
	}, nil
}
