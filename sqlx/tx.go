package sqlx

import (
	"context"
	"database/sql"
	"fmt"
	"sync"

	"github.com/junhwong/goost/runtime"
)

type txWarp struct {
	parent *connWrap
	tx     *sql.Tx
	err    error
	mu     sync.Mutex
}
type beginTx interface {
	BeginTx(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error)
}

func buildBegin(conn *connWrap, btx beginTx) func(ctx context.Context, opts []*sql.TxOptions) (TransactionInterface, error) {
	return func(ctx context.Context, opts []TxOptions) (TransactionInterface, error) {
		conn.mu.Lock()
		defer conn.mu.Unlock()

		if conn.raw == nil {
			return nil, fmt.Errorf("conn was closed")
		}
		var opt TxOptions
		tx, err := btx.BeginTx(ctx, opt)
		if err != nil {
			return nil, err
		}
		return &txWarp{
			// name: conn.name,
			parent: conn,
			tx:     tx,
		}, nil
	}
}

func (conn *txWarp) getInstance() string {
	if conn.parent != nil {
		return conn.parent.getInstance()
	}
	return ""
}

func (w *txWarp) internal() {}

func (c *txWarp) doClose() (err error) {
	if c.tx != nil {
		err = c.tx.Rollback()
	}
	c.tx = nil
	return
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

	if err := c.tx.Commit(); c.err != nil {
		c.err = wrapErr(err)
		return c.err
	}
	c.tx = nil // 防止再次 Rollback
	return c.err
}

func (c *txWarp) Do(fn func(Transaction) error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if err := c.check(); err != nil {
		return
	}

	defer runtime.HandleCrash(func(err error) {
		c.err = err
	})
	c.err = fn(c)
}

// func (w *txWarp) Exec(ctx context.Context, query string, args []interface{}) (ExecutedResult, error) {
// 	return sqlExec(ctx, w, w.tx, query, args)
// }

// func (w *txWarp) Query(ctx context.Context, query string, args []interface{}) (Rows, error) {
// 	return sqlQuery(ctx, w, w.tx, query, args)
// }

// func (conn *txWarp) QueryWithCallback(ctx context.Context, query string, args []interface{}, cb func(Rows) error) error {
// 	return sqlQueryWithCallback(ctx, conn, conn.tx, query, args, cb)
// }
func (conn *txWarp) Prepare(ctx context.Context, query string) (Stmt, error) {
	return prepareContext(ctx, conn, conn.tx, query)
}
