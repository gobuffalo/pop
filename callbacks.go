package pop

import (
	"reflect"

	"github.com/pkg/errors"
)

func (m *Model) runCallbacks(c *Connection, name string) error {
	rv := reflect.ValueOf(m.Value)
	mv := rv.MethodByName(name)
	if mv.IsValid() {
		if mv.Type().NumOut() != 1 {
			return errors.Errorf("%s does not have the correct method signature!", name)
		}
		out := mv.Call([]reflect.Value{reflect.ValueOf(c)})
		if !out[0].IsNil() {
			return out[0].Interface().(error)
		}
	}
	return nil
}

func (m *Model) beforeSave(c *Connection) error {
	return m.runCallbacks(c, "BeforeSave")
}

func (m *Model) beforeCreate(c *Connection) error {
	return m.runCallbacks(c, "BeforeCreate")
}

func (m *Model) beforeUpdate(c *Connection) error {
	return m.runCallbacks(c, "BeforeUpdate")
}

func (m *Model) beforeDestroy(c *Connection) error {
	return m.runCallbacks(c, "BeforeDestroy")
}

func (m *Model) afterDestroy(c *Connection) error {
	return m.runCallbacks(c, "AfterDestroy")
}

func (m *Model) afterUpdate(c *Connection) error {
	return m.runCallbacks(c, "AfterUpdate")
}

func (m *Model) afterCreate(c *Connection) error {
	return m.runCallbacks(c, "AfterCreate")
}

func (m *Model) afterSave(c *Connection) error {
	return m.runCallbacks(c, "AfterSave")
}
