package pop

import (
	"fmt"
	"reflect"
	"time"

	"github.com/markbates/going/validate"
	"github.com/markbates/inflect"
)

var tableMap = map[string]string{}

func MapTableName(name string, tableName string) {
	tableMap[name] = tableName
}

type Value interface{}

type Model struct {
	Value
	tableName string
}

func (m *Model) FieldByName(s string) (reflect.Value, error) {
	el := reflect.ValueOf(m.Value).Elem()
	fbn := el.FieldByName(s)
	if !fbn.IsValid() {
		return fbn, fmt.Errorf("Model does not have a field named %s", s)
	}
	return fbn, nil
}

func (m *Model) Validate(*Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

func (m *Model) ValidateNew(*Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

func (m *Model) ValidateUpdate(*Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

func (m *Model) ID() int {
	fbn, err := m.FieldByName("ID")
	if err != nil {
		return 0
	}
	return int(fbn.Int())
}

func (m *Model) SetID(i int) {
	fbn, err := m.FieldByName("ID")
	if err == nil {
		fbn.SetInt(int64(i))
	}
}

func (m *Model) TouchCreatedAt() {
	fbn, err := m.FieldByName("CreatedAt")
	if err == nil {
		fbn.Set(reflect.ValueOf(time.Now()))
	}
}

func (m *Model) TouchUpdatedAt() {
	fbn, err := m.FieldByName("UpdatedAt")
	if err == nil {
		fbn.Set(reflect.ValueOf(time.Now()))
	}
}

func (m *Model) TableName() string {
	if m.tableName != "" {
		return m.tableName
	}

	t := reflect.TypeOf(m.Value)
	kind := t.Kind().String()
	var name string
	switch kind {
	case "ptr":
		st := reflect.ValueOf(m.Value).Elem()
		name = st.Type().Name()
	case "string":
		return m.Value.(string)
	default:
		name = t.Name()
	}
	if tableMap[name] == "" {
		m.tableName = inflect.Tableize(name)
		tableMap[name] = m.tableName
	}
	return tableMap[name]
}
