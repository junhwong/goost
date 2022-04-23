package stmt

import (
	"context"
	"database/sql"
	"sync"

	"github.com/junhwong/goost/apm"
	"github.com/junhwong/goost/sqlx"
)

// var stmts = map[string]*Statement{}
var (
	stmts = sync.Map{}
)

type statementType string
type Statement struct {
	stype  statementType
	name   string
	query  string
	params sqlx.ParameterHolders
}

func New(name, query string, params sqlx.ParameterHolders) *Statement {
	return &Statement{
		query:  query,
		params: params,
		name:   name,
	}
}

func (stmt *Statement) Type() statementType {
	return stmt.stype
}

func Of(ctx context.Context, statementID string) context.Context {
	return context.WithValue(ctx, sqlx.StatementIDKey, statementID)
}

func GetStatement(id string) *Statement {
	if v, ok := stmts.Load(id); ok {
		if s, _ := v.(*Statement); s != nil {
			return s
		}
	}
	return nil
}

type RowInterface interface {
	Err() error
	Scan(dest ...interface{}) error
	ColumnTypes() ([]*sql.ColumnType, error)
	Columns() ([]string, error)
}

type RowIter = func(row RowInterface) error

func (s *Statement) Query(ctx context.Context, raw sqlx.Conn, getter sqlx.ParameterGetter,
	iter RowIter, nextResultIter ...RowIter) error {
	ctx, span := apm.Start(ctx)
	defer span.End()

	stmt, err := raw.Prepare(ctx, s.query)
	if err != nil {
		return err
	}
	defer apm.Close(stmt, span)

	args, err := s.params.Values(getter)
	if err != nil {
		return err
	}
	rows, err := stmt.Query(ctx, args...)
	if err != nil {
		return err
	}
	defer rows.Close()

	for err == nil && rows.Next() {
		err = iter(rows)
	}
	var nextIter RowIter
	for _, it := range nextResultIter {
		nextIter = it
	}

	for err == nil && nextIter != nil && rows.NextResultSet() {
		for err == nil && rows.Next() {
			err = nextIter(rows)
		}
	}

	return err
}

func (s *Statement) Exec(ctx context.Context, raw sqlx.Conn, getter sqlx.ParameterGetter) (sqlx.ExecutedResult, error) {
	ctx, span := apm.Start(ctx,
		apm.WithFields(),
		apm.WithCallDepth(2),
	)
	defer span.End()
	stmt, err := raw.Prepare(ctx, s.query)
	if err != nil {
		return nil, err
	}
	defer apm.Close(stmt, span)

	args, err := s.params.Values(getter)
	if err != nil {
		return nil, err
	}
	return stmt.Exec(ctx, args...)
}

type structedStmt struct {
	*Statement
	names map[string]int
}

func (s *structedStmt) NewParams(obj interface{}, filters ...ParamterFilter) (*structedParams, error) {
	return NewStructedParams(obj, s.names, filters...)
}

type StructedStmt interface {
	Query(ctx context.Context, raw sqlx.Conn, getter sqlx.ParameterGetter, iter func(row RowInterface) error, nextResultIter ...func(row RowInterface) error) error
	Exec(ctx context.Context, raw sqlx.Conn, getter sqlx.ParameterGetter) (sql.Result, error)
	NewParams(obj interface{}, filters ...ParamterFilter) (*structedParams, error)
}
