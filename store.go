package pop

import (
	"database/sql"

	"github.com/jmoiron/sqlx"
)

type Store interface {
	Select(dest interface{}, query string, args ...interface{}) error
	Get(dest interface{}, query string, args ...interface{}) error
	NamedExec(query string, arg interface{}) (sql.Result, error)
	Exec(query string, args ...interface{}) (sql.Result, error)
	PrepareNamed(query string) (*sqlx.NamedStmt, error)
	Transaction() (*TX, error)
	Rollback() error
	Commit() error
}
