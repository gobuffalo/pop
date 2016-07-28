package fizz_test

import (
	"testing"

	"github.com/markbates/pop/fizz"
	"github.com/stretchr/testify/require"
)

func Test_AddColumn(t *testing.T) {
	r := require.New(t)
	ddl := `add_column("users", "age", "integer", {"null": false, "default":0})`

	ch := fizz.AString(ddl)
	b := <-ch
	tl := b.Data.(*fizz.Table)
	r.Equal("users", tl.Name)

	c := tl.Columns[0]
	r.Equal("age", c.Name)
}
