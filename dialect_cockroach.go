package pop

import (
	"database/sql"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path"
	"strings"
	"sync"

	_ "github.com/cockroachdb/cockroach-go/crdb" // Load CockroachdbQL/postgres Go driver which also loads github.com/lib/pq
	"github.com/gobuffalo/fizz"
	"github.com/gobuffalo/fizz/translators"
	"github.com/gobuffalo/pop/columns"
	"github.com/gobuffalo/pop/logging"
	"github.com/jmoiron/sqlx"
	"github.com/markbates/going/defaults"
	"github.com/pkg/errors"
)

const nameCockroach = "cockroach"
const portCockroach = "26257"

func init() {
	AvailableDialects = append(AvailableDialects, nameCockroach)
	dialectSynonyms["cockroachdb"] = nameCockroach
	dialectSynonyms["crdb"] = nameCockroach
	finalizer[nameCockroach] = finalizerCockroach
	newConnection[nameCockroach] = newCockroach
}

var _ dialect = &cockroach{}

// ServerInfo holds informational data about connected database server.
type cockroachInfo struct {
	VersionString string `db:"version"`
	product       string `db:"-"`
	license       string `db:"-"`
	version       string `db:"-"`
	buildInfo     string `db:"-"`
	client        string `db:"-"`
}

type cockroach struct {
	translateCache    map[string]string
	mu                sync.Mutex
	ConnectionDetails *ConnectionDetails
	info              cockroachInfo
}

func (p *cockroach) Name() string {
	return nameCockroach
}

func (p *cockroach) Details() *ConnectionDetails {
	return p.ConnectionDetails
}

func (p *cockroach) Create(s store, model *Model, cols columns.Columns) error {
	keyType := model.PrimaryKeyType()
	switch keyType {
	case "int", "int64":
		cols.Remove("id")
		id := struct {
			ID int `db:"id"`
		}{}
		w := cols.Writeable()
		var query string
		if len(w.Cols) > 0 {
			query = fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s) returning id", model.TableName(), w.String(), w.SymbolizedString())
		} else {
			query = fmt.Sprintf("INSERT INTO %s DEFAULT VALUES returning id", model.TableName())
		}
		log(logging.SQL, query)
		stmt, err := s.PrepareNamed(query)
		if err != nil {
			return errors.WithStack(err)
		}
		err = stmt.Get(&id, model.Value)
		if err != nil {
			if err := stmt.Close(); err != nil {
				return errors.WithMessage(err, "failed to close statement")
			}
			return errors.WithStack(err)
		}
		model.setID(id.ID)
		return errors.WithMessage(stmt.Close(), "failed to close statement")
	}
	return genericCreate(s, model, cols)
}

func (p *cockroach) Update(s store, model *Model, cols columns.Columns) error {
	return genericUpdate(s, model, cols)
}

