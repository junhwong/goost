package sqlx

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"sync"
)

type connWrap struct {
	serverVersion string
	name          string
	beginTxFn     func(ctx context.Context, opts []TxOptions) (TransactionInterface, error)
	newConnFn     func(ctx context.Context) (ConnectionInterface, error)
	parent        *connWrap
	raw           sqlConn
	mu            sync.Mutex
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
		name: config.Driver + "://" + config.Name,
		raw:  db,
	}

	// if rows, err := db.Query("SELECT VERSION()"); err != nil {
	// 	return nil, err
	// } else {
	// 	defer rows.Close()
	// 	for rows.Next() {
	// 		if err := rows.Scan(&conn.serverVersion); err != nil {
	// 			return nil, err
	// 		}
	// 	}
	// }

	conn.beginTxFn = buildBegin(conn, db)
	conn.newConnFn = func(ctx context.Context) (ConnectionInterface, error) {
		conn.mu.Lock()
		defer conn.mu.Unlock()

		if conn.raw == nil {
			return nil, fmt.Errorf("conn was closed")
		}

		nc, err := db.Conn(ctx)
		if err != nil {
			return nil, err
		}
		nconn := &connWrap{
			parent:    conn,
			raw:       nc,
			newConnFn: conn.newConnFn,
		}
		nconn.beginTxFn = buildBegin(nconn, nc)

		return nconn, nil
	}

	return conn, nil
}
func (conn *connWrap) getInstance() string {
	if conn.parent != nil {
		return conn.parent.getInstance()
	}
	return conn.name
}

func (conn *connWrap) ServerVersion() string {
	return conn.serverVersion
}

func (conn *connWrap) Close() (err error) {
	conn.mu.Lock()
	defer conn.mu.Unlock()

	if conn.raw != nil {
		err = conn.raw.Close()
	}
	conn.raw = nil
	return
}

func (c *connWrap) Begin(ctx context.Context, opts ...TxOptions) (TransactionInterface, error) {
	return c.beginTxFn(ctx, opts)
}
func (c *connWrap) New(ctx context.Context) (ConnectionInterface, error) {
	return c.newConnFn(ctx)
}

func (conn *connWrap) Prepare(ctx context.Context, query string) (Stmt, error) {
	return prepareContext(ctx, conn, conn.raw, query)
}
