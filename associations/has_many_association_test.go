package associations_test

import (
	"reflect"
	"testing"

	"github.com/gobuffalo/pop/associations"
	"github.com/gobuffalo/uuid"
	"github.com/stretchr/testify/require"
)

type fooHasMany struct {
	ID           uuid.UUID    `db:"id"`
	BarHasManies barHasManies `has_many:"bar_has_manies"`
}

type barHasMany struct {
	FooHasManyID uuid.UUID `db:"foo_has_many_id"`
}

type barHasManies []barHasMany

func Test_Has_Many_Association(t *testing.T) {
	a := require.New(t)

	id, _ := uuid.NewV1()
	foo := fooHasMany{ID: id}

	as, err := associations.AssociationsForStruct(&foo)

	a.NoError(err)
	a.Equal(len(as), 1)
	a.Equal(reflect.Slice, as[0].Kind())

	where, args := as[0].Constraint()
	a.Equal("foo_has_many_id = ?", where)
	a.Equal(id, args[0].(uuid.UUID))
}