func (p *cockroach) Destroy(s store, model *Model) error {
	stmt := p.TranslateSQL(fmt.Sprintf("DELETE FROM %s WHERE %s", model.TableName(), model.whereID()))
	err := genericExec(s, stmt, model.ID())
	if err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func (p *cockroach) SelectOne(s store, model *Model, query Query) error {
	return genericSelectOne(s, model, query)
}

func (p *cockroach) SelectMany(s store, models *Model, query Query) error {
	return genericSelectMany(s, models, query)
}

func (p *cockroach) CreateDB() error {
	// createdb -h db -p 5432 -U cockroach enterprise_development
	deets := p.ConnectionDetails
	db, err := sql.Open(deets.Dialect, p.urlWithoutDb())
	if err != nil {
		return errors.Wrapf(err, "error creating Cockroach database %s", deets.Database)
	}
	defer db.Close()
	query := fmt.Sprintf("CREATE DATABASE \"%s\"", deets.Database)
	log(logging.SQL, query)

	_, err = db.Exec(query)
	if err != nil {
		return errors.Wrapf(err, "error creating Cockroach database %s", deets.Database)
	}

	log(logging.Info, "created database %s", deets.Database)
	return nil
}

func (p *cockroach) DropDB() error {
	deets := p.ConnectionDetails
	db, err := sql.Open(deets.Dialect, p.urlWithoutDb())
	if err != nil {
		return errors.Wrapf(err, "error dropping Cockroach database %s", deets.Database)
	}
	defer db.Close()
	query := fmt.Sprintf("DROP DATABASE \"%s\" CASCADE;", deets.Database)
	log(logging.SQL, query)

	_, err = db.Exec(query)
	if err != nil {
		return errors.Wrapf(err, "error dropping Cockroach database %s", deets.Database)
	}

	log(logging.Info, "dropped database %s", deets.Database)
	return nil
}

func (p *cockroach) URL() string {
	c := p.ConnectionDetails
	if c.URL != "" {
		return c.URL
	}
	s := "postgres://%s:%s@%s:%s/%s?%s"
	return fmt.Sprintf(s, c.User, c.Password, c.Host, c.Port, c.Database, c.OptionsString(""))
}

func (p *cockroach) urlWithoutDb() string {
	c := p.ConnectionDetails
	s := "postgres://%s:%s@%s:%s/?%s"
	return fmt.Sprintf(s, c.User, c.Password, c.Host, c.Port, c.OptionsString(""))
}

func (p *cockroach) MigrationURL() string {
	return p.URL()
}

func (p *cockroach) TranslateSQL(sql string) string {
	defer p.mu.Unlock()
	p.mu.Lock()

	if csql, ok := p.translateCache[sql]; ok {
		return csql
	}
	csql := sqlx.Rebind(sqlx.DOLLAR, sql)

	p.translateCache[sql] = csql
	return csql
}

func (p *cockroach) FizzTranslator() fizz.Translator {
	return translators.NewCockroach(p.URL(), p.Details().Database)
}

func (p *cockroach) Lock(fn func() error) error {
	return fn()
}

func (p *cockroach) DumpSchema(w io.Writer) error {
	cmd := exec.Command("cockroach", "dump", p.Details().Database, "--dump-mode=schema")

	c := p.ConnectionDetails
	if defaults.String(c.Options["sslmode"], "disable") == "disable" || strings.Contains(c.RawOptions, "sslmode=disable") {
		cmd.Args = append(cmd.Args, "--insecure")
	}
	return genericDumpSchema(p.Details(), cmd, w)
}

func (p *cockroach) LoadSchema(r io.Reader) error {
	return genericLoadSchema(p.ConnectionDetails, p.MigrationURL(), r)
}

func (p *cockroach) TruncateAll(tx *Connection) error {
	type table struct {
		TableName string `db:"table_name"`
	}

	tableQuery := "select table_name from information_schema.tables where table_schema = 'public' and table_type = 'BASE TABLE' and table_catalog = ?"
	if strings.HasPrefix(p.info.version, "v1") {
		tableQuery = "select table_name from information_schema.tables where table_schema = ?"
	}

	var tables []table
	if err := tx.RawQuery(tableQuery, tx.Dialect.Details().Database).All(&tables); err != nil {
		return err
	}

	if len(tables) == 0 {
		return nil
	}

	tableNames := make([]string, len(tables))
	for i, t := range tables {
		tableNames[i] = t.TableName
		//! work around for current limitation of DDL and DML at the same transaction.
		//  it should be fixed when cockroach support it or with other approach.
		//  https://www.cockroachlabs.com/docs/stable/known-limitations.html#schema-changes-within-transactions
		if err := tx.RawQuery(fmt.Sprintf("delete from %s", t.TableName)).Exec(); err != nil {
			return err
		}
	}
	return nil
	// TODO!
	// return tx3.RawQuery(fmt.Sprintf("truncate %s cascade;", strings.Join(tableNames, ", "))).Exec()
}

func (p *cockroach) afterOpen(c *Connection) error {
	if err := c.RawQuery(`select version() AS "version"`).First(&p.info); err != nil {
		return err
	}
	if s := strings.Split(p.info.VersionString, " "); len(s) > 3 {
		p.info.product = s[0]
		p.info.license = s[1]
		p.info.version = s[2]
		p.info.buildInfo = s[3]
	}
	log(logging.Debug, "server: %v %v %v", p.info.product, p.info.license, p.info.version)

	return nil
}

func newCockroach(deets *ConnectionDetails) (dialect, error) {
	deets.Dialect = "postgres"
	d := &cockroach{
		ConnectionDetails: deets,
		translateCache:    map[string]string{},
		mu:                sync.Mutex{},
	}
	d.info.client = deets.Options["application_name"]
	return d, nil
}

func finalizerCockroach(cd *ConnectionDetails) {
	appName := path.Base(os.Args[0])
	cd.Options["application_name"] = defaults.String(cd.Options["application_name"], appName)
	cd.Port = defaults.String(cd.Port, portCockroach)
}
