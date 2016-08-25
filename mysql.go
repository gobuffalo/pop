package pop

import (
	"fmt"
	"os/exec"

	_ "github.com/go-sql-driver/mysql"
	"github.com/markbates/going/clam"
	. "github.com/markbates/pop/columns"
	"github.com/markbates/pop/fizz"
	"github.com/markbates/pop/fizz/translators"
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
	return genericCreate(s, model, cols)
}

func (m *mysql) Update(s store, model *Model, cols Columns) error {
	return genericUpdate(s, model, cols)
}

func (m *mysql) Destroy(s store, model *Model) error {
	return genericDestroy(s, model)
}

func (m *mysql) SelectOne(s store, model *Model, query Query) error {
	return genericSelectOne(s, model, query)
}

func (m *mysql) SelectMany(s store, models *Model, query Query) error {
	return genericSelectMany(s, models, query)
}

func (m *mysql) CreateDB() error {
	c := m.ConnectionDetails
	cmd := exec.Command("mysql", "-u", c.User, "-p"+c.Password, "-e", fmt.Sprintf("create database %s", c.Database))
	return clam.RunAndListen(cmd, func(s string) {
		fmt.Println(s)
	})
}

func (m *mysql) DropDB() error {
	c := m.ConnectionDetails
	cmd := exec.Command("mysql", "-u", c.User, "-p"+c.Password, "-e", fmt.Sprintf("drop database %s", c.Database))
	return clam.RunAndListen(cmd, func(s string) {
		fmt.Println(s)
	})
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
	deets.Parse("3306")
	cd := &mysql{
		ConnectionDetails: deets,
	}

	return cd
}
