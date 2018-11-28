package pop

import (
	"github.com/gobuffalo/pop/associations"
	"github.com/gobuffalo/validate"
)

func (c *Connection) eagerCreate(model interface{}, excludeColumns ...string) error {
	assos, err := associations.ForStruct(model, c.eagerFields...)
	if err != nil {
		return err
	}

	c.disableEager()

	// No association, fallback to non-eager mode.
	if len(assos) == 0 {
		return c.Create(model, excludeColumns...)
	}

	// Try to create the associations the root model depends on.
	before := assos.AssociationsBeforeCreatable()
	for index := range before {
		i := before[index].BeforeInterface()
		if i == nil {
			continue
		}

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

		err = before[index].BeforeSetup()
		if err != nil {
			return err
		}
	}

	// Create the root model
	err = c.Create(model, excludeColumns...)
	if err != nil {
		return err
	}

	// Try to create the associations depending on the root model.
	after := assos.AssociationsAfterCreatable()
	for index := range after {
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
			return nil
		})

		if err != nil {
			return err
		}
	}

	stms := assos.AssociationsCreatableStatement()
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

	return err
}

func (c *Connection) eagerUpdate(model interface{}, excludeColumns ...string) error {
	assos, err := associations.ForStruct(model, c.eagerFields...)
	if err != nil {
		return err
	}

	c.disableEager()

	// No association, fallback to non-eager mode.
	if len(assos) == 0 {
		return c.Update(model, excludeColumns...)
	}
	// Try to update/create the associations the root model depends on.
	before := assos.AssociationsBeforeUpdatable()
	for index := range before {
		i := before[index].BeforeInterface()
		if i == nil {
			continue
		}

		sm := &Model{Value: i}
		err = sm.iterate(func(m *Model) error {
			id, err := m.fieldByName("ID")
			if err != nil {
				return err
			}

			if IsZeroOfUnderlyingType(id.Interface()) {
				return c.Create(m.Value)
			} else {
				return c.Update(m.Value)
			}
			return nil
		})

		if err != nil {
			return err
		}

		err = before[index].BeforeSetup()
		if err != nil {
			return err
		}

	}

	// Update the root model
	err = c.Update(model, excludeColumns...)
	if err != nil {
		return err
	}

	//Try to update/create the associations depending on the model.
	after := assos.AssociationsAfterUpdatable()
	for index := range after {

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
				err = c.Create(m.Value)

				if err != nil {
					return err
				}
			} else {
				err = c.Update(m.Value)
				if err != nil {
					return err
				}
			}

			return nil
		})

		 stm := after[index].AfterFixRelationships()

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

		if err != nil {
			return err
		}
	}

	// TODO  Should I fire Process relationships? or something like that?
	// TODO Need to clean up the has one when I update a has one then I need to

	stms := assos.AssociationsCreatableStatement()
	for index := range stms {
		statements := stms[index].Statements()

		// Create Associations
		// TODO need to check the existing associations for deletions
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
		//	Delete Associations.µ
	}

	return err

}

func (c *Connection) eagerValidateAndCreate(model interface{}, excludeColumns ...string) (*validate.Errors, error) {
	asoss, err := associations.ForStruct(model, c.eagerFields...)
	verrs := validate.NewErrors()

	if err != nil {
		return verrs, err
	}

	if len(asoss) == 0 {
		c.disableEager()
		return c.ValidateAndCreate(model, excludeColumns...)
	}

	before := asoss.AssociationsBeforeCreatable()
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

	after := asoss.AssociationsAfterCreatable()
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

	return verrs, c.eagerCreate(model, excludeColumns...)
}

//Eager Update

func (c *Connection) eagerValidateAndUpdate(model interface{}, excludeColumns ...string) (*validate.Errors, error) {
	// Get all the associations for the struct

	asos, err := associations.ForStruct(model, c.eagerFields...)
	verrs := validate.NewErrors()

	if err != nil {
		return verrs, err
	}

	if len(asos) == 0 {
		c.disableEager()
		return c.ValidateAndCreate(model, excludeColumns...)
	}

	before := asos.AssociationsBeforeUpdatable()
	for index := range before {
		i := before[index].BeforeInterface()
		if i == nil {
			continue
		}

		sm := &Model{Value: i}
		verrs, err := sm.validateAndOnlyUpdate(c)
		if err != nil || verrs.HasAny() {
			return verrs, err
		}
	}

	after := asos.AssociationsAfterUpdatable()
	for index := range after {
		i := after[index].AfterInterface()
		if i == nil {
			continue
		}

		sm := &Model{Value: i}
		verrs, err := sm.validateAndOnlyUpdate(c)
		if err != nil || verrs.HasAny() {
			return verrs, err
		}
	}

	sm := &Model{Value: model}
	verrs, err = sm.validateUpdate(c)
	if err != nil || verrs.HasAny() {
		return verrs, err
	}

	return verrs, c.eagerUpdate(model, excludeColumns...)
}
