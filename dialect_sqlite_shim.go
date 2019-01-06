// +build !sqlite

package pop

import (
	"errors"
)

const nameSQLite3 = "sqlite3"

func init() {
	AvailableDialects = append(AvailableDialects, nameSQLite3)
	dialectSynonyms["sqlite"] = nameSQLite3
	NewConnectionHook[nameSQLite3] = newSQLite
}

func newSQLite(deets *ConnectionDetails) (Dialect, error) {
	return nil, errors.New("sqlite3 support was not compiled into the binary")
}
