// +build !sqlite

package pop

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSqlite_NewDriver(t *testing.T) {
	_, err := newSQLiteDriver()
	require.Error(t, err)
}
