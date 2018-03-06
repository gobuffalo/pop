package associations

import (
	"fmt"
	"reflect"
	"time"

	"github.com/markbates/inflect"
)

type manyToManyAssociation struct {
	fieldType           reflect.Type
	fieldValue          reflect.Value
	model               reflect.Value
	manyToManyTableName string
	owner               interface{}
	fkID                string
	orderBy             string
	*associationSkipable
	*associationComposite
}

func init() {
	associationBuilders["many_to_many"] = func(p associationParams) (Association, error) {
		// Validates if model.ID is nil, this association will be skipped.
		var skipped bool
		model := p.modelValue
		if fieldIsNil(model.FieldByName("ID")) {
			skipped = true
		}

		return &manyToManyAssociation{
			fieldType:           p.modelValue.FieldByName(p.field.Name).Type(),
			fieldValue:          p.modelValue.FieldByName(p.field.Name),
			owner:               p.model,
			model:               model,
			manyToManyTableName: p.popTags.Find("many_to_many").Value,
			fkID:                p.popTags.Find("fk_id").Value,
			orderBy:             p.popTags.Find("order_by").Value,
			associationSkipable: &associationSkipable{
				skipped: skipped,
			},
			associationComposite: &associationComposite{innerAssociations: p.innerAssociations},
		}, nil
	}
}

func (m *manyToManyAssociation) Kind() reflect.Kind {
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

// Constraint returns the content for a where clause, and the args
// needed to execute it.
func (m *manyToManyAssociation) Constraint() (string, []interface{}) {
	modelColumnID := fmt.Sprintf("%s%s", inflect.Underscore(m.model.Type().Name()), "_id")

	var columnFieldID string
	i := reflect.Indirect(m.fieldValue)
	if i.Kind() == reflect.Slice || i.Kind() == reflect.Array {
		t := i.Type().Elem()
		columnFieldID = fmt.Sprintf("%s%s", inflect.Underscore(t.Name()), "_id")
	} else {
		columnFieldID = fmt.Sprintf("%s%s", inflect.Underscore(i.Type().Name()), "_id")
	}

	if m.fkID != "" {
		columnFieldID = m.fkID
	}

	subQuery := fmt.Sprintf("select %s from %s where %s = ?", columnFieldID, m.manyToManyTableName, modelColumnID)
	modelIDValue := m.model.FieldByName("ID").Interface()

	return fmt.Sprintf("id in (%s)", subQuery), []interface{}{modelIDValue}
}

func (m *manyToManyAssociation) OrderBy() string {
	return m.orderBy
}

func (m *manyToManyAssociation) CreatableDependencies() []interface{} {
	dependencies := []interface{}{}
	if m.fieldValue.Kind() == reflect.Ptr {
		dependencies = append(dependencies, m.fieldValue.Interface())
	} else {
		dependencies = append(dependencies, m.fieldValue.Addr().Interface())
	}

	modelID := reflect.Indirect(reflect.ValueOf(m.owner)).FieldByName("ID")
	if fieldIsNil(modelID) || isZero(modelID.Interface()) {
		dependencies = append(dependencies, m.owner)
	}
	return dependencies
}

func (m *manyToManyAssociation) Initialize() error {
	return nil
}

func (m *manyToManyAssociation) Statements() []AssociationStatement {
	statements := []AssociationStatement{}

	modelColumnID := fmt.Sprintf("%s%s", inflect.Underscore(m.model.Type().Name()), "_id")
	var columnFieldID string
	i := reflect.Indirect(m.fieldValue)
	if i.Kind() == reflect.Slice || i.Kind() == reflect.Array {
		t := i.Type().Elem()
		columnFieldID = fmt.Sprintf("%s%s", inflect.Underscore(t.Name()), "_id")
	} else {
		columnFieldID = fmt.Sprintf("%s%s", inflect.Underscore(i.Type().Name()), "_id")
	}

	for i := 0; i < m.fieldValue.Len(); i++ {
		v := m.fieldValue.Index(i)
		manyIDValue := v.FieldByName("ID").Interface()
		modelIDValue := m.model.FieldByName("ID").Interface()
		stm := "INSERT INTO %s (%s,%s,%s,%s) VALUES(?,?,?,?)"

		associationStm := AssociationStatement{
			Statement: fmt.Sprintf(stm, m.manyToManyTableName, modelColumnID, columnFieldID, "created_at", "updated_at"),
			Args:      []interface{}{modelIDValue, manyIDValue, time.Now(), time.Now()},
		}

		statements = append(statements, associationStm)
	}

	return statements
}
