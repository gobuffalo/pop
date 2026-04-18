package associations

import (
	"cmp"
	"fmt"
	"reflect"

	"github.com/gobuffalo/flect"
	"github.com/gobuffalo/nulls"

	"github.com/gobuffalo/pop/v6/columns"
)

// belongsToAssociation is the implementation for the belongs_to association type in a model.
type belongsToAssociation struct {
	*associationSkippable
	*associationComposite

	ownerModel reflect.Value
	ownerType  reflect.Type
	ownerID    reflect.Value
	primaryID  string
	ownedModel any

	primaryTableID string
}

func init() {
	associationBuilders["belongs_to"] = belongsToAssociationBuilder
}

func belongsToAssociationBuilder(p associationParams) (Association, error) {
	ownerVal := p.modelValue.FieldByName(p.field.Name)
	tags := p.popTags
	primaryIDField := cmp.Or(tags.Find("primary_id").Value, "ID")
	ownerIDField := fmt.Sprintf("%s%s", p.field.Name, "ID")

	if dbTag := tags.Find("fk_id").Value; dbTag != "" {
		if _, found := p.modelType.FieldByName(dbTag); !found {
			t := p.modelValue.Type()
			for i := range t.NumField() {
				f := t.Field(i)
				if f.Tag.Get("db") == dbTag {
					ownerIDField = f.Name
					break
				}
			}
		} else {
			ownerIDField = dbTag
		}
	}

	// belongs_to requires an holding field for the foreign model ID.
	if _, found := p.modelType.FieldByName(ownerIDField); !found {
		return nil, fmt.Errorf("there is no '%s' defined in model '%s'", ownerIDField, p.modelType.Name())
	}

	// If ownerIDField is nil, this association will be skipped.
	var skipped bool
	f := p.modelValue.FieldByName(ownerIDField)
	if fieldIsNil(f) || IsZeroOfUnderlyingType(f.Interface()) {
		skipped = true
	}
	// associated model
	ownerPk := "id"
	if primaryIDField != "ID" {
		ownerModel := reflect.Indirect(ownerVal)
		ownerPrimaryField, found := ownerModel.Type().FieldByName(primaryIDField)
		if !found {
			return nil, fmt.Errorf(
				"there is no primary field '%s' defined in model '%s'",
				primaryIDField,
				ownerModel.Type(),
			)
		}
		ownerPTags := columns.TagsFor(ownerPrimaryField)
		ownerPk = cmp.Or(ownerPTags.Find("db").Value, flect.Underscore(ownerPrimaryField.Name))
	}

	return &belongsToAssociation{
		ownerModel: ownerVal,
		ownerType:  ownerVal.Type(),
		ownerID:    f,
		primaryID:  primaryIDField,
		ownedModel: p.model,
		associationSkippable: &associationSkippable{
			skipped: skipped,
		},
		associationComposite: &associationComposite{innerAssociations: p.innerAssociations},
		primaryTableID:       ownerPk,
	}, nil
}

func (b *belongsToAssociation) Kind() reflect.Kind {
	if b.ownerType.Kind() == reflect.Pointer {
		return b.ownerType.Elem().Kind()
	}
	return b.ownerType.Kind()
}

func (b *belongsToAssociation) Interface() any {
	if b.ownerModel.Kind() == reflect.Pointer {
		val := reflect.New(b.ownerType.Elem())
		b.ownerModel.Set(val)
		return b.ownerModel.Interface()
	}
	return b.ownerModel.Addr().Interface()
}

// Constraint returns the content for a where clause, and the args
// needed to execute it.
func (b *belongsToAssociation) Constraint() (string, []any) {
	return b.primaryTableID + " = ?", []any{b.ownerID.Interface()}
}

func (b *belongsToAssociation) BeforeInterface() any {
	// if the owner field is set, don't try to create the association to prevent conflicts.
	if !b.skipped {
		return nil
	}

	m := b.ownerModel
	if m.Kind() == reflect.Pointer && !m.IsNil() {
		m = b.ownerModel.Elem()
	}

	if IsZeroOfUnderlyingType(m.Interface()) {
		return nil
	}

	return m.Addr().Interface()
}

func (b *belongsToAssociation) BeforeSetup() error {
	ownerID := reflect.Indirect(reflect.ValueOf(b.ownerModel.Interface())).FieldByName("ID")
	toSet := b.ownerID
	switch b.ownerID.Type().Name() {
	case "NullUUID":
		b.ownerID.FieldByName("Valid").Set(reflect.ValueOf(true))
		toSet = b.ownerID.FieldByName("UUID")
	}

	if toSet.CanSet() {
		if n := nulls.New(toSet.Interface()); n != nil {
			toSet.Set(reflect.ValueOf(n.Parse(ownerID.Interface())))
		} else if toSet.Kind() == reflect.Pointer {
			toSet.Set(ownerID.Addr())
		} else {
			toSet.Set(ownerID)
		}
		return nil
	}
	return fmt.Errorf("could not set '%s' to '%s'", ownerID, toSet)
}
