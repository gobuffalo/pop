package pop

import (
	"fmt"
	"reflect"
	"sync"
	"time"

	"github.com/markbates/inflect"
	"github.com/pkg/errors"
)

var tableMap = map[string]string{}
var tableMapMu = sync.RWMutex{}

// MapTableName allows for the customize table mapping
// between a name and the database. For example the value
// `User{}` will automatically map to "users".
// MapTableName would allow this to change.
//
//	m := &pop.Model{Value: User{}}
//	m.TableName() // "users"
//
//	pop.MapTableName("user", "people")
//	m = &pop.Model{Value: User{}}
//	m.TableName() // "people"
func MapTableName(name string, tableName string) {
	defer tableMapMu.Unlock()
	tableMapMu.Lock()
	tableMap[name] = tableName
}

type Value interface{}

// Model is used throughout Pop to wrap the end user interface
// that is passed in to many functions.
type Model struct {
	Value
	tableName string
}

// func (m *Model) Validate(*Connection) (*validate.Errors, error) {
// 	return validate.NewErrors(), nil
// }
//
// func (m *Model) ValidateNew(*Connection) (*validate.Errors, error) {
// 	return validate.NewErrors(), nil
// }
//
// func (m *Model) ValidateUpdate(*Connection) (*validate.Errors, error) {
// 	return validate.NewErrors(), nil
// }

// ID returns the ID of the Model. All models must have an `ID` field this is
// of type `int`.
func (m *Model) ID() int {
	fbn, err := m.fieldByName("ID")
	if err != nil {
		return 0
	}
	return int(fbn.Int())
}

// TableName returns the corresponding name of the underlying database table
// for a given `Model`. See also `MapTableName` to change the default name of
// the table.
func (m *Model) TableName() string {
	if m.tableName != "" {
		return m.tableName
	}

	t := reflect.TypeOf(m.Value)
	name := m.typeName(t)

	defer tableMapMu.Unlock()
	tableMapMu.Lock()

	if tableMap[name] == "" {
		m.tableName = inflect.Tableize(name)
		tableMap[name] = m.tableName
	}
	return tableMap[name]
}

func (m *Model) typeName(t reflect.Type) string {
	kind := t.Kind().String()
	switch kind {
	case "ptr":
		st := reflect.ValueOf(m.Value).Elem()
		return m.typeName(st.Type())
	case "string":
		return m.Value.(string)
	case "slice":
		if t.Elem().Kind().String() == "ptr" {
			return m.typeName(t.Elem().Elem())
		}
		return t.Elem().Name()
	default:
		return t.Name()
	}
}

func (m *Model) fieldByName(s string) (reflect.Value, error) {
	el := reflect.ValueOf(m.Value).Elem()
	fbn := el.FieldByName(s)
	if !fbn.IsValid() {
		return fbn, errors.Errorf("Model does not have a field named %s", s)
	}
	return fbn, nil
}

func (m *Model) associationName() string {
	tn := inflect.Singularize(m.TableName())
	return fmt.Sprintf("%s_id", tn)
}

func (m *Model) setID(i int) {
	fbn, err := m.fieldByName("ID")
	if err == nil {
		fbn.SetInt(int64(i))
	}
}

func (m *Model) touchCreatedAt() {
	fbn, err := m.fieldByName("CreatedAt")
	if err == nil {
		fbn.Set(reflect.ValueOf(time.Now()))
	}
}

func (m *Model) touchUpdatedAt() {
	fbn, err := m.fieldByName("UpdatedAt")
	if err == nil {
		fbn.Set(reflect.ValueOf(time.Now()))
	}
}
