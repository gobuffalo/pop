package associations

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/markbates/pop/columns"
)

// AssociationsForStruct returns all associations for
// the struct specified. It takes into account tags
// associations like has_many, belongs_to, has_one.
// it throws an error when it finds a field that does
// not exist for a model.
func AssociationsForStruct(s interface{}, fields ...string) (Associations, error) {
	associations := Associations{}
	t, v := getModelDefinition(s)
	fields = trimFields(fields)

	// validate if fields contains a non existing field in struct.
	for _, f := range fields {
		if _, ok := t.FieldByName(f); !ok {
			return associations, fmt.Errorf("field %s does not exist in model %s", f, t.Name())
		}
	}

	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)

		// ignores those fields not included in fields list.
		if len(fields) > 0 && fieldIgnoredIn(fields, f.Name) {
			continue
		}

		tags := columns.TagsFor(f)

		// Find has_many association.
		tag := tags.Find("has_many")
		if !tag.Empty() {
			associations = append(associations, &hasManyAssociation{
				fieldName: f.Name,
				tableName: tag.Value,
				field:     f,
				value:     v.Field(i),
				ownerName: t.Name(),
				ownerID:   v.FieldByName("ID").Interface(),
				fkID:      tags.Find("fk_id").Value,
				orderBy:   tags.Find("order_by").Value,
			})
			continue
		}

		// Find belongs_to association.
		tag = tags.Find("belongs_to")
		if !tag.Empty() {
			fval := v.FieldByName(f.Name)
			associations = append(associations, &belongsToAssociation{
				ownerModel: fval,
				ownerType:  fval.Type(),
				ownerID:    v.FieldByName(fmt.Sprintf("%s%s", fval.Type().Name(), "ID")),
			})
			continue
		}

		// Find has_one association.
		tag = tags.Find("has_one")
		if !tag.Empty() {
			fval := v.FieldByName(f.Name)
			associations = append(associations, &hasOneAssociation{
				ownedModel: fval,
				ownedType:  fval.Type(),
				ownerID:    v.FieldByName("ID").Interface(),
				ownerName:  t.Name(),
				fkID:       tags.Find("fk_id").Value,
			})
		}
	}
	return associations, nil
}

func getModelDefinition(s interface{}) (reflect.Type, reflect.Value) {
	v := reflect.ValueOf(s)
	v = reflect.Indirect(v)
	t := v.Type()
	return t, v
}

func trimFields(fields []string) []string {
	trimFields := []string{}
	for _, f := range fields {
		if strings.TrimSpace(f) != "" {
			trimFields = append(trimFields, strings.TrimSpace(f))
		}
	}
	return trimFields
}

func fieldIgnoredIn(fields []string, field string) bool {
	for _, f := range fields {
		if f == field {
			return false
		}
	}
	return true
}
