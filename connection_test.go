//go:build sqlite
// +build sqlite

package pop

import (
	"context"
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

func Test_Connection_Open_BadDriver(t *testing.T) {
	r := require.New(t)

	cd := &ConnectionDetails{
		URL: "sqlite:///foo.db",
	}
	c, err := NewConnection(cd)
	r.NoError(err)

	cd.Driver = "unknown"
	err = c.Open()
	r.Error(err)
}

func Test_Connection_Transaction(t *testing.T) {
	r := require.New(t)
	ctx := context.WithValue(context.Background(), "test", "test")

	c, err := NewConnection(&ConnectionDetails{
		URL: "sqlite://file::memory:?_fk=true",
	})
	r.NoError(err)
	r.NoError(c.Open())
	c = c.WithContext(ctx)

	t.Run("func=NewTransaction", func(t *testing.T) {
		r := require.New(t)
		tx, err := c.NewTransaction()
		r.NoError(err)

		// has transaction and context
		r.NotNil(tx.TX)
		r.Nil(c.TX)
		r.Equal(ctx, tx.Context())

		// does not start a new transaction
		ntx, err := tx.NewTransaction()
		r.Equal(tx, ntx)

		r.NoError(tx.TX.Rollback())
	})

	t.Run("func=NewTransactionContext", func(t *testing.T) {
		r := require.New(t)
		nctx := context.WithValue(ctx, "nested", "test")
		tx, err := c.NewTransactionContext(nctx)
		r.NoError(err)

		// has transaction and context
		r.NotNil(tx.TX)
		r.Nil(c.TX)
		r.Equal(nctx, tx.Context())

		r.NoError(tx.TX.Rollback())
	})
}
