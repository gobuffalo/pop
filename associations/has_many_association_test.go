package associations_test

import (
	"reflect"
	"testing"

	"github.com/gobuffalo/pop/associations"
	"github.com/gobuffalo/pop/nulls"
	"github.com/stretchr/testify/require"
)

type FooHasMany struct {
	ID           int           `db:"id"`
	BarHasManies *barHasManies `has_many:"bar_has_manies"`
}

type barHasMany struct {
	Title        string    `db:"title"`
	FooHasManyID nulls.Int `db:"foo_has_many_id"`
}

type barHasManies []barHasMany

func Test_Has_Many_Association(t *testing.T) {
	a := require.New(t)

	id := 1
	foo := FooHasMany{ID: 1}

	as, err := associations.ForStruct(&foo)

	a.NoError(err)
	a.Equal(len(as), 1)
	a.Equal(reflect.Slice, as[0].Kind())

	where, args := as[0].Constraint()
	a.Equal("foo_has_many_id = ?", where)
	a.Equal(id, args[0].(int))
}

func Test_Has_Many_SetValue(t *testing.T) {
	a := require.New(t)
	foo := FooHasMany{ID: 1, BarHasManies: &barHasManies{{Title: "bar"}}}

	as, _ := associations.ForStruct(&foo)
	a.Equal(len(as), 1)

	ca, ok := as[0].(associations.AssociationAfterCreatable)
	a.True(ok)

	ca.AfterSetup()
	a.Equal(foo.ID, (*foo.BarHasManies)[0].FooHasManyID.Interface().(int))
}
