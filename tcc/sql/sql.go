package sql

import (
	gosql "database/sql"

	"github.com/junhwong/goost/tcc"
)

func Try(ctx tcc.TransactionContext, phase string, fn tcc.TryFn) error {
	p := ProviderWithContext(ctx)
	p.BeginLocalTx()
	state, err := fn(ctx)
	if err != nil {
		return err
	}
	if state.Name == "" {

	}
	return p.SaveState(state)
}

type Provider struct {
	*gosql.DB
}

func (p *Provider) Name() string {
	return "sql"
}
func (p *Provider) BeginLocalTx() (*gosql.Tx, error) {
	return p.Begin()
}
func (p *Provider) SaveState(tcc.State) error {
	return nil
}
func NewSqlProvider(options interface{}) *Provider {
	return &Provider{}
}

func ProviderWithContext(ctx tcc.TransactionContext) *Provider {
	return nil
}

func Do(db *gosql.DB) {

}

type DBExector interface {
	Query(query string, args ...interface{}) (*gosql.Rows, error)
	Exec(query string, args ...interface{}) (gosql.Result, error)
}
type DBInstance interface {
	DBExector
	Close()
}
type DBTransactionExector interface {
	Do(DBExector)
}
type DBTx interface {
	DBTransactionExector
	Done()
	Commit() error
}
