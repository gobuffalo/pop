package associations

import (
	"reflect"

	"github.com/markbates/pop/columns"
	"github.com/markbates/pop/nulls"
)

// Association represents a definition of a model association
// field. It can represent a association of the type has_many
// belongs_to or has_one, and other customized types.
type Association interface {
	Kind() reflect.Kind
	Interface() interface{}
	Constraint() (string, []interface{})
	Dependencies() []interface{}
	SetValue([]interface{}) error
	Skipped() bool
}

// associationSkipable is a helper struct that helps
// to include skippable behavior in associations.
type associationSkipable struct {
	skipped bool
}

func (a *associationSkipable) Skipped() bool {
	return a.skipped
}

// AssociationSortable a type to be sortable.
type AssociationSortable interface {
	OrderBy() string
	Association
}

// AssociationCreatableStatement a association that stablish
// some statements on create.
type AssociationCreatableStatement interface {
	Statements() []AssociationStatement
}

// AssociationStatement a type that represents a statement to be
// executed.
type AssociationStatement struct {
	Statement string
	Args      []interface{}
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

// fieldIsNil validates if a field has a nil reference. Also
// it validates if a field implements nullable interface and
// it has a nil value.
func fieldIsNil(f reflect.Value) bool {
	if n := nulls.New(f.Interface()); n != nil {
		return n.Interface() == nil
	}
	return f.Interface() == nil
}

func isZero(i interface{}) bool {
	v := reflect.ValueOf(i)
	return v.Interface() == reflect.Zero(v.Type()).Interface()
}
