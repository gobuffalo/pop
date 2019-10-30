package pop

import (
	"reflect"

	"github.com/gobuffalo/pop/associations"
	"github.com/gobuffalo/pop/columns"
	"github.com/gobuffalo/pop/logging"
	"github.com/gobuffalo/validate"
	"github.com/gofrs/uuid"
	"github.com/pkg/errors"
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
//
// If model is a slice, each item of the slice is validated then saved in the database.
func (c *Connection) ValidateAndSave(model interface{}, excludeColumns ...string) (*validate.Errors, error) {
	sm := &Model{Value: model}
	if err := sm.beforeValidate(c); err != nil {
		return nil, err
	}
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
//
// If model is a slice, each item of the slice is saved in the database.
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
//
// If model is a slice, each item of the slice is validated then created in the database.
func (c *Connection) ValidateAndCreate(model interface{}, excludeColumns ...string) (*validate.Errors, error) {
	sm := &Model{Value: model}
	if err := sm.beforeValidate(c); err != nil {
		return nil, err
	}
	verrs, err := sm.validateCreate(c)
	if err != nil {
		return verrs, err
	}
	if verrs.HasAny() {
		return verrs, nil
	}

	if c.eager {
		asos, err := associations.ForStruct(model, c.eagerFields...)
		if err != nil {
			return verrs, errors.Wrap(err, "could not retrieve associations")
		}

		if len(asos) == 0 {
			log(logging.Debug, "no associations found for given struct, disable eager mode")
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
//
// If model is a slice, each item of the slice is created in the database.
//
// Create support two modes:
// * Flat (default): Associate existing nested objects only. NO creation or update of nested objects.
// * Eager: Associate existing nested objects and create non-existent objects. NO change to existing objects.
func (c *Connection) Create(model interface{}, excludeColumns ...string) error {
	var isEager = c.eager

	c.disableEager()

	sm := &Model{Value: model}
	return sm.iterate(func(m *Model) error {
		return c.timeFunc("Create", func() error {
			var localIsEager = isEager
			asos, err := associations.ForStruct(m.Value, c.eagerFields...)
			if err != nil {
				return errors.Wrap(err, "could not retrieve associations")
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
							}
							exists, errE := Q(c).Exists(i)
							if errE != nil || !exists {
								return c.Create(m.Value)
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
//
// If model is a slice, each item of the slice is validated then updated in the database.
func (c *Connection) ValidateAndUpdate(model interface{}, excludeColumns ...string) (*validate.Errors, error) {
	sm := &Model{Value: model}
	if err := sm.beforeValidate(c); err != nil {
		return nil, err
	}
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
//
// If model is a slice, each item of the slice is updated in the database.
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

// UpdateColumns writes changes from an entry to the database, including only the given columns
// or all columns if no column names are provided.
// It updates the `updated_at` column automatically if you include `updated_at` in columnNames.
//
// If model is a slice, each item of the slice is updated in the database.
func (c *Connection) UpdateColumns(model interface{}, columnNames ...string) error {
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

			cols := columns.Columns{}
			if len(columnNames) > 0 && tn == sm.TableName() {
				cols = columns.NewColumnsWithAlias(tn, m.As)
				cols.Add(columnNames...)

			} else {
				cols = columns.ForStructWithAlias(model, tn, m.As)
			}
			cols.Remove("id", "created_at")

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

// Destroy deletes a given entry from the database.
//
// If model is a slice, each item of the slice is deleted from the database.
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
