package pop

import (
	"reflect"

	"github.com/gobuffalo/pop/associations"
	"github.com/gobuffalo/validate"
)

func (c *Connection) eagerCreate(model interface{}, excludeColumns ...string) error {
	asos, err := associations.AssociationsForStruct(model, c.eagerFields...)
	if err != nil {
		return err
	}

	c.eager = false
	for _, a := range asos {
		asoCreatable, ok := a.(associations.AssociationCreatable)
		if !ok {
			continue
		}

		// Create all dependencies first.
		dependencies := asoCreatable.CreatableDependencies()
		for _, d := range dependencies {
			if reflect.TypeOf(d) == reflect.TypeOf(model) {
				err = c.Create(d, excludeColumns...)
			} else {
				err = c.Create(d)
			}

			if err != nil {
				return err
			}
		}

		asoCreatable.Initialize()

		if acs, ok := a.(associations.AssociationCreatableStatement); ok {
			stms := acs.Statements()
			for _, stm := range stms {
				_, err = c.TX.Exec(c.Dialect.TranslateSQL(stm.Statement), stm.Args...)
				if err != nil {
					return err
				}
			}
			continue
		}

		i := a.Interface()
		if reflect.TypeOf(i) == reflect.TypeOf(model) {
			err = c.Create(i, excludeColumns...)
		} else {
			err = c.Create(i)
		}

		if err != nil {
			return err
		}
	}

	return err
}

func (c *Connection) eagerValidateAndCreate(model interface{}, excludeColumns ...string) (*validate.Errors, error) {
	asos, err := associations.AssociationsForStruct(model, c.eagerFields...)
	verrs := validate.NewErrors()

	if err != nil {
		return verrs, err
	}

	c.eager = false
	for _, a := range asos {
		asoCreatable, ok := a.(associations.AssociationCreatable)
		if !ok {
			continue
		}

		// Validate and create all dependencies first.
		dependencies := asoCreatable.CreatableDependencies()
		for _, d := range dependencies {
			verrs, err = c.ValidateAndCreate(d)
			if err != nil || verrs.HasAny() {
				return verrs, err
			}
		}

		verrs, err = c.ValidateAndCreate(a.Interface())
		if err != nil || verrs.HasAny() {
			return verrs, err
		}
	}

	return verrs, err
}
