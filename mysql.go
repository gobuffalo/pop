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

type MySQL struct {
	ConnectionDetails *ConnectionDetails
}

func (m *MySQL) Details() *ConnectionDetails {
	return m.ConnectionDetails
}

func (m *MySQL) URL() string {
	c := m.ConnectionDetails
	if c.URL != "" {
		return c.URL
	}

	s := "%s:%s@(%s:%s)/%s?parseTime=true&multiStatements=true"
	return fmt.Sprintf(s, c.User, c.Password, c.Host, c.Port, c.Database)
}

func (m *MySQL) MigrationURL() string {
	return m.URL()
}

func (m *MySQL) Create(store Store, model *Model, cols Columns) error {
	return genericCreate(store, model, cols)
}

func (m *MySQL) Update(store Store, model *Model, cols Columns) error {
	return genericUpdate(store, model, cols)
}

func (m *MySQL) Destroy(store Store, model *Model) error {
	return genericDestroy(store, model)
}

func (m *MySQL) SelectOne(store Store, model *Model, query Query) error {
	return genericSelectOne(store, model, query)
}

func (m *MySQL) SelectMany(store Store, models *Model, query Query) error {
	return genericSelectMany(store, models, query)
}

func (m *MySQL) CreateDB() error {
	c := m.ConnectionDetails
	cmd := exec.Command("mysql", "-u "+c.User, "--password"+c.Password, "-e", fmt.Sprintf("create database %s", c.Database))
	fmt.Printf("### cmd -> %#v\n", cmd)
	return clam.RunAndListen(cmd, func(s string) {
		fmt.Println(s)
	})
}

func (m *MySQL) DropDB() error {
	c := m.ConnectionDetails
	cmd := exec.Command("mysql", "-u"+c.User, "-p"+c.Password, "-e", fmt.Sprintf("drop database %s", c.Database))
	return clam.RunAndListen(cmd, func(s string) {
		fmt.Println(s)
	})
}

func (m *MySQL) TranslateSQL(sql string) string {
	return sql
}

func (m *MySQL) FizzTranslator() fizz.Translator {
	t := translators.NewMySQL(m.URL(), m.Details().Database)
	return t
}

func NewMySQL(deets *ConnectionDetails) Dialect {
	deets.Parse("3306")
	cd := &MySQL{
		ConnectionDetails: deets,
	}

	return cd
}
