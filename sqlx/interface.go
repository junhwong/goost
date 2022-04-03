package sqlx

import (
	"context"
	"database/sql"
	"io"
)

type Rows = *sql.Rows
type ExecutedResult = sql.Result

type ExecutableInterface interface {
	Exec(ctx context.Context, query string, args []interface{}) (ExecutedResult, error)
}

type QueryableInterface interface {
	Query(ctx context.Context, query string, args []interface{}) (Rows, error)
}

// 数据库连接
type Connection interface {
	ExecutableInterface
	QueryableInterface
}

// 事务控制器, 它是线程安全的。
type TransactionInterface interface {
	Close() error                    // 如果没有提交则 rollback(回滚), 否则什么都不做
	Commit() error                   // 提交事务, 只能提交一次
	Do(func(conn Transaction) error) // 在事务中执行操作
}

// 事务连接
type Transaction interface {
	Connection
	internal() // 用于区分 Connection 和 Transaction 的接口标识
}

// 连接接口
type ConnectionInterface interface {
	Connection
	io.Closer
	Begin(ctx context.Context) (TransactionInterface, error)
}

// 创建连接
type Factory interface {
	New(ctx context.Context) (ConnectionInterface, error)
}

type DBInterface interface {
	ServerVersion() string
	ConnectionInterface
	Factory
}
