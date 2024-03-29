package stmt

import (
	"fmt"
	"strings"

	"github.com/junhwong/goost/sqlx/named"
)

type UpdateBuildingOption interface {
	applyUpdateBuild(*updateBuildSettings)
}
type updateBuildSettings struct {
	buildOpts
	filterBuildOpts
	pkFieldsOption
	rcasVerBuildOpts
}

func (f buildApplyFn) applyUpdateBuild(opt *updateBuildSettings)          { f(&opt.buildOpts) }
func (f pkFieldsOptionApplyFn) applyUpdateBuild(opt *updateBuildSettings) { f(&opt.pkFieldsOption) }
func (f filterBuildApplyFn) applyUpdateBuild(opt *updateBuildSettings)    { f(&opt.filterBuildOpts) }
func (f rcasVerBuildApplyFn) applyUpdateBuild(opt *updateBuildSettings)   { f(&opt.rcasVerBuildOpts) }

func BuildUpdateSQL(obj interface{}, options ...UpdateBuildingOption) (*structedStmt, error) {
	settings := updateBuildSettings{}
	for _, opt := range options {
		if opt != nil {
			opt.applyUpdateBuild(&settings)
		}
	}

	names, _, err := getNames(obj, settings.filters)
	if err != nil {
		return nil, err
	}
	if len(names) == 0 {
		return nil, fmt.Errorf("无可用字段")
	}

	if len(settings.pkFields) == 0 {
		t, _ := obj.(PrimaryKeys)
		if t == nil {
			return nil, newParameterInvalidErr("未实现PrimaryKeys接口")
		}
		settings.pkFields = t.PrimaryKeyFields()
		if len(settings.pkFields) == 0 {
			return nil, newParameterInvalidErr("主键字段未提供")
		}
	}
	for _, k := range settings.pkFields {
		if _, ok := names[k]; !ok {
			return nil, newParameterInvalidErr("主键%q不在字段集中", k)
		}
	}

	settings.table, err = getTable(obj, settings.table)
	if err != nil {
		return nil, err
	}

	sets := []string{}
	cond := []string{}
	for name := range names {
		// 乐观锁
		if settings.casName == name {
			val := ":" + name
			if settings.casVal != "" {
				val = settings.casVal
			}
			cond = append(cond, fmt.Sprintf("%q=:%s", name, name))
			sets = append(sets, fmt.Sprintf("%q=%s", name, val))
			continue
		}
		if _, ok := names[name]; ok {
			cond = append(cond, fmt.Sprintf("%q=:%s", name, name))
			continue
		}
		sets = append(sets, fmt.Sprintf("%q=:%s", name, name))
	}

	if _, ok := names[settings.casName]; !ok && settings.casName != "" {
		return nil, newParameterInvalidErr("行乐观锁字段%q已定义, 但字段中不存在", settings.casName)
	}

	tpl := fmt.Sprintf("update %s set %s where %s", settings.table,
		strings.Join(sets, ","),
		strings.Join(cond, " and "))
	sql, p, err := named.BuildNamedQuery(tpl, named.WithPostgreSQLPlaceholder()) //
	if err != nil {
		return nil, err
	}
	return &structedStmt{
		Statement: New(settings.stmtName, sql, p),
		names:     names,
	}, nil
}
