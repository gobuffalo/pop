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

func (c *Connection) Save(model interface{}, excludeColumns ...string) error {
	sm := &Model{Value: model}
	id := sm.ID()
	if id == 0 || (fmt.Sprint(id) == emptyUUID) {
		return c.Create(model, excludeColumns...)
	} else {
		return c.Update(model, excludeColumns...)
	}
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

func (c *Connection) Create(model interface{}, excludeColumns ...string) error {
	return c.timeFunc("Create", func() error {
		sm := &Model{Value: model}

		cols := ColumnsForStruct(model, sm.TableName())
		cols.Remove(excludeColumns...)

		sm.touchCreatedAt()
		sm.touchUpdatedAt()

		return c.Dialect.Create(c.Store, sm, cols)
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

func (c *Connection) Update(model interface{}, excludeColumns ...string) error {
	return c.timeFunc("Update", func() error {
		sm := &Model{Value: model}

		cols := ColumnsForStruct(model, sm.TableName())
		cols.Remove("id", "created_at")
		cols.Remove(excludeColumns...)

		sm.touchUpdatedAt()

		return c.Dialect.Update(c.Store, sm, cols)
	})
}

func (c *Connection) Destroy(model interface{}) error {
	return c.timeFunc("Destroy", func() error {
		sm := &Model{Value: model}

		return c.Dialect.Destroy(c.Store, sm)
	})
}
