package associations

import (
	"fmt"
	"reflect"

	"github.com/markbates/pop/columns"
)

// AssociationsForStruct returns all associations for
// the struct specified. It takes into account tags
// associations like has_many, belongs_to, has_one.
func AssociationsForStruct(s interface{}) Associations {
	associations := Associations{}
	t, v := getModelDefinition(s)

	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
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
		}
	}
	return associations
}

func getModelDefinition(s interface{}) (reflect.Type, reflect.Value) {
	t := reflect.TypeOf(s)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	if t.Kind() == reflect.Slice {
		t = t.Elem()
		if t.Kind() == reflect.Ptr {
			t = t.Elem()
		}
	}

	v := reflect.ValueOf(s)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	return t, v
}
