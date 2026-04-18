package pop

import (
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"
)

// Store is an interface that must be implemented in order for Pop
// to be able to use the value as a way of talking to a datastore.
type store interface {
	Select(any, string, ...any) error
	Get(any, string, ...any) error
	NamedExec(string, any) (sql.Result, error)
	NamedQuery(query string, arg any) (*sqlx.Rows, error)
	Exec(string, ...any) (sql.Result, error)
	PrepareNamed(string) (*sqlx.NamedStmt, error)
	Transaction() (*Tx, error)
	Rollback() error
	Commit() error
	Close() error

	// Context versions to wrap with contextStore
	SelectContext(context.Context, any, string, ...any) error
	GetContext(context.Context, any, string, ...any) error
	NamedExecContext(context.Context, string, any) (sql.Result, error)
	NamedQueryContext(ctx context.Context, query string, arg any) (*sqlx.Rows, error)
	ExecContext(context.Context, string, ...any) (sql.Result, error)
	PrepareNamedContext(context.Context, string) (*sqlx.NamedStmt, error)
	TransactionContext(context.Context) (*Tx, error)
	TransactionContextOptions(context.Context, *sql.TxOptions) (*Tx, error)
}

// ContextStore wraps a store with a Context, so passes it with the functions that don't take it.
type contextStore struct {
	store
	ctx context.Context
}

func (s contextStore) Transaction() (*Tx, error) {
	return s.TransactionContext(s.ctx)
}

func (s contextStore) Select(dest any, query string, args ...any) error {
	return s.SelectContext(s.ctx, dest, query, args...)
}

func (s contextStore) Get(dest any, query string, args ...any) error {
	return s.GetContext(s.ctx, dest, query, args...)
}

func (s contextStore) NamedExec(query string, arg any) (sql.Result, error) {
	return s.NamedExecContext(s.ctx, query, arg)
}

func (s contextStore) Exec(query string, args ...any) (sql.Result, error) {
	return s.ExecContext(s.ctx, query, args...)
}

func (s contextStore) PrepareNamed(query string) (*sqlx.NamedStmt, error) {
	return s.PrepareNamedContext(s.ctx, query)
}

func (s contextStore) Context() context.Context {
	return s.ctx
}
