package columns_test

import (
	"fmt"
	"testing"

	"github.com/gobuffalo/pop/columns"
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

	c1 := columns.ForStruct(&foo{}, "foo")
	c2 := columns.ForStruct(&foos{}, "foo")
	r.Equal(c1.String(), c2.String())
}

func Test_Columns_Basics(t *testing.T) {
	r := require.New(t)

	for _, f := range []interface{}{foo{}, &foo{}} {
		c := columns.ForStruct(f, "foo")
		r.Equal(len(c.Cols), 4)
		r.Equal(c.Cols["first_name"], &columns.Column{Name: "first_name", Writeable: false, Readable: true, SelectSQL: "first_name as f"})
		r.Equal(c.Cols["LastName"], &columns.Column{Name: "LastName", Writeable: true, Readable: true, SelectSQL: "foo.LastName"})
		r.Equal(c.Cols["read"], &columns.Column{Name: "read", Writeable: false, Readable: true, SelectSQL: "foo.read"})
		r.Equal(c.Cols["write"], &columns.Column{Name: "write", Writeable: true, Readable: false, SelectSQL: "foo.write"})
	}
}

func Test_Columns_Add(t *testing.T) {
	r := require.New(t)

	for _, f := range []interface{}{foo{}, &foo{}} {
		c := columns.ForStruct(f, "foo")
		r.Equal(len(c.Cols), 4)
		c.Add("foo", "first_name")
		r.Equal(len(c.Cols), 5)
		r.Equal(c.Cols["foo"], &columns.Column{Name: "foo", Writeable: true, Readable: true, SelectSQL: "foo.foo"})
	}
}

func Test_Columns_Remove(t *testing.T) {
	r := require.New(t)

	for _, f := range []interface{}{foo{}, &foo{}} {
		c := columns.ForStruct(f, "foo")
		r.Equal(len(c.Cols), 4)
		c.Remove("foo", "first_name")
		r.Equal(len(c.Cols), 3)
	}
}

type fooWithSuffix struct {
	Amount      float64 `db:"amount"`
	AmountUnits string  `db:"amount_units"`
}
type fooQuoter struct{}

func (fooQuoter) Quote(key string) string {
	return fmt.Sprintf("`%v`", key)
}

func Test_Columns_Sorted(t *testing.T) {
	r := require.New(t)

	c := columns.ForStruct(fooWithSuffix{}, "fooWithSuffix")
	r.Equal(len(c.Cols), 2)
	r.Equal(c.SymbolizedString(), ":amount, :amount_units")
	r.Equal(c.String(), "amount, amount_units")
	r.Equal(c.QuotedString(fooQuoter{}), "`amount`, `amount_units`")
}
