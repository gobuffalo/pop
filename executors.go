package pop

import (
	"fmt"

	. "github.com/markbates/pop/columns"
	"github.com/markbates/validate"
	uuid "github.com/satori/go.uuid"
)

func (c *Connection) Reload(model interface{}) error {
	sm := Model{Value: model}
	return c.Find(model, sm.ID())
}

func (q *Query) Exec() error {
	return q.Connection.timeFunc("Exec", func() error {
		sql, args := q.ToSQL(nil)
		Log(sql, args...)
		_, err := q.Connection.Store.Exec(sql, args...)
		return err
	})
}

func (c *Connection) ValidateAndSave(model interface{}, excludeColumns ...string) (*validate.Errors, error) {
	sm := &Model{Value: model}
	verrs, err := sm.validateSave(c)
	if err != nil {
		return verrs, err
	}
	if verrs.HasAny() {
		return verrs, nil
	}
	return verrs, c.Save(model, excludeColumns...)
}

var emptyUUID = uuid.Nil.String()

func (c *Connection) Save(model interface{}, excludeColumns ...string) (err error) {
	sm := &Model{Value: model}
	id := sm.ID()

	err = sm.beforeSave(c)
	if err != nil {
		return err
	}

	if id == 0 || (fmt.Sprint(id) == emptyUUID) {
		err = c.Create(model, excludeColumns...)
	} else {
		err = c.Update(model, excludeColumns...)
	}

	if err != nil {
		return err
	}

	return sm.afterSave(c)
}

func (c *Connection) ValidateAndCreate(model interface{}, excludeColumns ...string) (*validate.Errors, error) {
	sm := &Model{Value: model}
	verrs, err := sm.validateCreate(c)
	if err != nil {
		return verrs, err
	}
	if verrs.HasAny() {
		return verrs, nil
	}
	return verrs, c.Create(model, excludeColumns...)
}

func (c *Connection) Create(model interface{}, excludeColumns ...string) (err error) {
	return c.timeFunc("Create", func() error {
		sm := &Model{Value: model}
		err = sm.beforeCreate(c)
		if err != nil {
			return err
		}

		cols := ColumnsForStruct(model, sm.TableName())
		cols.Remove(excludeColumns...)

		sm.touchCreatedAt()
		sm.touchUpdatedAt()

		err = c.Dialect.Create(c.Store, sm, cols)
		if err != nil {
			return err
		}

		return sm.afterCreate(c)
	})
}

func (c *Connection) ValidateAndUpdate(model interface{}, excludeColumns ...string) (*validate.Errors, error) {
	sm := &Model{Value: model}
	verrs, err := sm.validateUpdate(c)
	if err != nil {
		return verrs, err
	}
	if verrs.HasAny() {
		return verrs, nil
	}
	return verrs, c.Update(model, excludeColumns...)
}

func (c *Connection) Update(model interface{}, excludeColumns ...string) (err error) {
	return c.timeFunc("Update", func() error {
		sm := &Model{Value: model}
		err = sm.beforeUpdate(c)
		if err != nil {
			return err
		}

		cols := ColumnsForStruct(model, sm.TableName())
		cols.Remove("id", "created_at")
		cols.Remove(excludeColumns...)

		sm.touchUpdatedAt()

		err = c.Dialect.Update(c.Store, sm, cols)
		if err != nil {
			return err
		}

		return sm.afterUpdate(c)
	})
}

func (c *Connection) Destroy(model interface{}) (err error) {
	return c.timeFunc("Destroy", func() error {
		sm := &Model{Value: model}
		err = sm.beforeDestroy(c)
		if err != nil {
			return err
		}

		err = c.Dialect.Destroy(c.Store, sm)
		if err != nil {
			return err
		}

		return sm.afterDestroy(c)
	})
}
