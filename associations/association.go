package associations

import (
	"reflect"
)

// Association represents a definition of a model association
// field. It can represent a association of the type has_many
// belongs_to or has_one, and other customized types.
type Association interface {
	TableName() string
	FieldName() string
	Type() reflect.Kind
	Interface() interface{}
	SQLConstraint() (string, []interface{})
}

// Associations a group of model associations.
type Associations []Association
