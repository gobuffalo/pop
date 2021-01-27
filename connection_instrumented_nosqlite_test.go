// +build !sqlite

package pop

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestInstrumentation_WithoutSqlite(t *testing.T) {
	_, _, err := instrumentDriver(&ConnectionDetails{
		URL: "sqlite://:memory:",
	}, "sqlite")
	require.NoError(t, err)
}
