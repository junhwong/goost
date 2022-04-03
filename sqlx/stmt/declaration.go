package stmt

import (
	"fmt"
	"io/ioutil"
	"regexp"
	"strings"

	"github.com/junhwong/goost/sqlx"
	"gopkg.in/yaml.v3"
)

type statementType string

const (
	QueryStatement statementType = "Query"
	ExecStatement  statementType = "EXEC"
)

type Declaration struct {
	Type      string
	Table     string
	Columns   []interface{}
	Statement string
}

var (
	parameterExtract  = regexp.MustCompile(`\$\{[\w_\|\:]+\}|\$\:?[\w_]+`).ReplaceAllStringFunc
	removeBlankString = regexp.MustCompile(`(\r|\n|\s|\t)+`).ReplaceAllString
)

func Store(desc map[string]Declaration) error {
	for name, def := range desc {
		st := QueryStatement
		switch strings.ToLower(def.Type) {
		case "select":
		case "insert", "delete", "update":
			st = ExecStatement
		}
		holders := sqlx.ParameterHolders{}
		replacer := func(s string) string {
			if strings.HasPrefix(s, "${") {
				s = s[2 : len(s)-1]
				// TODO 处理管道符
			}
			s = strings.TrimPrefix(s, "$")
			if strings.HasPrefix(s, ":") {
				s = s[1:]
				// TODO 内置函数
				f, ok := functions[s]
				if !ok {
					panic(fmt.Errorf("sqlx.stmt: function not defined: %s", s))
				}
				holders = append(holders, func(getter sqlx.ParameterGetter) (interface{}, error) {
					return f()
				})
				return "?"
			}
			holders = append(holders, func(getter sqlx.ParameterGetter) (interface{}, error) {
				return getter.Get(s)
			})
			return "?"
		}
		sql := parameterExtract(def.Statement, replacer)
		stmts[name] = &Statement{
			stype:  st,
			query:  strings.TrimSpace(removeBlankString(removeBlankString(sql, " "), " ")),
			params: holders,
		}
	}
	return nil
}

func StoreWithMapper(f string) error {
	mapperData, err := ioutil.ReadFile(f) // todo
	if err != nil {
		return err
	}
	mapper := map[string]Declaration{}
	if err := yaml.Unmarshal(mapperData, &mapper); err != nil {
		return err
	}
	return Store(mapper)
}
