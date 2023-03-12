package sqlx

import (
	"context"
	"database/sql"
	"fmt"
	"io"

	"github.com/junhwong/goost/apm"
	"github.com/junhwong/goost/apm/field"
)

type Config struct {
	Name   string
	Driver string `json:"driver" yml:"driver"`
	DSN    string `json:"dsn" yml:"dsn"`
}

var (
	_, dbInstance  = field.String("db.instance")
	_, dbType      = field.String("db.type")
	_, dbStmt      = field.String("db.statement")
	_, dbStage     = field.String("db.stage")
	_, dbPrepareID = field.String("db.prepareid")
	_, dbArgs      = field.String("db.arguments")
	StatementIDKey = struct{ name string }{"$$statement_id"}
)

type sqlPrepare interface {
	PrepareContext(ctx context.Context, query string) (*sql.Stmt, error)
}
type sqlConn interface {
	io.Closer
	sqlPrepare
}

// 裁剪参数以使日志减小
func argSlim(a []interface{}) []interface{} {
	target := []interface{}{}
	for _, it := range a {
		// 大于100字节和30个字符
		if s, ok := it.(string); ok && len(s) > 100 {
			str := []rune(s)
			if len(str) > 30 {
				target = append(target, string(str[:10])+"..."+string(str[len(str)-10:]))
				continue
			}
		} else if _, ok := it.([]byte); ok {
			target = append(target, "<binary>")
			continue
		}

		target = append(target, it)
	}
	return target
}

func prepareContext(ctx context.Context, meta connMeta, raw sqlPrepare, query string) (*stmtWrap, error) {
	if raw == nil {
		return nil, fmt.Errorf("conn was closed")
	}
	// var id string
	// if s, ok := ctx.Value(StatementIDKey).(string); ok {
	// 	id = s
	// }
	// if id == "" {
	// 	id = "sql"
	// }
	panic("todo")
	ctx, span := apm.Start(ctx,
		// apm.WithFields(
		// 	dbType("sql"),
		// 	dbStage("prepare"),
		// 	dbInstance(meta.getInstance()),
		// 	dbStmt(query),
		// ),
		apm.WithCallDepth(3),
		// , apm.WithName(id), apm.WithCallDepth(3)
	)
	defer span.End()

	// fmt.Println("sql", query)
	stmt, err := raw.PrepareContext(ctx, query)
	if err != nil {
		span.Fail(err)
		return nil, err
	}

	return &stmtWrap{stmt: stmt, prepareid: span.Context().GetSpanID()}, nil
}

type stmtWrap struct {
	stmt      *sql.Stmt
	prepareid string
}

func (stmt *stmtWrap) Exec(ctx context.Context, args ...interface{}) (ExecutedResult, error) {
	panic("todo")
	ctx, span := apm.Start(ctx,
		// apm.WithFields(
		// 	dbType("sql"),
		// 	dbStage("exec"),
		// 	dbPrepareID(stmt.prepareid),
		// 	dbArgs(argSlim(args)...),
		// ),
		apm.WithCallDepth(3),
	)
	defer span.End()

	return stmt.stmt.ExecContext(ctx, args...)
}

func (stmt *stmtWrap) Query(ctx context.Context, args ...interface{}) (Rows, error) {
	panic("todo")
	ctx, span := apm.Start(ctx) // apm.WithFields(
	// 	dbType("sql"),
	// 	dbStage("query"),
	// 	dbPrepareID(stmt.prepareid),
	// 	dbArgs(argSlim(args)...),
	// ),

	defer span.End()

	return stmt.stmt.QueryContext(ctx, args...)
}

func (stmt *stmtWrap) QueryRow(ctx context.Context, args ...interface{}) (*sql.Row, error) {
	panic("todo")
	ctx, span := apm.Start(ctx) // apm.WithFields(
	// 	dbType("sql"),
	// 	dbStage("query"),
	// 	dbPrepareID(stmt.prepareid),
	// 	dbArgs(argSlim(args)...),
	// ),

	defer span.End()

	row := stmt.stmt.QueryRowContext(ctx, args...)
	return row, nil
}

func (stmt *stmtWrap) Close() error {
	return nil
}
