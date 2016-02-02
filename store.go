package pop

import (
	"database/sql"

	"github.com/jmoiron/sqlx"
)

// Store is an interface that must be implemented in order for Pop
// to be able to use the value as a way of talking to a datastore.
type Store interface {
	Select(dest interface{}, query string, args ...interface{}) error
	Get(dest interface{}, query string, args ...interface{}) error
	NamedExec(query string, arg interface{}) (sql.Result, error)
	Exec(query string, args ...interface{}) (sql.Result, error)
	PrepareNamed(query string) (*sqlx.NamedStmt, error)
	Transaction() (*tX, error)
	Rollback() error
	Commit() error
}
