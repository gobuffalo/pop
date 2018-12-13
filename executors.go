package pop

import (
	"reflect"

	"github.com/gobuffalo/pop/associations"
	"github.com/gobuffalo/pop/columns"
	"github.com/gobuffalo/pop/logging"
	"github.com/gobuffalo/uuid"
	"github.com/gobuffalo/validate"
)

// Reload fetch fresh data for a given model, using its ID.
func (c *Connection) Reload(model interface{}) error {
	sm := Model{Value: model}
	return sm.iterate(func(m *Model) error {
		return c.Find(m.Value, m.ID())
	})
}

// Exec runs the given query.
func (q *Query) Exec() error {
	return q.Connection.timeFunc("Exec", func() error {
		sql, args := q.ToSQL(nil)
		log(logging.SQL, sql, args...)
		_, err := q.Connection.Store.Exec(sql, args...)
		return err
	})
}

// ExecWithCount runs the given query, and returns the amount of
// affected rows.
func (q *Query) ExecWithCount() (int, error) {
	count := int64(0)
	return int(count), q.Connection.timeFunc("Exec", func() error {
		sql, args := q.ToSQL(nil)
		log(logging.SQL, sql, args...)
		result, err := q.Connection.Store.Exec(sql, args...)
		if err != nil {
			return err
		}

		count, err = result.RowsAffected()
		return err
	})
}

// ValidateAndSave applies validation rules on the given entry, then save it
// if the validation succeed, excluding the given columns.
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

// IsZeroOfUnderlyingType will check if the value of anything is the equal to the Zero value of that type.
func IsZeroOfUnderlyingType(x interface{}) bool {
	return reflect.DeepEqual(x, reflect.Zero(reflect.TypeOf(x)).Interface())
}

// Save wraps the Create and Update methods. It executes a Create if no ID is provided with the entry;
// or issues an Update otherwise.
func (c *Connection) Save(model interface{}, excludeColumns ...string) error {
	sm := &Model{Value: model}
	return sm.iterate(func(m *Model) error {
		id, err := m.fieldByName("ID")
		if err != nil {
			return err
		}
		if IsZeroOfUnderlyingType(id.Interface()) {
			return c.Create(m.Value, excludeColumns...)
		}
		return c.Update(m.Value, excludeColumns...)
	})
}

// ValidateAndCreate applies validation rules on the given entry, then creates it
// if the validation succeed, excluding the given columns.
func (c *Connection) ValidateAndCreate(model interface{}, excludeColumns ...string) (*validate.Errors, error) {
	sm := &Model{Value: model}
	verrs, err := sm.validateCreate(c)
	if err != nil {
		return verrs, err
	}
	if verrs.HasAny() {
		return verrs, nil
	}

	if c.eager {
		asos, err2 := associations.ForStruct(model, c.eagerFields...)

		if err2 != nil {
			return verrs, err2
		}

		if len(asos) == 0 {
			c.disableEager()
			return verrs, c.Create(model, excludeColumns...)
		}

		before := asos.AssociationsBeforeCreatable()
		for index := range before {
			i := before[index].BeforeInterface()
			if i == nil {
				continue
			}

			sm := &Model{Value: i}
			verrs, err := sm.validateAndOnlyCreate(c)
			if err != nil || verrs.HasAny() {
				return verrs, err
			}
		}

		after := asos.AssociationsAfterCreatable()
		for index := range after {
			i := after[index].AfterInterface()
			if i == nil {
				continue
			}

			sm := &Model{Value: i}
			verrs, err := sm.validateAndOnlyCreate(c)
			if err != nil || verrs.HasAny() {
				return verrs, err
			}
		}

		sm := &Model{Value: model}
		verrs, err = sm.validateCreate(c)
		if err != nil || verrs.HasAny() {
			return verrs, err
		}
	}

	return verrs, c.Create(model, excludeColumns...)
}

