package sqlx

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"sync"

	"github.com/junhwong/goost/apm"
	"github.com/junhwong/goost/pkg/field"
	"github.com/junhwong/goost/runtime"
)

type Config struct {
	Name   string
	Driver string `json:"driver" yml:"driver"`
	DSN    string `json:"dsn" yml:"dsn"`
}

type connWrap struct {
	serverVersion string
	name          string
	db            *sql.DB
	conn          *sql.Conn
	beginTxFn     func(ctx context.Context) (TransactionInterface, error)
	newConnFn     func(ctx context.Context) (ConnectionInterface, error)
}

func New(config Config) (DBInterface, error) {
	db, err := sql.Open(config.Driver, config.DSN)
	if err != nil {
		return nil, err
	}
	// db.Driver()
	err = db.Ping()
	if err != nil {
		return nil, err
	}

	if config.Name == "" {
		s := ""
		arr := strings.SplitN(config.DSN, "/", 2)
		if len(arr) == 2 {
			index := strings.Index(arr[0], "(")
			if index > 0 {
				s = arr[0][index+1 : len(arr[0])-1]
			} else {
				s = "localhost"
			}
			dbs := arr[1]
			index = strings.Index(dbs, "?")
			if index > 0 {
				dbs = dbs[:index]
			}
			s = s + "/" + dbs
		} else {
			s = "unnamed"
		}
		config.Name = s
	}

	conn := &connWrap{
		db:   db,
		name: config.Driver + "://" + config.Name,
	}

	if rows, err := db.Query("SELECT VERSION()"); err != nil {
		return nil, err
	} else {
		defer rows.Close()
		for rows.Next() {
			if err := rows.Scan(&conn.serverVersion); err != nil {
				return nil, err
			}
		}
	}

	conn.beginTxFn = func(ctx context.Context) (TransactionInterface, error) {
		if conn.db == nil {
			return nil, fmt.Errorf("conn was closed")
		}
		tx, err := conn.db.Begin()
		if err != nil {
			return nil, err
		}
		return &txWarp{
			name: conn.name,
			tx:   tx,
		}, nil
	}
	conn.newConnFn = func(ctx context.Context) (ConnectionInterface, error) {
		if conn.db == nil {
			return nil, fmt.Errorf("conn was closed")
		}
		nc, err := conn.db.Conn(ctx)
		if err != nil {
			return nil, err
		}
		return &connWrap{
			name:      conn.name,
			conn:      nc,
			beginTxFn: conn.beginTxFn,
			newConnFn: conn.newConnFn,
		}, nil
	}
	return conn, nil
}

func (c *connWrap) ServerVersion() string {
	return c.serverVersion
}

func (w *connWrap) Exec(ctx context.Context, query string, args []interface{}) (ExecutedResult, error) {
	var conn sqlPrepare = w.db
	if conn == nil {
		conn = w.conn
	}
	return sqlExec(ctx, w.name, conn, query, args)
}
func (w *connWrap) Query(ctx context.Context, query string, args []interface{}) (Rows, error) {
	var conn sqlPrepare = w.db
	if conn == nil {
		conn = w.conn
	}
	return sqlQuery(ctx, w.name, conn, query, args)
}
func (c *connWrap) Close() error {
	defer func() {
		c.db = nil
		c.conn = nil
	}()

	// TODO err
	if c.db != nil {
		c.db.Close()
	}
	if c.conn != nil {
		c.conn.Close()
	}
	return nil
}

func (c *connWrap) Begin(ctx context.Context) (TransactionInterface, error) {
	return c.beginTxFn(ctx)
}
func (c *connWrap) New(ctx context.Context) (ConnectionInterface, error) {
	return c.newConnFn(ctx)
}

type txWarp struct {
	name string
	tx   *sql.Tx
	err  error
	mu   sync.Mutex
}

func (w *txWarp) internal() {}

func (w *txWarp) Exec(ctx context.Context, query string, args []interface{}) (ExecutedResult, error) {
	return sqlExec(ctx, w.name, w.tx, query, args)
}
func (w *txWarp) Query(ctx context.Context, query string, args []interface{}) (Rows, error) {
	return sqlQuery(ctx, w.name, w.tx, query, args)
}

func (c *txWarp) doClose() error {
	defer func() {
		c.tx = nil
	}()
	if c.tx == nil {
		return nil
	}
	return c.tx.Rollback()
}

func (c *txWarp) Close() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if err := c.doClose(); err != nil {
		fmt.Println("sqlx: 关闭发生错误:", err)
	}
	return c.err
}

func (c *txWarp) check() error {
	if c.err != nil {
		return c.err
	}
	if c.tx == nil {
		return fmt.Errorf("tx was closed")
	}
	return c.err
}

func (c *txWarp) Commit() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	defer c.doClose()
	if err := c.check(); err != nil {
		return err
	}
	c.err = wrapErr(c.tx.Commit())
	if c.err == nil {
		// 防止再次 Rollback
		c.tx = nil
	}
	return c.err
}

func (c *txWarp) Do(run func(conn Transaction) error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if err := c.check(); err != nil {
		return
	}

	defer runtime.HandleCrash(func(r interface{}) {
		c.err = fmt.Errorf("panic with tx: %v", r)
	})
	c.err = run(c)
}

var (
	_, dbInstance  = field.String("db.instance")
	_, dbType      = field.String("db.type")
	_, dbStmt      = field.String("db.statement")
	_, dbArgs      = field.Slice("db.arguments", field.StringKind)
	StatementIDKey = struct{ name string }{"$$statement_id"}
)

type sqlPrepare interface {
	PrepareContext(ctx context.Context, query string) (*sql.Stmt, error)
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

// https://github.com/opentracing/specification/blob/master/semantic_conventions.md
func prepareContext(ctx context.Context, instance string, conn sqlPrepare, query string, args []interface{}) (context.Context, apm.Span, *sql.Stmt, error) {
	var id string
	if s, ok := ctx.Value(StatementIDKey).(string); ok {
		id = s
	}
	// if id == "" {
	// 	id = "sql"
	// }
	ctx, span := apm.Start(ctx, apm.WithName(id), apm.WithCalldepth(4))

	if conn == nil {
		return ctx, span, nil, fmt.Errorf("conn was closed")
	}
	span.Debug(
		dbType("sql"),
		dbInstance(instance),
		dbStmt(query),
		dbArgs(argSlim(args)...),
	)
	stmt, err := conn.PrepareContext(ctx, query)

	return ctx, span, stmt, err
}

func sqlExec(ctx context.Context, instance string, conn sqlPrepare, query string, args []interface{}) (ExecutedResult, error) {
	ctx, span, stmt, err := prepareContext(ctx, instance, conn, query, args)
	defer span.End(apm.WithCalldepth(2))
	if span.FailIf(err) {
		return nil, err
	}
	defer stmt.Close()

	result, err := stmt.ExecContext(ctx, args...)
	if span.FailIf(err) {
		return nil, err
	}
	return result, err
}
func sqlQuery(ctx context.Context, instance string, conn sqlPrepare, query string, args []interface{}) (Rows, error) {
	ctx, span, stmt, err := prepareContext(ctx, instance, conn, query, args)
	defer span.End(apm.WithCalldepth(2))
	if span.FailIf(err) {
		return nil, err
	}
	defer stmt.Close()

	result, err := stmt.QueryContext(ctx, args...)
	if span.FailIf(err) {
		return nil, err
	}
	return result, err
}
