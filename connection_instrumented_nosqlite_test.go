// +build !sqlite

package pop

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestInstrumentation_WithoutSqlite(t *testing.T) {
	_, _, err := instrumentDriver(&ConnectionDetails{
		URL: "sqlite://:memory:",
	}, "sqlite")
	require.NoError(t, err)
}
