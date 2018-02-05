package associations_test

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/markbates/pop/associations"

	"github.com/satori/go.uuid"
)

type fooBelongsTo struct {
	ID         uuid.UUID   `db:"id"`
	NestedBars []nestedBar `has_many:"nested_bars"`
}

func (f fooBelongsTo) TableName() string {
	return "foosy"
}

type barBelongsTo struct {
	FooBelongsToID uuid.UUID    `db:"foo_id"`
	Foo            fooBelongsTo `belongs_to:"foo"`
	NestedBar      nestedBar    `has_one:"nestedBar"`
}

type nestedBar struct {
	ID             uuid.UUID `db:"id"`
	fooBelongsToID uuid.UUID `db:"foo_belongs_to_id"`
}

func Test_Belongs_To_Association(t *testing.T) {
	a := require.New(t)

	id, _ := uuid.NewV1()
	bar := barBelongsTo{FooBelongsToID: id}

	as, err := associations.AssociationsForStruct(&bar, "Foo.NestedBars")
	a.NoError(err)
	a.Equal(len(as), 1)
	a.Equal(reflect.Struct, as[0].Kind())

	where, args := as[0].Constraint()
	a.Equal("id = ?", where)
	a.Equal(id, args[0].(uuid.UUID))
}
