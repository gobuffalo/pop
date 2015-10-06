package pop_test

import (
	"testing"

	"github.com/markbates/pop"
	"github.com/stretchr/testify/require"
)

type foo struct {
	FirstName string `db:"first_name" select:"first_name as f"`
	LastName  string
	Unwanted  string `db:"-"`
	ReadOnly  string `db:"read" rw:"r"`
	WriteOnly string `db:"write" rw:"w"`
}

type foos []foo

func Test_Column_MapsSlice(t *testing.T) {
	r := require.New(t)

	c1 := pop.ColumnsForStruct(&foo{})
	c2 := pop.ColumnsForStruct(&foos{})
	r.Equal(c1.String(), c2.String())
}

func Test_Column_UpdateString(t *testing.T) {
	r := require.New(t)
	c := pop.Column{Name: "foo"}
	r.Equal(c.UpdateString(), "foo = :foo")
}

func Test_Columns_UpdateString(t *testing.T) {
	r := require.New(t)
	for _, f := range []interface{}{foo{}, &foo{}} {
		c := pop.ColumnsForStruct(f)
		u := c.Writeable().UpdateString()
		r.Equal(u, "LastName = :LastName, write = :write")
	}
}

func Test_Columns_WriteableString(t *testing.T) {
	r := require.New(t)
	for _, f := range []interface{}{foo{}, &foo{}} {
		c := pop.ColumnsForStruct(f)
		u := c.Writeable().String()
		r.Equal(u, "LastName, write")
	}
}

func Test_Columns_ReadableString(t *testing.T) {
	r := require.New(t)
	for _, f := range []interface{}{foo{}, &foo{}} {
		c := pop.ColumnsForStruct(f)
		u := c.Readable().String()
		r.Equal(u, "LastName, first_name, read")
	}
}

func Test_Columns_Readable_SelectString(t *testing.T) {
	r := require.New(t)
	for _, f := range []interface{}{foo{}, &foo{}} {
		c := pop.ColumnsForStruct(f)
		u := c.Readable().SelectString()
		r.Equal(u, "LastName, first_name as f, read")
	}
}

func Test_Columns_WriteableString_Symbolized(t *testing.T) {
	r := require.New(t)
	for _, f := range []interface{}{foo{}, &foo{}} {
		c := pop.ColumnsForStruct(f)
		u := c.Writeable().SymbolizedString()
		r.Equal(u, ":LastName, :write")
	}
}

func Test_Columns_ReadableString_Symbolized(t *testing.T) {
	r := require.New(t)
	for _, f := range []interface{}{foo{}, &foo{}} {
		c := pop.ColumnsForStruct(f)
		u := c.Readable().SymbolizedString()
		r.Equal(u, ":LastName, :first_name, :read")
	}
}
func Test_Columns_Basics(t *testing.T) {
	r := require.New(t)

	for _, f := range []interface{}{foo{}, &foo{}} {
		c := pop.ColumnsForStruct(f)
		r.Equal(len(c.Cols), 4)
		r.Equal(c.Cols["first_name"], &pop.Column{Name: "first_name", Writeable: false, Readable: true, SelectSQL: "first_name as f"})
		r.Equal(c.Cols["LastName"], &pop.Column{Name: "LastName", Writeable: true, Readable: true, SelectSQL: "LastName"})
		r.Equal(c.Cols["read"], &pop.Column{Name: "read", Writeable: false, Readable: true, SelectSQL: "read"})
		r.Equal(c.Cols["write"], &pop.Column{Name: "write", Writeable: true, Readable: false, SelectSQL: "write"})
	}
}

func Test_Columns_Add(t *testing.T) {
	r := require.New(t)

	for _, f := range []interface{}{foo{}, &foo{}} {
		c := pop.ColumnsForStruct(f)
		r.Equal(len(c.Cols), 4)
		c.Add("foo", "first_name")
		r.Equal(len(c.Cols), 5)
		r.Equal(c.Cols["foo"], &pop.Column{Name: "foo", Writeable: true, Readable: true, SelectSQL: "foo"})
	}
}

func Test_Columns_Remove(t *testing.T) {
	r := require.New(t)

	for _, f := range []interface{}{foo{}, &foo{}} {
		c := pop.ColumnsForStruct(f)
		r.Equal(len(c.Cols), 4)
		c.Remove("foo", "first_name")
		r.Equal(len(c.Cols), 3)
	}
}
