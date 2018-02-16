package associations

import (
	"fmt"
	"reflect"

	"github.com/markbates/inflect"
	"github.com/markbates/pop/nulls"
)

// belongsToAssociation is the implementation for the belongs_to
// association type in a model.
type belongsToAssociation struct {
	ownerModel reflect.Value
	ownerType  reflect.Type
	ownerID    reflect.Value
	owner      interface{}
	skipped    bool
}

func init() {
	associationBuilders["belongs_to"] = belongsToAssociationBuilder
}

func belongsToAssociationBuilder(p associationParams) (Association, error) {
	fval := p.modelValue.FieldByName(p.field.Name)
	ownerIDField := fmt.Sprintf("%s%s", inflect.Capitalize(fval.Type().Name()), "ID")

	if _, found := p.modelType.FieldByName(ownerIDField); !found {
		return nil, fmt.Errorf("there is no '%s' defined in model '%s'", ownerIDField, p.modelType.Name())
	}

	// Validates if ownerIDField is nil, this association will be skipped.
	var skipped bool
	f := p.modelValue.FieldByName(ownerIDField)
	if fieldIsNil(f) {
		skipped = true
	}

	return &belongsToAssociation{
		ownerModel: fval,
		ownerType:  fval.Type(),
		ownerID:    f,
		owner:      p.model,
		skipped:    skipped,
	}, nil
}

func (b *belongsToAssociation) Kind() reflect.Kind {
	if b.ownerType.Kind() == reflect.Ptr {
		return b.ownerType.Elem().Kind()
	}
	return b.ownerType.Kind()
}

func (b *belongsToAssociation) Interface() interface{} {
	if b.skipped {
		return b.owner
	}

	if b.ownerModel.Kind() == reflect.Ptr {
		val := reflect.New(b.ownerType.Elem())
		b.ownerModel.Set(val)
		return b.ownerModel.Interface()
	}
	return b.ownerModel.Addr().Interface()
}

// Constraint returns the content for a where clause, and the args
// needed to execute it.
func (b *belongsToAssociation) Constraint() (string, []interface{}) {
	return "id = ?", []interface{}{b.ownerID.Interface()}
}

func (b *belongsToAssociation) Dependencies() []interface{} {
	if b.skipped {
		if b.ownerModel.Kind() == reflect.Ptr {
			return []interface{}{b.ownerModel.Interface()}
		}
		return []interface{}{b.ownerModel.Addr().Interface()}
	}
	return []interface{}{}
}

func (b *belongsToAssociation) SetValue(i []interface{}) error {
	ownerID := reflect.Indirect(reflect.ValueOf(i[0])).FieldByName("ID").Interface()

	if b.ownerID.CanSet() {
		if n := nulls.New(b.ownerID.Interface()); n != nil {
			b.ownerID.Set(reflect.ValueOf(n.Parse(ownerID)))
		} else {
			b.ownerID.Set(reflect.ValueOf(ownerID))
		}
	} else {
		return fmt.Errorf("could not set '%s' to '%s'", ownerID, b.ownerID)
	}

	return nil
}

func (b *belongsToAssociation) Skipped() bool {
	return b.skipped
}
