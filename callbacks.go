package pop

import (
	"reflect"

	"github.com/pkg/errors"
)

func (m *Model) runCallbacks(c *Connection, name string) error {
	rv := reflect.ValueOf(m.Value)
	mv := rv.MethodByName(name)
	if mv.IsValid() {
		typ := mv.Type()
		if typ.NumIn() == 1 && typ.In(0) == reflect.TypeOf(c) {
			if typ.NumOut() != 1 {
				return errors.Errorf("%s function should return error!", name)
			}
			out := mv.Call([]reflect.Value{reflect.ValueOf(c)})
			if !out[0].IsNil() {
				return out[0].Interface().(error)
			}
		} else {
			return errors.Errorf("%s function should take 1 argument of type '*pop.Connection'", name)
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
