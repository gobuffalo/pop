package associations

import (
	"fmt"
	"reflect"

	"github.com/markbates/pop/nulls"

	"github.com/markbates/inflect"
)

// hasManyAssociation is the implementation for the has_many
// association type in a model.
type hasManyAssociation struct {
	tableName string
	field     reflect.StructField
	value     reflect.Value
	ownerName string
	ownerID   interface{}
	owner     interface{}
	fkID      string
	orderBy   string
	skipped   bool
}

func init() {
	associationBuilders["has_many"] = hasManyAssociationBuilder
}

func hasManyAssociationBuilder(p associationParams) (Association, error) {
	// Validates if ownerID is nil, this association will be skipped.
	var skipped bool
	ownerID := p.modelValue.FieldByName("ID")
	if fieldIsNil(ownerID) {
		skipped = true
	}

	return &hasManyAssociation{
		owner:     p.model,
		tableName: p.popTags.Find("has_many").Value,
		field:     p.field,
		value:     p.modelValue.FieldByName(p.field.Name),
		ownerName: p.modelType.Name(),
		ownerID:   ownerID.Interface(),
		fkID:      p.popTags.Find("fk_id").Value,
		orderBy:   p.popTags.Find("order_by").Value,
		skipped:   skipped,
	}, nil
}

func (a *hasManyAssociation) Kind() reflect.Kind {
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

// Constraint returns the content for a where clause, and the args
// needed to execute it.
func (a *hasManyAssociation) Constraint() (string, []interface{}) {
	tn := inflect.Underscore(a.ownerName)
	condition := fmt.Sprintf("%s_id = ?", tn)
	if a.fkID != "" {
		condition = fmt.Sprintf("%s = ?", a.fkID)
	}
	return condition, []interface{}{a.ownerID}
}

func (a *hasManyAssociation) OrderBy() string {
	return a.orderBy
}

// if it is skipped if the only dependency is present:
// owner ID.
func (a *hasManyAssociation) Dependencies() []interface{} {
	ownerID := reflect.Indirect(reflect.ValueOf(a.owner)).FieldByName("ID").Interface()
	if a.skipped || isZero(ownerID) {
		return []interface{}{a.owner}
	}
	return []interface{}{}
}

// SetValue for has many association loop over every item in the
// value associated and set his foreign key reference with the
// val passed as parameter.
func (a *hasManyAssociation) SetValue(val []interface{}) error {
	var ownerID interface{}
	if a.skipped {
		ownerID = reflect.Indirect(reflect.ValueOf(val[0])).FieldByName("ID").Interface()
	} else {
		ownerID = reflect.Indirect(reflect.ValueOf(a.owner)).FieldByName("ID").Interface()
	}

	v := a.value
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	for i := 0; i < v.Len(); i++ {
		fval := v.Index(i).FieldByName(a.ownerName + "ID")
		if fval.CanSet() {
			if n := nulls.New(fval.Interface()); n != nil {
				fval.Set(reflect.ValueOf(n.Parse(ownerID)))
			} else {
				fval.Set(reflect.ValueOf(ownerID))
			}
		} else {
			return fmt.Errorf("could not set '%s' in '%s'", ownerID, fval)
		}
	}
	return nil
}

func (a *hasManyAssociation) Skipped() bool {
	return a.skipped
}
