package associations

import (
	"fmt"
	"reflect"

	"github.com/gobuffalo/flect"
	"github.com/gobuffalo/nulls"
	"github.com/jmoiron/sqlx"
)

// hasManyAssociation is the implementation for the has_many
// association type in a model.
type hasManyAssociation struct {
	tableName string
	field     reflect.StructField
	value     reflect.Value
	ownerName string
	ownerID   any
	owner     any
	fkID      string
	orderBy   string
	*associationSkipable
	*associationComposite
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
		associationSkipable: &associationSkipable{
			skipped: skipped,
		},
		associationComposite: &associationComposite{innerAssociations: p.innerAssociations},
	}, nil
}

func (a *hasManyAssociation) Kind() reflect.Kind {
	if a.field.Type.Kind() == reflect.Pointer {
		return a.field.Type.Elem().Kind()
	}
	return a.field.Type.Kind()
}

func (a *hasManyAssociation) Interface() any {
	if a.value.Kind() == reflect.Pointer {
		val := reflect.New(a.field.Type.Elem())
		a.value.Set(val)
		return a.value.Interface()
	}

	// This piece of code clears a slice in case it is filled with elements.
	if a.value.Kind() == reflect.Slice || a.value.Kind() == reflect.Array {
		valPointer := a.value.Addr()
		valPointer.Elem().Set(reflect.MakeSlice(valPointer.Type().Elem(), 0, valPointer.Elem().Cap()))
		return valPointer.Interface()
	}

	return a.value.Addr().Interface()
}

// Constraint returns the content for a where clause, and the args
// needed to execute it.
func (a *hasManyAssociation) Constraint() (string, []any) {
	tn := flect.Underscore(a.ownerName)
	condition := tn + "_id = ?"
	if a.fkID != "" {
		condition = a.fkID + " = ?"
	}
	return condition, []any{a.ownerID}
}

func (a *hasManyAssociation) OrderBy() string {
	return a.orderBy
}

func (a *hasManyAssociation) AfterInterface() any {
	if a.value.Kind() == reflect.Pointer {
		return a.value.Interface()
	}
	return a.value.Addr().Interface()
}

func (a *hasManyAssociation) AfterSetup() error {
	ownerID := reflect.Indirect(reflect.ValueOf(a.owner)).FieldByName("ID").Interface()

	v := a.value
	if v.Kind() == reflect.Pointer {
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
			return fmt.Errorf(
				"could not set field '%s' in table '%s' to value '%s' for 'has_many' relation",
				a.ownerName+"ID",
				a.tableName,
				ownerID,
			)
		}
	}
	return nil
}

func (a *hasManyAssociation) AfterProcess() AssociationStatement {
	v := a.value
	if v.Kind() == reflect.Pointer {
		v = v.Elem()
	}

	belongingIDFieldName := "ID"

	ownerIDFieldName := "ID"
	ownerID := reflect.Indirect(reflect.ValueOf(a.owner)).FieldByName(ownerIDFieldName).Interface()

	var ids []any

	for i := 0; i < v.Len(); i++ {
		id := v.Index(i).FieldByName(belongingIDFieldName).Interface()
		if !IsZeroOfUnderlyingType(id) {
			ids = append(ids, id)
		}
	}
	if len(ids) == 0 {
		return AssociationStatement{
			Statement: "",
			Args:      []any{},
		}
	}

	fk := a.fkID
	if fk == "" {
		fk = flect.Underscore(a.ownerName) + "_id"
	}

	// This will be used to update all of our owned models' foreign keys to our ID.
	ret := fmt.Sprintf("UPDATE %s SET %s = ? WHERE %s in (?);", a.tableName, fk, belongingIDFieldName)

	update, args, err := sqlx.In(ret, ownerID, ids)
	if err != nil {
		return AssociationStatement{
			Statement: "",
			Args:      []any{},
		}
	}

	return AssociationStatement{
		Statement: update,
		Args:      args,
	}
}
