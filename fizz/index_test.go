package fizz_test

import (
	"testing"

	"github.com/markbates/pop/fizz"
	"github.com/stretchr/testify/require"
)

func Test_AddIndex(t *testing.T) {
	r := require.New(t)

	ddl := `add_index("users", "email", {})`

	b := <-fizz.AString(ddl).Bubbles
	r.Equal(fizz.E_ADD_INDEX, b.BubbleType)

	i := b.Data.(*fizz.Index)

	r.Equal("users_email_idx", i.Name)
	r.False(i.Unique)
	r.Equal([]string{"email"}, i.Columns)
	r.Equal("users", i.TableName)
}

func Test_AddIndex_CustomName(t *testing.T) {
	r := require.New(t)

	ddl := `add_index("users", "email", {"name": "email_index"})`

	b := <-fizz.AString(ddl).Bubbles
	i := b.Data.(*fizz.Index)

	r.Equal("email_index", i.Name)
}

func Test_AddIndex_MultipleColumns(t *testing.T) {
	r := require.New(t)

	ddl := `add_index("users", ["email", "username"], {})`

	b := <-fizz.AString(ddl).Bubbles
	i := b.Data.(*fizz.Index)

	r.Equal("users_email_username_idx", i.Name)
}

func Test_AddIndex_Unique(t *testing.T) {
	r := require.New(t)

	ddl := `add_index("users", "email", {"unique": true})`

	b := <-fizz.AString(ddl).Bubbles
	i := b.Data.(*fizz.Index)

	r.True(i.Unique)
}

func Test_DropIndex(t *testing.T) {
	r := require.New(t)

	ddl := `drop_index("users", "users_email_idx")`

	b := <-fizz.AString(ddl).Bubbles
	r.Equal(fizz.E_DROP_INDEX, b.BubbleType)

	i := b.Data.(*fizz.Index)
	r.Equal("users", i.TableName)
	r.Equal("users_email_idx", i.Name)
}

func Test_RenameIndex(t *testing.T) {
	r := require.New(t)

	ddl := `rename_index("users", "users_email_idx", "email_ix")`

	b := <-fizz.AString(ddl).Bubbles
	r.Equal(fizz.E_RENAME_INDEX, b.BubbleType)

	idx := b.Data.([]*fizz.Index)
	r.Len(idx, 2)

	i := idx[0]
	r.Equal("users", i.TableName)
	r.Equal("users_email_idx", i.Name)

	i = idx[1]
	r.Equal("users", i.TableName)
	r.Equal("email_ix", i.Name)
}
