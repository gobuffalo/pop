package associations_test

import (
	"reflect"
	"testing"

	"github.com/gobuffalo/pop/v6/associations"
	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/require"
)

type fooManyToMany struct {
	ID              uuid.UUID       `db:"id"`
	BarManyToManies barManyToManies `many_to_many:"foos_and_bars"`
}

type fooManyToMany2 struct {
	ID              uuid.UUID       `db:"id"`
	BarManyToManies barManyToManies `many_to_many:"foos_and_bars" primary_id:"fufu_id"`
}

type barManyToMany struct {
	ID uuid.UUID `db:"id"`
}

type barManyToManies []barManyToMany

func (b barManyToManies) TableName() string {
	return "bars"
}

func Test_Many_To_Many_Association(t *testing.T) {
	a := require.New(t)

	id, _ := uuid.NewV1()
	foo := fooManyToMany{ID: id}

	as, err := associations.ForStruct(&foo)

	a.NoError(err)
	a.Equal(len(as), 1)

	a.Equal(reflect.Slice, as[0].Kind())

	where, args := as[0].Constraint()
	a.Equal("id in (select bar_many_to_many_id from foos_and_bars where foo_many_to_many_id = ?)", where)
	a.Equal(id, args[0].(uuid.UUID))
}

func Test_Many_To_Many_PrimaryId_Tag(t *testing.T) {
	a := require.New(t)

	id, _ := uuid.NewV1()
	foo := fooManyToMany2{ID: id}

	as, err := associations.ForStruct(&foo)

	a.NoError(err)
	a.Equal(len(as), 1)

	a.Equal(reflect.Slice, as[0].Kind())

	where, args := as[0].Constraint()
	a.Equal("id in (select bar_many_to_many_id from foos_and_bars where fufu_id = ?)", where)
	a.Equal(id, args[0].(uuid.UUID))
}
