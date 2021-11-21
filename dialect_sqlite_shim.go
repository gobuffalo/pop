// +build !sqlite

package pop

import (
	"database/sql/driver"
	"errors"
)

const nameSQLite3 = "sqlite3"

func init() {
	dialectSynonyms["sqlite"] = nameSQLite3
	newConnection[nameSQLite3] = newSQLite
}

func newSQLite(deets *ConnectionDetails) (dialect, error) {
	return nil, errors.New("sqlite3 support was not compiled into the binary")
}

func newSQLiteDriver() (driver.Driver, error) {
	return nil, errors.New("sqlite3 support was not compiled into the binary")
}
