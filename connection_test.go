package pop

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_Connection_SimpleFlow(t *testing.T) {
	r := require.New(t)

	cd := &ConnectionDetails{
		URL: "sqlite:///foo.db",
	}
	c, err := NewConnection(cd)
	r.NoError(err)

	err = c.Open()
	r.NoError(err)
	err = c.Open() // open again
	r.NoError(err)
	err = c.Close()
	r.NoError(err)
}

func Test_Connection_Open_NoDialect(t *testing.T) {
	r := require.New(t)

	cd := &ConnectionDetails{
		URL: "sqlite:///foo.db",
	}
	c, err := NewConnection(cd)
	r.NoError(err)

	c.Dialect = nil
	err = c.Open()
	r.Error(err)
}

func Test_Connection_Open_BadDialect(t *testing.T) {
	r := require.New(t)

	cd := &ConnectionDetails{
		URL: "sqlite:///foo.db",
	}
	c, err := NewConnection(cd)
	r.NoError(err)

	cd.Dialect = "unknown"
	err = c.Open()
	r.Error(err)
}
