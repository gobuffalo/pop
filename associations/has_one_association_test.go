package associations_test

import (
	"reflect"
	"testing"

	"github.com/gobuffalo/pop/associations"
	"github.com/gobuffalo/uuid"
	"github.com/stretchr/testify/require"
)

type fooHasOne struct {
	ID        uuid.UUID `db:"id"`
	BarHasOne barHasOne `has_one:"barHasOne"`
}

type barHasOne struct {
	FooHasOneID uuid.UUID `db:"foo_has_one_id"`
}

func Test_Has_One_Association(t *testing.T) {
	a := require.New(t)

	id, _ := uuid.NewV1()
	foo := fooHasOne{ID: id}

	as, err := associations.AssociationsForStruct(&foo)

	a.NoError(err)
	a.Equal(len(as), 1)
	a.Equal(reflect.Struct, as[0].Kind())

	where, args := as[0].Constraint()
	a.Equal("foo_has_one_id = ?", where)
	a.Equal(id, args[0].(uuid.UUID))
}
