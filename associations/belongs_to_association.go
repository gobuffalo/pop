package associations

import (
	"reflect"

	"github.com/markbates/inflect"
)

// belongsToAssociation is the implementation for the belongs_to
// association type in a model.
type belongsToAssociation struct {
	ownerModel reflect.Value
	ownerType  reflect.Type
	ownerID    reflect.Value
}

func (b *belongsToAssociation) TableName() string {
	return inflect.Tableize(b.ownerType.Name())
}

func (b *belongsToAssociation) FieldName() string {
	return b.ownerType.Name()
}

func (b *belongsToAssociation) Type() reflect.Kind {
	if b.ownerType.Kind() == reflect.Ptr {
		return b.ownerType.Elem().Kind()
	}
	return b.ownerType.Kind()
}

func (b *belongsToAssociation) Interface() interface{} {
	if b.ownerModel.Kind() == reflect.Ptr {
		val := reflect.New(b.ownerType.Elem())
		b.ownerModel.Set(val)
		return b.ownerModel.Interface()
	}
	return b.ownerModel.Addr().Interface()
}

// SQLConstraint returns the content for a where clause, and the args
// needed to execute it.
func (b *belongsToAssociation) SQLConstraint() (string, []interface{}) {
	return "id = ?", []interface{}{b.ownerID.Interface()}
}
