package fizz_test

import (
	"testing"

	"github.com/markbates/pop/fizz"
	"github.com/stretchr/testify/require"
)

func Test_CreateTable(t *testing.T) {
	r := require.New(t)
	ddl := `
create_table("users", func(t) {
	t.Column("email", "string", {"null": false})
	t.Column("name", "string", {})
	t.Column("token", "string", {})
})
`

	ch := fizz.AString(ddl)
	b := <-ch
	tl := b.Data.(*fizz.Table)
	r.Equal("users", tl.Name)
	r.Equal(6, len(tl.Columns))

	for i, n := range []string{"id", "created_at", "updated_at", "email", "name", "token"} {
		c := tl.Columns[i]
		r.Equal(n, c.Name)
	}
}

func Test_DropTable(t *testing.T) {
	r := require.New(t)
	ddl := `drop_table("users")`

	ch := fizz.AString(ddl)
	b := <-ch
	tl := b.Data.(*fizz.Table)
	r.Equal("users", tl.Name)
}
