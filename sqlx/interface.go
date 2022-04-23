package sqlx

import (
	"context"
	"database/sql"
	"io"
)

type Rows = *sql.Rows
type ExecutedResult = sql.Result
type TxOptions = *sql.TxOptions

type Conn interface {
	Prepare(ctx context.Context, query string) (Stmt, error)
}

type Stmt interface {
	io.Closer
	Exec(ctx context.Context, args ...interface{}) (ExecutedResult, error)
	Query(ctx context.Context, args ...interface{}) (Rows, error)
	QueryRow(ctx context.Context, args ...interface{}) (*sql.Row, error)
}

// 事务控制器, 它是线程安全的。
type TransactionInterface interface {
	io.Closer                        // 如果没有提交则 rollback(回滚), 否则什么都不做
	Commit() error                   // 提交事务, 只能提交一次
	Do(func(conn Transaction) error) // 在事务中执行操作
}

// 事务连接
type Transaction interface {
	Conn
	internal() // 用于区分 Conn 和 Transaction 的接口标识
}

// 连接接口
type ConnectionInterface interface {
	io.Closer
	Conn

	Begin(ctx context.Context, opts ...TxOptions) (TransactionInterface, error)
}

type DBInterface interface {
	ConnectionInterface
	New(ctx context.Context) (ConnectionInterface, error)

	ServerVersion() string
}
