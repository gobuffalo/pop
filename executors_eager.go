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

func (c *Connection) eagerValidateAndCreate(model interface{}, excludeColumns ...string) (*validate.Errors, error) {
	asos, err := associations.ForStruct(model, c.eagerFields...)
	verrs := validate.NewErrors()

	if err != nil {
		return verrs, err
	}

	if len(asos) == 0 {
		c.disableEager()
		return c.ValidateAndCreate(model, excludeColumns...)
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

	return verrs, c.eagerCreate(model, excludeColumns...)
}
