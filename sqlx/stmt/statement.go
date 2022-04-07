package stmt

import (
	"context"
	"database/sql"
	"fmt"
	"sync"

	"github.com/junhwong/goost/sqlx"
)

// var stmts = map[string]*Statement{}
var (
	stmts = sync.Map{}
)

type Statement struct {
	stype  statementType
	query  string
	params sqlx.ParameterHolders
}

func (stmt *Statement) Type() statementType {
	return stmt.stype
}

func Of(ctx context.Context, statementID string) context.Context {
	return context.WithValue(ctx, sqlx.StatementIDKey, statementID)
}

func (stmt *Statement) Exec(ctx context.Context, conn sqlx.ExecutableInterface, getter sqlx.ParameterGetter) (sqlx.ExecutedResult, error) {
	args, err := stmt.params.Values(getter)
	if err != nil {
		return nil, err
	}
	return conn.Exec(ctx, stmt.query, args)
}

func (stmt *Statement) Query(ctx context.Context, conn sqlx.QueryableInterface, getter sqlx.ParameterGetter) (sqlx.Rows, error) {
	args, err := stmt.params.Values(getter)
	if err != nil {
		return nil, err
	}
	return conn.Query(ctx, stmt.query, args)
}

func GetStatement(id string) *Statement {
	if v, ok := stmts.Load(id); ok {
		if s, _ := v.(*Statement); s != nil {
			return s
		}
	}

	return nil
}

func getStatementWithCtx(ctx context.Context, key string, stype statementType) (context.Context, *Statement, error) {
	// key, _ := ctx.Value(sqlx.StatementIDKey).(string)
	// if key == "" {
	// 	return ctx,nil, fmt.Errorf("stmt: Cannot get Statement key from context")
	// }
	stmt := GetStatement(key)
	if stmt == nil {
		return ctx, nil, fmt.Errorf("stmt: Statement %q undefined", key)
	}
	if stmt.stype != stype {
		return ctx, nil, fmt.Errorf("stmt: Statement %q defined, but mismatch type: %q != %q", key, stype, stmt.stype)
	}
	return Of(ctx, key), stmt, nil
}

type RowInterface interface {
	Scan(dest ...interface{}) error
	ColumnTypes() ([]*sql.ColumnType, error)
	Columns() ([]string, error)
}

type RowIter = func(row RowInterface) error

func Query(ctx context.Context, key string, conn sqlx.QueryableInterface, getter sqlx.ParameterGetter,
	iter RowIter, nextResultIterator ...RowIter) error {
	ctx, stmt, err := getStatementWithCtx(ctx, key, QueryStatement)
	if err != nil {
		return err
	}
	rows, err := stmt.Query(ctx, conn, getter)
	if err != nil {
		return err
	}
	defer rows.Close()

	for err == nil && rows.Next() {
		err = iter(rows)
	}

	var nextIter RowIter
	for _, it := range nextResultIterator {
		nextIter = it
	}

	for err == nil && nextIter != nil && rows.NextResultSet() {
		for err == nil && rows.Next() {
			err = nextIter(rows)
		}
	}

	return err
}

func Exec(ctx context.Context, key string, conn sqlx.ExecutableInterface, getter sqlx.ParameterGetter) (sqlx.ExecutedResult, error) {
	ctx, stmt, err := getStatementWithCtx(ctx, key, ExecStatement)
	if err != nil {
		return nil, err
	}
	return stmt.Exec(ctx, conn, getter)
}
