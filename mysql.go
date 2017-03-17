package pop

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"

	_ "github.com/go-sql-driver/mysql"
	"github.com/markbates/going/clam"
	. "github.com/markbates/pop/columns"
	"github.com/markbates/pop/fizz"
	"github.com/markbates/pop/fizz/translators"
	"github.com/pkg/errors"
)

var _ dialect = &mysql{}

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

func (m *mysql) DumpSchema(w io.Writer) error {
	// mysqldump -d -h localhost -P 3306 -u root --password=root coke_development
	deets := m.Details()
	cmd := exec.Command("mysqldump", "-d", "-h", deets.Host, "-P", deets.Port, "-u", deets.User, fmt.Sprintf("--password=%s", deets.Password), deets.Database)
	cmd.Stdout = w
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func (m *mysql) LoadSchema(r io.Reader) error {
	// mysql -u root --password=root -h localhost -P 3306 -D coke_test < schema.sql
	tmp, err := ioutil.TempFile("", "mysql-dump")
	if err != nil {
		return err
	}
	defer os.Remove(tmp.Name())
	_, err = io.Copy(tmp, r)
	if err != nil {
		return err
	}

	deets := m.Details()
	cmd := exec.Command("mysql", "-d", "-h", deets.Host, "-P", deets.Port, "-u", deets.User, fmt.Sprintf("--password=%s", deets.Password), "-D", deets.Database, "<", tmp.Name())
	return cmd.Run()
}

func newMySQL(deets *ConnectionDetails) dialect {
	cd := &mysql{
		ConnectionDetails: deets,
	}

	return cd
}
