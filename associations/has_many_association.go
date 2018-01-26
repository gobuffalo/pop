package associations

import (
	"fmt"
	"reflect"
	"strings"
)

// hasManyAssociation is the implementation for the has_many
// association type in a model.
type hasManyAssociation struct {
	fieldName string
	tableName string
	field     reflect.StructField
	value     reflect.Value
	ownerName string
	ownerID   interface{}
}

func (a *hasManyAssociation) TableName() string {
	return a.tableName
}

func (a *hasManyAssociation) FieldName() string {
	return a.fieldName
}

func (a *hasManyAssociation) Type() reflect.Kind {
	if a.field.Type.Kind() == reflect.Ptr {
		return a.field.Type.Elem().Kind()
	}
	return a.field.Type.Kind()
}

func (a *hasManyAssociation) Interface() interface{} {
	if a.value.Kind() == reflect.Ptr {
		val := reflect.New(a.field.Type.Elem())
		a.value.Set(val)
		return a.value.Interface()
	}
	return a.value.Addr().Interface()
}

// SQLConstraint returns the content for a where clause, and the args
// needed to execute it.
func (a *hasManyAssociation) SQLConstraint() (string, []interface{}) {
	tn := strings.ToLower(a.ownerName)
	return fmt.Sprintf("%s_id = ?", tn), []interface{}{a.ownerID}
}
