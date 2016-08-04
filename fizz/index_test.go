package fizz_test

import (
	"testing"

	"github.com/markbates/pop/fizz"
	"github.com/stretchr/testify/require"
)

func Test_AddIndex(t *testing.T) {
	r := require.New(t)

	ddl := `add_index("users", "email", {})`

	bub, _ := fizz.AString(ddl)
	b := bub.Bubbles[0]
	r.Equal(fizz.E_ADD_INDEX, b.BubbleType)

	tl := b.Data.(fizz.Table)
	r.Equal("users", tl.Name)

	i := tl.Indexes[0]
	r.Equal("users_email_idx", i.Name)
	r.False(i.Unique)
	r.Equal([]string{"email"}, i.Columns)
}

func Test_AddIndex_CustomName(t *testing.T) {
	r := require.New(t)

	ddl := `add_index("users", "email", {"name": "email_index"})`

	bub, _ := fizz.AString(ddl)
	b := bub.Bubbles[0]

	tl := b.Data.(fizz.Table)
	i := tl.Indexes[0]

	r.Equal("email_index", i.Name)
}

func Test_AddIndex_MultipleColumns(t *testing.T) {
	r := require.New(t)

	ddl := `add_index("users", ["email", "username"], {})`

	bub, _ := fizz.AString(ddl)
	b := bub.Bubbles[0]

	tl := b.Data.(fizz.Table)

	i := tl.Indexes[0]

	r.Equal("users_email_username_idx", i.Name)
}

func Test_AddIndex_Unique(t *testing.T) {
	r := require.New(t)

	ddl := `add_index("users", "email", {"unique": true})`

	bub, _ := fizz.AString(ddl)
	b := bub.Bubbles[0]
	tl := b.Data.(fizz.Table)

	i := tl.Indexes[0]
	r.True(i.Unique)
}

func Test_DropIndex(t *testing.T) {
	r := require.New(t)

	ddl := `drop_index("users_email_idx")`

	bub, _ := fizz.AString(ddl)
	b := bub.Bubbles[0]
	r.Equal(fizz.E_DROP_INDEX, b.BubbleType)

	tl := b.Data.(fizz.Index)

	r.Equal("users_email_idx", tl.Name)
}

func Test_RenameIndex(t *testing.T) {
	r := require.New(t)

	ddl := `rename_index("users_email_idx", "email_ix")`

	bub, _ := fizz.AString(ddl)
	b := bub.Bubbles[0]
	r.Equal(fizz.E_RENAME_INDEX, b.BubbleType)

	idx := b.Data.([]fizz.Index)

	r.Len(idx, 2)

	i := idx[0]
	r.Equal("users_email_idx", i.Name)

	i = idx[1]
	r.Equal("email_ix", i.Name)
}
