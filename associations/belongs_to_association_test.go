package associations_test

import (
	"reflect"
	"testing"

	"github.com/gobuffalo/pop/v6/associations"
	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/require"
)

type fooBelongsTo struct {
	ID uuid.UUID `db:"id"`
}

func (f fooBelongsTo) TableName() string {
	return "foosy"
}

type barBelongsTo struct {
	FooID uuid.UUID    `db:"foo_id"`
	Foo   fooBelongsTo `belongs_to:"foo"`
}

type barBelongsToNullable struct {
	FooID uuid.NullUUID `db:"foo_id"`
	Foo   *fooBelongsTo `belongs_to:"foo"`
}

func Test_Belongs_To_Association(t *testing.T) {
	a := require.New(t)

	id, _ := uuid.NewV1()
	bar := barBelongsTo{FooID: id}

	as, err := associations.ForStruct(&bar, "Foo")
	a.NoError(err)
	a.Equal(len(as), 1)
	a.Equal(reflect.Struct, as[0].Kind())

	where, args := as[0].Constraint()
	a.Equal("id = ?", where)
	a.Equal(id, args[0].(uuid.UUID))

	bar2 := barBelongsTo{FooID: uuid.Nil}
	as, err = associations.ForStruct(&bar2, "Foo")

	a.NoError(err)
	a.Equal(len(as), 1)
	a.Equal(reflect.Struct, as[0].Kind())

	before := as.AssociationsBeforeCreatable()

	for index := range before {
		a.Equal(nil, before[index].BeforeInterface())
	}
}

func Test_Belongs_To_Nullable_Association(t *testing.T) {
	a := require.New(t)
	id, _ := uuid.NewV1()

	bar := barBelongsToNullable{Foo: &fooBelongsTo{id}}
	as, err := associations.ForStruct(&bar, "Foo")
	a.NoError(err)

	before := as.AssociationsBeforeCreatable()
	for index := range before {
		a.Equal(nil, before[index].BeforeSetup())
	}
}
