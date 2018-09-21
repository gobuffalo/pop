package pop

import (
	"fmt"
	"reflect"
	"sync"
	"time"

	"github.com/gobuffalo/flect"
	nflect "github.com/gobuffalo/flect/name"
	"github.com/gobuffalo/uuid"
	"github.com/pkg/errors"
)

var tableMap = map[string]string{}
var tableMapMu = sync.RWMutex{}

// Value is the contents of a `Model`.
type Value interface{}

type modelIterable func(*Model) error

// Model is used throughout Pop to wrap the end user interface
// that is passed in to many functions.
type Model struct {
	Value
	tableName string
	As        string
}

// ID returns the ID of the Model. All models must have an `ID` field this is
// of type `int`,`int64` or of type `uuid.UUID`.
func (m *Model) ID() interface{} {
	fbn, err := m.fieldByName("ID")
	if err != nil {
		return 0
	}
	if m.PrimaryKeyType() == "UUID" {
		return fbn.Interface().(uuid.UUID).String()
	}
	return fbn.Interface()
}

// PrimaryKeyType gives the primary key type of the `Model`.
func (m *Model) PrimaryKeyType() string {
	fbn, err := m.fieldByName("ID")
	if err != nil {
		return "int"
	}
	return fbn.Type().Name()
}

// TableNameAble interface allows for the customize table mapping
// between a name and the database. For example the value
// `User{}` will automatically map to "users". Implementing `TableNameAble`
// would allow this to change to be changed to whatever you would like.
type TableNameAble interface {
	TableName() string
}

// TableName returns the corresponding name of the underlying database table
// for a given `Model`. See also `TableNameAble` to change the default name of the table.
func (m *Model) TableName() string {
	if s, ok := m.Value.(string); ok {
		return s
	}
	if n, ok := m.Value.(TableNameAble); ok {
		return n.TableName()
	}

	if m.tableName != "" {
		return m.tableName
	}

	t := reflect.TypeOf(m.Value)
	name := m.typeName(t)

	defer tableMapMu.Unlock()
	tableMapMu.Lock()

	if tableMap[name] == "" {
		m.tableName = nflect.Tableize(name)
		tableMap[name] = m.tableName
	}
	return tableMap[name]
}

func (m *Model) typeName(t reflect.Type) string {
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	switch t.Kind() {
	case reflect.Slice, reflect.Array:
		el := t.Elem()
		if el.Kind() == reflect.Ptr {
			el = el.Elem()
		}

		// validates if the elem of slice or array implements TableNameAble interface.
		tableNameAble := (*TableNameAble)(nil)
		if el.Implements(reflect.TypeOf(tableNameAble).Elem()) {
			v := reflect.New(el)
			out := v.MethodByName("TableName").Call([]reflect.Value{})
			name := out[0].String()
			if tableMap[el.Name()] == "" {
				tableMap[el.Name()] = name
			}
		}

		return el.Name()
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
	tn := flect.Singularize(m.TableName())
	return fmt.Sprintf("%s_id", tn)
}

func (m *Model) setID(i interface{}) {
	fbn, err := m.fieldByName("ID")
	if err == nil {
		v := reflect.ValueOf(i)
		switch fbn.Kind() {
		case reflect.Int, reflect.Int64:
			fbn.SetInt(v.Int())
		default:
			fbn.Set(reflect.ValueOf(i))
		}
	}
}

func (m *Model) touchCreatedAt() {
	fbn, err := m.fieldByName("CreatedAt")
	if err == nil {
		now := time.Now()
		switch fbn.Kind() {
		case reflect.Int, reflect.Int64:
			fbn.SetInt(now.Unix())
		default:
			fbn.Set(reflect.ValueOf(now))
		}
	}
}

func (m *Model) touchUpdatedAt() {
	fbn, err := m.fieldByName("UpdatedAt")
	if err == nil {
		now := time.Now()
		switch fbn.Kind() {
		case reflect.Int, reflect.Int64:
			fbn.SetInt(now.Unix())
		default:
			fbn.Set(reflect.ValueOf(now))
		}
	}
}

func (m *Model) whereID() string {
	return fmt.Sprintf("%s.id = ?", m.TableName())
}

func (m *Model) whereNamedID() string {
	return fmt.Sprintf("%s.id = :id", m.TableName())
}

func (m *Model) isSlice() bool {
	v := reflect.Indirect(reflect.ValueOf(m.Value))
	return v.Kind() == reflect.Slice || v.Kind() == reflect.Array
}

func (m *Model) iterate(fn modelIterable) error {
	if m.isSlice() {
		v := reflect.Indirect(reflect.ValueOf(m.Value))
		for i := 0; i < v.Len(); i++ {
			val := v.Index(i)
			newModel := &Model{Value: val.Addr().Interface()}
			err := fn(newModel)

			if err != nil {
				return err
			}
		}
		return nil
	}

	return fn(m)
}
