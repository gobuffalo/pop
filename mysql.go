package pop

import (
	"fmt"
	"os/exec"

	_ "github.com/go-sql-driver/mysql"
	"github.com/markbates/going/clam"
	. "github.com/markbates/pop/columns"
	"github.com/markbates/pop/fizz"
	"github.com/markbates/pop/fizz/translators"
	"github.com/pkg/errors"
)

type mysql struct {
	ConnectionDetails *ConnectionDetails
}

func (m *mysql) Details() *ConnectionDetails {
	return m.ConnectionDetails
}

func (m *mysql) URL() string {
	c := m.ConnectionDetails
	if c.URL != "" {
		return c.URL
	}

	s := "%s:%s@(%s:%s)/%s?parseTime=true&multiStatements=true&readTimeout=1s"
	return fmt.Sprintf(s, c.User, c.Password, c.Host, c.Port, c.Database)
}

func (m *mysql) MigrationURL() string {
	return m.URL()
}

func (m *mysql) Create(s store, model *Model, cols Columns) error {
	return errors.Wrap(genericCreate(s, model, cols), "mysql create")
}

func (m *mysql) Update(s store, model *Model, cols Columns) error {
	return errors.Wrap(genericUpdate(s, model, cols), "mysql update")
}

func (m *mysql) Destroy(s store, model *Model) error {
	return errors.Wrap(genericDestroy(s, model), "mysql destroy")
}

func (m *mysql) SelectOne(s store, model *Model, query Query) error {
	return errors.Wrap(genericSelectOne(s, model, query), "mysql select one")
}

func (m *mysql) SelectMany(s store, models *Model, query Query) error {
	return errors.Wrap(genericSelectMany(s, models, query), "mysql select many")
}

func (m *mysql) CreateDB() error {
	c := m.ConnectionDetails
	cmd := exec.Command("mysql", "-u", c.User, "-p"+c.Password, "-h", c.Host, "-P", c.Port, "-e", fmt.Sprintf("create database %s", c.Database))
	err := clam.RunAndListen(cmd, func(s string) {
		fmt.Println(s)
	})
	return errors.Wrapf(err, "error creating MySQL database %s", c.Database)
}

func (m *mysql) DropDB() error {
	c := m.ConnectionDetails
	cmd := exec.Command("mysql", "-u", c.User, "-p"+c.Password, "-h", c.Host, "-P", c.Port, "-e", fmt.Sprintf("drop database %s", c.Database))
	err := clam.RunAndListen(cmd, func(s string) {
		fmt.Println(s)
	})
	return errors.Wrapf(err, "error dropping MySQL database %s", c.Database)
}

func (m *mysql) TranslateSQL(sql string) string {
	return sql
}

func (m *mysql) FizzTranslator() fizz.Translator {
	t := translators.NewMySQL(m.URL(), m.Details().Database)
	return t
}

func (m *mysql) Lock(fn func() error) error {
	return fn()
}

func newMySQL(deets *ConnectionDetails) dialect {
	cd := &mysql{
		ConnectionDetails: deets,
	}

	return cd
}
