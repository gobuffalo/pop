package pop

func (c *Connection) Reload(model interface{}) error {
	sm := Model{Value: model}
	return c.Find(model, sm.ID())
}

func (q *Query) Exec() error {
	sql, args := q.ToSQL(nil)
	_, err := q.Connection.Store.Exec(sql, args...)
	return err
}

func (c *Connection) Save(model interface{}, excludeColumns ...string) error {
	sm := &Model{Value: model}
	if sm.ID() == 0 {
		return c.Create(model, excludeColumns...)
	} else {
		return c.Update(model, excludeColumns...)
	}
}

func (c *Connection) Create(model interface{}, excludeColumns ...string) error {
	sm := &Model{Value: model}

	cols := ColumnsForStruct(model)
	cols.Remove(excludeColumns...)

	sm.TouchCreatedAt()
	sm.TouchUpdatedAt()

	return c.Dialect.Create(c.Store, sm, cols)
}

func (c *Connection) Update(model interface{}, excludeColumns ...string) error {
	sm := &Model{Value: model}

	cols := ColumnsForStruct(model)
	cols.Remove("id", "created_at")
	cols.Remove(excludeColumns...)

	sm.TouchUpdatedAt()

	return c.Dialect.Update(c.Store, sm, cols)
}

func (c *Connection) Destroy(model interface{}) error {
	sm := &Model{Value: model}

	return c.Dialect.Destroy(c.Store, sm)
}
