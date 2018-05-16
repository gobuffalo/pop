package pop

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"sync"

	"github.com/jmoiron/sqlx"

	"github.com/gobuffalo/pop/columns"
	"github.com/gobuffalo/pop/fizz"
	"github.com/gobuffalo/pop/fizz/translators"
	"github.com/markbates/going/defaults"
	"github.com/pkg/errors"
)

var _ dialect = &postgresql{}

type postgresql struct {
	translateCache    map[string]string
	mu                sync.Mutex
	ConnectionDetails *ConnectionDetails
}

func (p *postgresql) Name() string {
	return "postgresql"
}

func (p *postgresql) Details() *ConnectionDetails {
	return p.ConnectionDetails
}

func (p *postgresql) Create(s store, model *Model, cols columns.Columns) error {
	keyType := model.PrimaryKeyType()
	switch keyType {
	case "int", "int64":
		cols.Remove("id")
		id := struct {
			ID int `db:"id"`
		}{}
		w := cols.Writeable()
		query := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s) returning id", model.TableName(), w.String(), w.SymbolizedString())
		Log(query)
		stmt, err := s.PrepareNamed(query)
		if err != nil {
			return errors.WithStack(err)
		}
		err = stmt.Get(&id, model.Value)
		if err != nil {
			return errors.WithStack(err)
		}
		model.setID(id.ID)
		return nil
	case "UUID":
		return genericCreate(s, model, cols)
	}
	return errors.Errorf("can not use %s as a primary key type!", keyType)
}

func (p *postgresql) Update(s store, model *Model, cols columns.Columns) error {
	return genericUpdate(s, model, cols)
}

func (p *postgresql) Destroy(s store, model *Model) error {
	return genericDestroy(s, model)
}

func (p *postgresql) SelectOne(s store, model *Model, query Query) error {
	return genericSelectOne(s, model, query)
}

func (p *postgresql) SelectMany(s store, models *Model, query Query) error {
	return genericSelectMany(s, models, query)
}

func (p *postgresql) CreateDB() error {
	// createdb -h db -p 5432 -U postgres enterprise_development
	deets := p.ConnectionDetails
	db, err := sqlx.Open(deets.Dialect, p.urlWithoutDb())
	if err != nil {
		return errors.Wrapf(err, "error creating PostgreSQL database %s", deets.Database)
	}
	defer db.Close()
	query := fmt.Sprintf("CREATE DATABASE \"%s\"", deets.Database)
	Log(query)

	_, err = db.Exec(query)
	if err != nil {
		return errors.Wrapf(err, "error creating PostgreSQL database %s", deets.Database)
	}

	fmt.Printf("created database %s\n", deets.Database)
	return nil
}

func (p *postgresql) DropDB() error {
	deets := p.ConnectionDetails
	db, err := sqlx.Open(deets.Dialect, p.urlWithoutDb())
	if err != nil {
		return errors.Wrapf(err, "error dropping PostgreSQL database %s", deets.Database)
	}
	defer db.Close()
	query := fmt.Sprintf("DROP DATABASE \"%s\"", deets.Database)
	Log(query)

	_, err = db.Exec(query)
	if err != nil {
		return errors.Wrapf(err, "error dropping PostgreSQL database %s", deets.Database)
	}

	fmt.Printf("dropped database %s\n", deets.Database)
	return nil
}

func (p *postgresql) URL() string {
	c := p.ConnectionDetails
	if c.URL != "" {
		return c.URL
	}
	ssl := defaults.String(c.Options["sslmode"], "disable")

	s := "postgres://%s:%s@%s:%s/%s?sslmode=%s"
	return fmt.Sprintf(s, c.User, c.Password, c.Host, c.Port, c.Database, ssl)
}

func (p *postgresql) urlWithoutDb() string {
	c := p.ConnectionDetails
	ssl := defaults.String(c.Options["sslmode"], "disable")

	// https://github.com/gobuffalo/buffalo/issues/836
	// If the db is not precised, postgresql takes the username as the database to connect on.
	// To avoid a connection problem if the user db is not here, we use the default "postgres"
	// db, just like the other client tools do.
	s := "postgres://%s:%s@%s:%s/postgres?sslmode=%s"
	return fmt.Sprintf(s, c.User, c.Password, c.Host, c.Port, ssl)
}

func (p *postgresql) MigrationURL() string {
	return p.URL()
}

func (p *postgresql) TranslateSQL(sql string) string {
	defer p.mu.Unlock()
	p.mu.Lock()

	if csql, ok := p.translateCache[sql]; ok {
		return csql
	}
	csql := sqlx.Rebind(sqlx.DOLLAR, sql)

	p.translateCache[sql] = csql
	return csql
}

func (p *postgresql) FizzTranslator() fizz.Translator {
	return translators.NewPostgres()
}

func (p *postgresql) Lock(fn func() error) error {
	return fn()
}

func (p *postgresql) DumpSchema(w io.Writer) error {
	cmd := exec.Command("pg_dump", "-s", fmt.Sprintf("--dbname=%s", p.URL()))
	Log(strings.Join(cmd.Args, " "))
	cmd.Stdout = w
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		return err
	}

	fmt.Printf("dumped schema for %s\n", p.Details().Database)
	return nil
}

func (p *postgresql) LoadSchema(r io.Reader) error {
	// Open DB connection on the target DB
	deets := p.ConnectionDetails
	db, err := sqlx.Open(deets.Dialect, p.MigrationURL())
	if err != nil {
		return errors.WithMessage(err, fmt.Sprintf("unable to load schema for %s", deets.Database))
	}
	defer db.Close()

	// Get reader contents
	contents, err := ioutil.ReadAll(r)
	if err != nil {
		return err
	}

	if len(contents) == 0 {
		fmt.Printf("schema is empty for %s, skipping\n", deets.Database)
		return nil
	}

	// From the sqlx package docs, this works with pq driver
	_, err = db.Exec(string(contents))
	if err != nil {
		return errors.WithMessage(err, fmt.Sprintf("unable to load schema for %s", deets.Database))
	}

	fmt.Printf("loaded schema for %s\n", deets.Database)
	return nil
}

func (p *postgresql) TruncateAll(tx *Connection) error {
	return tx.RawQuery(pgTruncate).Exec()
}

func newPostgreSQL(deets *ConnectionDetails) dialect {
	cd := &postgresql{
		ConnectionDetails: deets,
		translateCache:    map[string]string{},
		mu:                sync.Mutex{},
	}
	return cd
}

const pgTruncate = `DO
$func$
DECLARE
   _tbl text;
   _sch text;
BEGIN
   FOR _sch, _tbl IN
      SELECT schemaname, tablename
      FROM   pg_tables
      WHERE    schemaname = 'public'
   LOOP
      --RAISE ERROR '%',
      EXECUTE  -- dangerous, test before you execute!
         format('TRUNCATE TABLE %I.%I CASCADE', _sch, _tbl);
   END LOOP;
END
$func$;`
