package associations_test

import (
	"reflect"
	"testing"

	"github.com/gobuffalo/nulls"
	"github.com/gobuffalo/pop/associations"
	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/require"
)

type FooHasOne struct {
	ID        uuid.UUID `db:"id"`
	BarHasOne barHasOne `has_one:"barHasOne"`
}

type barHasOne struct {
	Title       string     `db:"title"`
	FooHasOneID nulls.UUID `db:"foo_has_one_id"`
	BazHasOneID nulls.UUID `db:"baz_id"`
}

type bazHasOne struct {
	ID        uuid.UUID  `db:"id"`
	BarHasOne *barHasOne `has_one:"barHasOne" fk_id:"baz_id"`
}

func Test_Has_One_Association(t *testing.T) {
	a := require.New(t)

	id, _ := uuid.NewV1()
	foo := FooHasOne{ID: id}

	as, err := associations.ForStruct(&foo)

	a.NoError(err)
	a.Equal(len(as), 1)
	a.Equal(reflect.Struct, as[0].Kind())

	where, args := as[0].Constraint()
	a.Equal("foo_has_one_id = ?", where)
	a.Equal(id, args[0].(uuid.UUID))

	foo2 := FooHasOne{}

	as, err = associations.ForStruct(&foo2)
	a.NoError(err)
	after := as.AssociationsAfterCreatable()
	for index := range after {
		a.Equal(nil, after[index].AfterInterface())
	}

	baz := bazHasOne{ID: id}

	as, err = associations.ForStruct(&baz)

	a.NoError(err)
	a.Equal(len(as), 1)
	a.Equal(reflect.Struct, as[0].Kind())

	where, args = as[0].Constraint()
	a.Equal("baz_id = ?", where)
	a.Equal(id, args[0].(uuid.UUID))
}

func Test_Has_One_SetValue(t *testing.T) {
	a := require.New(t)
	id, _ := uuid.NewV1()
	foo := FooHasOne{ID: id, BarHasOne: barHasOne{Title: "bar"}}

	as, _ := associations.ForStruct(&foo)
	a.Equal(len(as), 1)

	ca, ok := as[0].(associations.AssociationAfterCreatable)
	a.True(ok)

	a.NoError(ca.AfterSetup())
	a.Equal(foo.ID, foo.BarHasOne.FooHasOneID.Interface().(uuid.UUID))
}
