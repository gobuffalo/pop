// +build !sqlite

package pop

import (
	"errors"
	"fmt"
)

func init() {
	fmt.Println("in shim")
}

func newSQLite(deets *ConnectionDetails) (dialect, error) {
	return nil, errors.New("sqlite3 was not compiled into the binary")
}
