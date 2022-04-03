package stmt

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/junhwong/goost/sqlx"
)

var stmts = map[string]*Statement{}

type Statement struct {
	stype  statementType
	query  string
	params sqlx.ParameterHolders
}

func (stmt *Statement) Type() statementType {
	return stmt.stype
}

// Deprecated: 重命名 Key
func Of(ctx context.Context, statementID string) context.Context {
	return context.WithValue(ctx, sqlx.StatementIDKey, statementID)
}

func Key(ctx context.Context, statementID string) context.Context {
	return context.WithValue(ctx, sqlx.StatementIDKey, statementID)
}

func (stmt *Statement) Exec(ctx context.Context, conn sqlx.ExecutableInterface, getter sqlx.ParameterGetter) (sqlx.ExecutedResult, error) {
	values, err := stmt.params.Values(getter)
	if err != nil {
		return nil, err
	}
	return conn.Exec(ctx, stmt.query, values)
}

func (stmt *Statement) Query(ctx context.Context, conn sqlx.QueryableInterface, getter sqlx.ParameterGetter) (sqlx.Rows, error) {
	values, err := stmt.params.Values(getter)
	if err != nil {
		return nil, err
	}
	return conn.Query(ctx, stmt.query, values)
}

func GetStatement(id string) *Statement {
	return stmts[id]
}

func getStatementWithCtx(ctx context.Context, stype statementType) (*Statement, error) {
	var id string
	if s, ok := ctx.Value(sqlx.StatementIDKey).(string); ok {
		id = s
	}
	stmt := GetStatement(id)
	if stmt == nil {
		return nil, fmt.Errorf("Statement %q undefined", id)
	}
	if stmt.stype != stype {
		return nil, fmt.Errorf("Statement %q defined, but mismatch type: %q != %q", id, stype, stmt.stype)
	}
	return stmt, nil
}

type RowInterface interface {
	Scan(dest ...interface{}) error
	ColumnTypes() ([]*sql.ColumnType, error)
	Columns() ([]string, error)
}

func Query(ctx context.Context, conn sqlx.QueryableInterface, getter sqlx.ParameterGetter,
	iterator func(row RowInterface) error) error {
	stmt, err := getStatementWithCtx(ctx, QueryStatement)
	if err != nil {
		return err
	}
	rows, err := stmt.Query(ctx, conn, getter)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		err = iterator(rows)
		if err != nil {
			break
		}
	}
	return err
}

func Exec(ctx context.Context, conn sqlx.ExecutableInterface, getter sqlx.ParameterGetter) (sqlx.ExecutedResult, error) {
	stmt, err := getStatementWithCtx(ctx, ExecStatement)
	if err != nil {
		return nil, err
	}
	return stmt.Exec(ctx, conn, getter)
}
