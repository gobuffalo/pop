package associations

import (
	"reflect"

	"github.com/markbates/pop/columns"
)

// Association represents a definition of a model association
// field. It can represent a association of the type has_many
// belongs_to or has_one, and other customized types.
type Association interface {
	Kind() reflect.Kind
	Interface() interface{}
	Constraint() (string, []interface{})
}

// AssociationSortable a type to be sortable.
type AssociationSortable interface {
	OrderBy() string
	Association
}

// Associations a group of model associations.
type Associations []Association

// associationParams a wrapper for associations definition
// and creation.
type associationParams struct {
	field      reflect.StructField // an association field defined in model.
	modelType  reflect.Type        // the model type where this field is defined.
	modelValue reflect.Value       // the model value where this field is defined.
	popTags    columns.Tags        // the tags defined in this association field.
	model      interface{}         // the model, owner of the association.
}

// associationBuilder is a type representing an association builder implementation.
// see the builder defined in ./has_many_association.go as a guide of how to use it.
type associationBuilder func(associationParams) (Association, error)
