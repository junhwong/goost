package sql

import (
	"context"
	"database/sql"
	"fmt"
	"regexp"
	"strings"

	"github.com/junhwong/goost/apm"
	"github.com/junhwong/goost/pkg/field"
)

//VARBINARY(16)
// INET6_ATON INET6_NTOA INET_ATON INET_NTOA

type ParameterGetter interface {
	GetParameter(name string) interface{}
}

type Parameters []func(func(string) interface{}) (interface{}, error)

func (p Parameters) Values(getter ParameterGetter) ([]interface{}, error) {
	values := []interface{}{}
	for _, it := range p {
		val, err := it(getter.GetParameter)
		if err != nil {
			return nil, err
		}
		values = append(values, val)
	}
	return values, nil
}

// 预置、字段/表、参数
var parameterExtract = regexp.MustCompile(`(\#|\:)[\w_]+`).ReplaceAllStringFunc
var removeBlankString = regexp.MustCompile(`(\r|\n|\s|\t)+`).ReplaceAllString
var format = func(tpl string) string {
	tpl = removeBlankString(tpl, " ")
	tpl = removeBlankString(tpl, " ")
	return strings.TrimSpace(tpl)
}

func BuildSQL(tpl string) (sql string, params Parameters, err error) {
	params = Parameters{}
	replacer := func(s string) string {
		name := s[1:]
		switch s[:1] {
		case "#":
			return "`" + name + "`"
		case ":":
			params = append(params, func(get func(string) interface{}) (interface{}, error) {
				val := get(name)
				if val == nil {
					return nil, fmt.Errorf("parameter %q is required", name)
				}
				return val, nil
			})
			return "?"
		}
		return s
	}

	sql = format(tpl)
	sql = parameterExtract(sql, replacer)

	return

}

type ParameterHolder struct {
	params map[string]interface{}
}

func (h *ParameterHolder) GetParameter(name string) interface{} {
	return h.params[name]
}

func (h *ParameterHolder) Set(name string, value interface{}) {
	if value == nil {
		delete(h.params, name)
		return
	}
	h.params[name] = value
}
func MakeHolder() *ParameterHolder {
	return &ParameterHolder{params: make(map[string]interface{})}
}

type ExecutableInterface interface {
	Exec(ctx context.Context, query string, args ...interface{}) (ExecutedResult, error)
}

type QueryableInterface interface {
	Query(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
}

type ExecutedResult interface {
	sql.Result
}
type DBFactory interface {
	New() (*DBConn, error)
	Instance() *DBConn
	Close() error
}
type Configuration struct {
	Driver string `json:"driver" yml:"driver"`
	DSN    string `json:"dsn" yml:"dsn"`
}

func NewDBFactory(config Configuration) (DBFactory, error) {
	f := &dbFactory{
		driver: config.Driver,
		dsn:    config.DSN,
	}

	db, err := f.New()
	if err != nil {
		return nil, err
	}
	if err := db.db.Ping(); err != nil {
		return nil, err
	}
	db.isSingleton = true
	f.singleton = db
	return f, nil
}

type dbFactory struct {
	driver    string
	dsn       string
	singleton *DBConn
}

func (f *dbFactory) New() (*DBConn, error) {
	db, err := sql.Open(f.driver, f.dsn)
	if err != nil {
		return nil, err
	}
	// TODO: 设置db 参数
	// db.SetConnMaxLifetime(time.Minute * 3)
	// db.SetMaxOpenConns(1)
	// db.SetMaxIdleConns(1)
	return &DBConn{db: db}, nil
}

func (f *dbFactory) Instance() *DBConn {
	if f.singleton == nil {
		panic("db/sql: factory was closed")
	}
	return f.singleton
}
func (f *dbFactory) Close() error {
	//TODO
	return nil
}

type DBConn struct {
	db          *sql.DB
	isSingleton bool
}

func (c *DBConn) Close() error {
	if c == nil || c.db == nil || c.isSingleton {
		return nil
	}
	return c.db.Close()
}

var (
	statement = field.String("db.statement")
	arguments = field.Strings("db.arguments")
)

func (c *DBConn) Exec(ctx context.Context, query string, args ...interface{}) (ExecutedResult, error) {
	span := apm.Start(ctx,
		apm.WithName("sql"),
		apm.WithFields(
			statement(query),
			arguments(args...),
		),
	)
	defer span.Finish()

	stmt, err := c.db.Prepare(query)
	if err != nil {
		span.Fail()
		return nil, err
	}
	defer stmt.Close()

	result, err := stmt.Exec(args...)
	if err != nil {
		span.Fail()
	}
	return result, err
}

func (c *DBConn) Query(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	span := apm.Start(ctx,
		apm.WithName("sql"),
		apm.WithFields(
			statement(query),
			arguments(args...),
		),
	)
	defer span.Finish()

	stmt, err := c.db.Prepare(query)
	if err != nil {
		span.Fail()
		return nil, err
	}
	defer stmt.Close()

	result, err := stmt.Query(args...)
	if err != nil {
		span.Fail()
	}
	return result, err
}
