package pop_test

import (
	"testing"

	"github.com/gobuffalo/pop"
	"github.com/stretchr/testify/require"
)

func Test_LoadsConnectionsFromConfig(t *testing.T) {
	r := require.New(t)

	conns := pop.Connections
	r.Equal(5, len(conns))
}

func Test_AddLookupPaths(t *testing.T) {
	r := require.New(t)
	pop.AddLookupPaths("./foo")
	r.Contains(pop.LookupPaths(), "./foo")
}
