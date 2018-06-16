package pop

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"

	// Load MySQL Go driver
	_ "github.com/go-sql-driver/mysql"
	"github.com/gobuffalo/pop/columns"
	"github.com/gobuffalo/pop/fizz"
	"github.com/gobuffalo/pop/fizz/translators"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

var _ dialect = &mysql{}

type mysql struct {
	ConnectionDetails *ConnectionDetails
}

func (m *mysql) Name() string {
	return "mysql"
}

func (m *mysql) Details() *ConnectionDetails {
	return m.ConnectionDetails
}

func (m *mysql) URL() string {
	c := m.ConnectionDetails
	if m.ConnectionDetails.URL != "" {
		return strings.TrimPrefix(m.ConnectionDetails.URL, "mysql://")
	}
	s := "%s:%s@(%s:%s)/%s?parseTime=true&multiStatements=true&readTimeout=1s"
	return fmt.Sprintf(s, c.User, c.Password, c.Host, c.Port, c.Database)
}

func (m *mysql) urlWithoutDb() string {
	c := m.ConnectionDetails
	if m.ConnectionDetails.URL != "" {
		// respect user's own URL definition (with options).
		url := strings.TrimPrefix(m.ConnectionDetails.URL, "mysql://")
		return strings.Replace(url, "/"+c.Database+"?", "/?", 1)
	}
	s := "%s:%s@(%s:%s)/?parseTime=true&multiStatements=true&readTimeout=1s"
	return fmt.Sprintf(s, c.User, c.Password, c.Host, c.Port)
}

func (m *mysql) MigrationURL() string {
	return m.URL()
}

func (m *mysql) Create(s store, model *Model, cols columns.Columns) error {
	return errors.Wrap(genericCreate(s, model, cols), "mysql create")
}

func (m *mysql) Update(s store, model *Model, cols columns.Columns) error {
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

// CreateDB creates a new database, from the given connection credentials
func (m *mysql) CreateDB() error {
	deets := m.ConnectionDetails
	db, err := sqlx.Open(deets.Dialect, m.urlWithoutDb())
	if err != nil {
		return errors.Wrapf(err, "error creating MySQL database %s", deets.Database)
	}
	defer db.Close()
	query := fmt.Sprintf("CREATE DATABASE `%s` DEFAULT COLLATE `utf8_general_ci`", deets.Database)
	Log(query)

	_, err = db.Exec(query)
	if err != nil {
		return errors.Wrapf(err, "error creating MySQL database %s", deets.Database)
	}

	fmt.Printf("created database %s\n", deets.Database)
	return nil
}

// DropDB drops an existing database, from the given connection credentials
func (m *mysql) DropDB() error {
	deets := m.ConnectionDetails
	db, err := sqlx.Open(deets.Dialect, m.urlWithoutDb())
	if err != nil {
		return errors.Wrapf(err, "error dropping MySQL database %s", deets.Database)
	}
	defer db.Close()
	query := fmt.Sprintf("DROP DATABASE `%s`", deets.Database)
	Log(query)

	_, err = db.Exec(query)
	if err != nil {
		return errors.Wrapf(err, "error dropping MySQL database %s", deets.Database)
	}

	fmt.Printf("dropped database %s\n", deets.Database)
	return nil
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
	deets := m.Details()
	cmd := exec.Command("mysqldump", "-d", "-h", deets.Host, "-P", deets.Port, "-u", deets.User, fmt.Sprintf("--password=%s", deets.Password), deets.Database)
	if deets.Port == "socket" {
		cmd = exec.Command("mysqldump", "-d", "-S", deets.Host, "-u", deets.User, fmt.Sprintf("--password=%s", deets.Password), deets.Database)
	}
	Log(strings.Join(cmd.Args, " "))
	cmd.Stdout = w
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		return err
	}

	fmt.Printf("dumped schema for %s\n", m.Details().Database)
	return nil
}

func (m *mysql) LoadSchema(r io.Reader) error {
	deets := m.Details()
	cmd := exec.Command("mysql", "-u", deets.User, fmt.Sprintf("--password=%s", deets.Password), "-h", deets.Host, "-P", deets.Port, "-D", deets.Database)
	if deets.Port == "socket" {
		cmd = exec.Command("mysql", "-u", deets.User, fmt.Sprintf("--password=%s", deets.Password), "-S", deets.Host, "-D", deets.Database)
	}
	in, err := cmd.StdinPipe()
	if err != nil {
		return err
	}
	go func() {
		defer in.Close()
		io.Copy(in, r)
	}()
	Log(strings.Join(cmd.Args, " "))
	err = cmd.Start()
	if err != nil {
		return err
	}

	err = cmd.Wait()
	if err != nil {
		return err
	}

	fmt.Printf("loaded schema for %s\n", m.Details().Database)
	return nil
}

func (m *mysql) TruncateAll(tx *Connection) error {
	stmts := []string{}
	err := tx.RawQuery(mysqlTruncate, m.Details().Database).All(&stmts)
	if err != nil {
		return err
	}
	if len(stmts) == 0 {
		return nil
	}

	var qb bytes.Buffer
	// #49: Disable foreign keys before truncation
	qb.WriteString("SET SESSION FOREIGN_KEY_CHECKS = 0; ")
	qb.WriteString(strings.Join(stmts, " "))
	// #49: Re-enable foreign keys after truncation
	qb.WriteString(" SET SESSION FOREIGN_KEY_CHECKS = 1;")

	return tx.RawQuery(qb.String()).Exec()
}

func newMySQL(deets *ConnectionDetails) dialect {
	cd := &mysql{
		ConnectionDetails: deets,
	}

	return cd
}

const mysqlTruncate = "SELECT concat('TRUNCATE TABLE `', TABLE_NAME, '`;') as stmt FROM INFORMATION_SCHEMA.TABLES WHERE table_schema = ? AND table_type <> 'VIEW'"
