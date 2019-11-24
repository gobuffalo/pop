package pop

import (
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"
)

// Store is an interface that must be implemented in order for Pop
// to be able to use the value as a way of talking to a datastore.
type store interface {
	Select(interface{}, string, ...interface{}) error
	SelectContext(context.Context, interface{}, string, ...interface{}) error

	Get(interface{}, string, ...interface{}) error
	GetContext(context.Context,interface{}, string, ...interface{}) error

	NamedExec(string, interface{}) (sql.Result, error)
	NamedExecContext(context.Context,string, interface{}) (sql.Result, error)

	Exec(string, ...interface{}) (sql.Result, error)
	ExecContext(context.Context,string, ...interface{}) (sql.Result, error)

	PrepareNamed(string) (*sqlx.NamedStmt, error)
	PrepareNamedContext(context.Context,string) (*sqlx.NamedStmt, error)

	Transaction() (*Tx, error)
	TransactionContext(context.Context) (*Tx, error)

	Rollback() error
	Commit() error
	Close() error
}
