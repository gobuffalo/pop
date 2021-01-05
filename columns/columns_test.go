package columns_test

import (
	"fmt"
	"testing"

	"github.com/gobuffalo/pop/v5/columns"
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

	c1 := columns.ForStruct(&foo{}, "foo", "id")
	c2 := columns.ForStruct(&foos{}, "foo", "id")
	r.Equal(c1.String(), c2.String())
}

func Test_Columns_Basics(t *testing.T) {
	r := require.New(t)

	for _, f := range []interface{}{foo{}, &foo{}} {
		c := columns.ForStruct(f, "foo", "id")
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
		c := columns.ForStruct(f, "foo", "id")
		r.Equal(len(c.Cols), 4)
		c.Add("foo", "first_name")
		r.Equal(len(c.Cols), 5)
		r.Equal(c.Cols["foo"], &columns.Column{Name: "foo", Writeable: true, Readable: true, SelectSQL: "foo.foo"})
	}
}

func Test_Columns_Remove(t *testing.T) {
	r := require.New(t)

	for _, f := range []interface{}{foo{}, &foo{}} {
		c := columns.ForStruct(f, "foo", "id")
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

	c := columns.ForStruct(fooWithSuffix{}, "fooWithSuffix", "id")
	r.Equal(len(c.Cols), 2)
	r.Equal(c.SymbolizedString(), ":amount, :amount_units")
	r.Equal(c.String(), "amount, amount_units")
	r.Equal(c.QuotedString(fooQuoter{}), "`amount`, `amount_units`")
}

func Test_Columns_IDField(t *testing.T) {
	type withID struct {
		ID string `db:"id"`
	}

	r := require.New(t)
	c := columns.ForStruct(withID{}, "with_id", "id")
	r.Equal(1, len(c.Cols), "%+v", c)
	r.Equal(&columns.Column{Name: "id", Writeable: false, Readable: true, SelectSQL: "with_id.id"}, c.Cols["id"])
}

func Test_Columns_IDField_Readonly(t *testing.T) {
	type withIDReadonly struct {
		ID string `db:"id" rw:"r"`
	}

	r := require.New(t)
	c := columns.ForStruct(withIDReadonly{}, "with_id_readonly", "id")
	r.Equal(1, len(c.Cols), "%+v", c)
	r.Equal(&columns.Column{Name: "id", Writeable: false, Readable: true, SelectSQL: "with_id_readonly.id"}, c.Cols["id"])
}

func Test_Columns_ID_Field_Not_ID(t *testing.T) {
	type withNonStandardID struct {
		PK string `db:"notid"`
	}

	r := require.New(t)

	c := columns.ForStruct(withNonStandardID{}, "non_standard_id", "notid")
	r.Equal(1, len(c.Cols), "%+v", c)
	r.Equal(&columns.Column{Name: "notid", Writeable: false, Readable: true, SelectSQL: "non_standard_id.notid"}, c.Cols["notid"])
}
