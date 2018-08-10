package pop

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_LoadsConnectionsFromConfig(t *testing.T) {
	r := require.New(t)

	conns := Connections
	r.Equal(5, len(conns))
}

func Test_AddLookupPaths(t *testing.T) {
	r := require.New(t)
	AddLookupPaths("./foo")
	r.Contains(LookupPaths(), "./foo")
}
