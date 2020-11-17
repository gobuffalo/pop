package columns_test

import (
	"testing"

	"github.com/gobuffalo/pop/v5/columns"
	"github.com/stretchr/testify/require"
)

func Test_Columns_WriteableString_Symbolized(t *testing.T) {
	r := require.New(t)
	for _, f := range []interface{}{foo{}, &foo{}} {
		c := columns.ForStruct(f, "foo", "id")
		u := c.Writeable().SymbolizedString()
		r.Equal(u, ":LastName, :write")
	}
}

func Test_Columns_UpdateString(t *testing.T) {
	r := require.New(t)
	for _, f := range []interface{}{foo{}, &foo{}} {
		c := columns.ForStruct(f, "foo", "id")
		u := c.Writeable().UpdateString()
		r.Equal(u, "LastName = :LastName, write = :write")
	}
}

type testQuoter struct{}

func (testQuoter) Quote(col string) string {
	return `"` + col + `"`
}

func Test_Columns_QuotedUpdateString(t *testing.T) {
	r := require.New(t)
	q := testQuoter{}
	for _, f := range []interface{}{foo{}, &foo{}} {
		c := columns.ForStruct(f, "foo", "id")
		u := c.Writeable().QuotedUpdateString(q)
		r.Equal(u, "\"LastName\" = :LastName, \"write\" = :write")
	}
}

func Test_Columns_WriteableString(t *testing.T) {
	r := require.New(t)
	for _, f := range []interface{}{foo{}, &foo{}} {
		c := columns.ForStruct(f, "foo", "id")
		u := c.Writeable().String()
		r.Equal(u, "LastName, write")
	}
}