// Create add a new given entry to the database, excluding the given columns.
// It updates `created_at` and `updated_at` columns automatically.
func (c *Connection) Create(model interface{}, excludeColumns ...string) error {
	var isEager = c.eager

	c.disableEager()

	sm := &Model{Value: model}
	return sm.iterate(func(m *Model) error {
		return c.timeFunc("Create", func() error {
			var localIsEager = isEager
			asos, err := associations.ForStruct(m.Value, c.eagerFields...)
			if err != nil {
				return err
			}

			if localIsEager && len(asos) == 0 {
				// No association, fallback to non-eager mode.
				localIsEager = false
			}

			if err = m.beforeSave(c); err != nil {
				return err
			}

			if err = m.beforeCreate(c); err != nil {
				return err
			}

			processAssoc := len(asos) > 0

			if processAssoc {
				before := asos.AssociationsBeforeCreatable()
				for index := range before {
					i := before[index].BeforeInterface()
					if i == nil {
						continue
					}

					if localIsEager {
						sm := &Model{Value: i}
						err = sm.iterate(func(m *Model) error {
							id, err := m.fieldByName("ID")
							if err != nil {
								return err
							}
							if IsZeroOfUnderlyingType(id.Interface()) {
								return c.Create(m.Value)
							}
							return nil
						})

						if err != nil {
							return err
						}
					}

					err = before[index].BeforeSetup()
					if err != nil {
						return err
					}
				}
			}

			tn := m.TableName()
			cols := columns.ForStructWithAlias(m.Value, tn, m.As)

			if tn == sm.TableName() {
				cols.Remove(excludeColumns...)
			}

			m.touchCreatedAt()
			m.touchUpdatedAt()

			if err = c.Dialect.Create(c.Store, m, cols); err != nil {
				return err
			}

			if processAssoc {
				after := asos.AssociationsAfterCreatable()
				for index := range after {
					if localIsEager {
						err = after[index].AfterSetup()
						if err != nil {
							return err
						}

						i := after[index].AfterInterface()
						if i == nil {
							continue
						}

						sm := &Model{Value: i}
						err = sm.iterate(func(m *Model) error {
							fbn, err := m.fieldByName("ID")
							if err != nil {
								return err
							}
							id := fbn.Interface()
							if IsZeroOfUnderlyingType(id) {
								return c.Create(m.Value)
							} else {
								exists, errE := Q(c).Exists(i)
								if errE != nil || !exists {
									return c.Create(m.Value)
								}
							}
							return nil
						})

						if err != nil {
							return err
						}
					}
					stm := after[index].AfterProcess()
					if c.TX != nil && !stm.Empty() {
						_, err := c.TX.Exec(c.Dialect.TranslateSQL(stm.Statement), stm.Args...)
						if err != nil {
							return err
						}
					}
				}

				stms := asos.AssociationsCreatableStatement()
				for index := range stms {
					statements := stms[index].Statements()
					for _, stm := range statements {
						if c.TX != nil {
							_, err := c.TX.Exec(c.Dialect.TranslateSQL(stm.Statement), stm.Args...)
							if err != nil {
								return err
							}
							continue
						}
						_, err = c.Store.Exec(c.Dialect.TranslateSQL(stm.Statement), stm.Args...)
						if err != nil {
							return err
						}
					}
				}
			}

			if err = m.afterCreate(c); err != nil {
				return err
			}

			return m.afterSave(c)
		})
	})
}

// ValidateAndUpdate applies validation rules on the given entry, then update it
// if the validation succeed, excluding the given columns.
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

// Update writes changes from an entry to the database, excluding the given columns.
// It updates the `updated_at` column automatically.
func (c *Connection) Update(model interface{}, excludeColumns ...string) error {
	sm := &Model{Value: model}
	return sm.iterate(func(m *Model) error {
		return c.timeFunc("Update", func() error {
			var err error

			if err = m.beforeSave(c); err != nil {
				return err
			}
			if err = m.beforeUpdate(c); err != nil {
				return err
			}

			tn := m.TableName()
			cols := columns.ForStructWithAlias(model, tn, m.As)
			cols.Remove("id", "created_at")

			if tn == sm.TableName() {
				cols.Remove(excludeColumns...)
			}

			m.touchUpdatedAt()

			if err = c.Dialect.Update(c.Store, m, cols); err != nil {
				return err
			}
			if err = m.afterUpdate(c); err != nil {
				return err
			}

			return m.afterSave(c)
		})
	})
}

// Destroy deletes a given entry from the database
func (c *Connection) Destroy(model interface{}) error {
	sm := &Model{Value: model}
	return sm.iterate(func(m *Model) error {
		return c.timeFunc("Destroy", func() error {
			var err error

			if err = m.beforeDestroy(c); err != nil {
				return err
			}
			if err = c.Dialect.Destroy(c.Store, m); err != nil {
				return err
			}

			return m.afterDestroy(c)
		})
	})
}
