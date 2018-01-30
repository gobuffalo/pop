package associations

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/markbates/inflect"
)

type manyToManyAssociation struct {
	fieldType           reflect.Type
	fieldValue          reflect.Value
	model               reflect.Value
	manyToManyTableName string
	fkID                string
}

func init() {
	associationBuilders["many_to_many"] = func(p associationParams) (Association, error) {
		return &manyToManyAssociation{
			fieldType:           p.modelValue.FieldByName(p.field.Name).Type(),
			fieldValue:          p.modelValue.FieldByName(p.field.Name),
			model:               p.modelValue,
			manyToManyTableName: p.popTags.Find("many_to_many").Value,
			fkID:                p.popTags.Find("fk_id").Value,
		}, nil
	}
}

func (m *manyToManyAssociation) TableName() string {
	method := m.fieldValue.MethodByName("TableName")
	if method.IsValid() {
		out := method.Call([]reflect.Value{})
		return out[0].String()
	}
	return inflect.Tableize(m.fieldType.Name())
}

func (m *manyToManyAssociation) Type() reflect.Kind {
	return m.fieldType.Kind()
}

func (m *manyToManyAssociation) Interface() interface{} {
	if m.fieldValue.Kind() == reflect.Ptr {
		val := reflect.New(m.fieldType.Elem())
		m.fieldValue.Set(val)
		return m.fieldValue.Interface()
	}
	return m.fieldValue.Addr().Interface()
}

// SQLConstraint returns the content for a where clause, and the args
// needed to execute it.
func (m *manyToManyAssociation) SQLConstraint() (string, []interface{}) {
	modelColumnID := fmt.Sprintf("%s%s", strings.ToLower(m.model.Type().Name()), "_id")

	var columnFieldID string
	i := reflect.Indirect(m.fieldValue)
	if i.Kind() == reflect.Slice {
		t := i.Type().Elem()
		columnFieldID = fmt.Sprintf("%s%s", strings.ToLower(t.Name()), "_id")
	} else {
		columnFieldID = fmt.Sprintf("%s%s", strings.ToLower(i.Type().Name()), "_id")
	}

	if m.fkID != "" {
		columnFieldID = m.fkID
	}

	subQuery := fmt.Sprintf("select %s from %s where %s = ?", columnFieldID, m.manyToManyTableName, modelColumnID)
	modelIDValue := m.model.FieldByName("ID").Interface()

	return fmt.Sprintf("id in (%s)", subQuery), []interface{}{modelIDValue}
}
