package associations

import (
	"reflect"

	"github.com/gobuffalo/pop/columns"
	"github.com/gobuffalo/pop/nulls"
)

// Association represents a definition of a model association
// field. It can represent a association of the type has_many
// belongs_to or has_one, and other customized types.
type Association interface {
	Kind() reflect.Kind
	Interface() interface{}
	Constraint() (string, []interface{})
	InnerAssociations() InnerAssociations
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

// associationComposite adds the ability for a Association to
// have nested associations.
type associationComposite struct {
	innerAssociations InnerAssociations
}

func (a *associationComposite) InnerAssociations() InnerAssociations {
	return a.innerAssociations
}

// InnerAssociation is a struct that represents a deep level
// association. per example Song.Composer, Composer is an inner
// association for Song.
type InnerAssociation struct {
	Name   string
	Fields string
}

// InnerAssociations is a group of InnerAssociation.
type InnerAssociations []InnerAssociation

// AssociationSortable a type to be sortable.
type AssociationSortable interface {
	OrderBy() string
	Association
}

// AssociationCreatable is an association that
// can be created by executors.
type AssociationCreatable interface {
	// CreatableDependencies returns all dependencies for
	// this associations that needs to be created before.
	// A common example is a belongs to association type, where
	// the owner ID needs to be defined before it can be created.
	CreatableDependencies() []interface{}

	// Initialize is called before this association can be created
	// by an executor. Just as a way to setups fields or validates
	// any.
	Initialize() error
	Association
}

// AssociationCreatableStatement a association that defines
// create statements on database.
type AssociationCreatableStatement interface {
	Statements() []AssociationStatement
	Association
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
	field             reflect.StructField // an association field defined in model.
	modelType         reflect.Type        // the model type where this field is defined.
	modelValue        reflect.Value       // the model value where this field is defined.
	popTags           columns.Tags        // the tags defined in this association field.
	model             interface{}         // the model, owner of the association.
	innerAssociations InnerAssociations   // the data for the deep level associations.
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
