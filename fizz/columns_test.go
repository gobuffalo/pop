package fizz_test

import (
	"testing"

	"github.com/markbates/pop/fizz"
	"github.com/stretchr/testify/require"
)

func Test_AddColumn(t *testing.T) {
	r := require.New(t)
	ddl := `add_column("users", "age", "integer", {"null": false, "default":0})`

	bub, _ := fizz.AString(ddl)
	b := bub.Bubbles[0]
	r.Equal(fizz.E_ADD_COLUMN, b.BubbleType)

	tl := b.Data.(fizz.Table)
	r.Equal("users", tl.Name)

	c := tl.Columns[0]
	r.Equal("age", c.Name)
}

func Test_DropColumn(t *testing.T) {
	r := require.New(t)
	ddl := `drop_column("users", "email")`
	bub, _ := fizz.AString(ddl)
	b := bub.Bubbles[0]
	r.Equal(fizz.E_DROP_COLUMN, b.BubbleType)

	tl := b.Data.(fizz.Table)
	r.Equal("users", tl.Name)

	c := tl.Columns[0]
	r.Equal("email", c.Name)
}

func Test_RenameColumn(t *testing.T) {
	r := require.New(t)
	ddl := `rename_column("users", "email", "email_address")`
	bub, _ := fizz.AString(ddl)
	b := bub.Bubbles[0]
	r.Equal(fizz.E_RENAME_COLUMN, b.BubbleType)

	tl := b.Data.(fizz.Table)
	r.Equal("users", tl.Name)

	c := tl.Columns[0]
	r.Equal("email", c.Name)
	c = tl.Columns[1]
	r.Equal("email_address", c.Name)
}
