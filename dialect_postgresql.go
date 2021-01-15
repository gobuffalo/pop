package pop

import (
	"fmt"
	"io"
	"os/exec"
	"sync"

	// Load pgx driver
	_ "github.com/jackc/pgx/v4/stdlib"

	"github.com/gobuffalo/fizz"
	"github.com/gobuffalo/fizz/translators"
	"github.com/jackc/pgconn"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"

	"github.com/gobuffalo/pop/v5/columns"
	"github.com/gobuffalo/pop/v5/internal/defaults"
	"github.com/gobuffalo/pop/v5/logging"
)

const namePostgreSQL = "postgres"
const portPostgreSQL = "5432"

func init() {
	AvailableDialects = append(AvailableDialects, namePostgreSQL)
	dialectSynonyms["postgresql"] = namePostgreSQL
	dialectSynonyms["pg"] = namePostgreSQL
	dialectSynonyms["pgx"] = namePostgreSQL
	urlParser[namePostgreSQL] = urlParserPostgreSQL
	finalizer[namePostgreSQL] = finalizerPostgreSQL
	newConnection[namePostgreSQL] = newPostgreSQL
}

var _ dialect = &postgresql{}

type postgresql struct {
	commonDialect
	translateCache map[string]string
	mu             sync.Mutex
}

func (p *postgresql) Name() string {
	return namePostgreSQL
}

func (p *postgresql) DefaultDriver() string {
	return "pgx"
}

func (p *postgresql) Details() *ConnectionDetails {
	return p.ConnectionDetails
}

func (p *postgresql) Create(s store, model *Model, cols columns.Columns) error {
	keyType, err := model.PrimaryKeyType()
	if err != nil {
		return err
	}
	switch keyType {
	case "int", "int64":
		cols.Remove(model.IDField())
		w := cols.Writeable()
		var query string
		if len(w.Cols) > 0 {
			query = fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s) returning %s", p.Quote(model.TableName()), w.QuotedString(p), w.SymbolizedString(), model.IDField())
		} else {
			query = fmt.Sprintf("INSERT INTO %s DEFAULT VALUES returning %s", p.Quote(model.TableName()), model.IDField())
		}
		log(logging.SQL, query)
		stmt, err := s.PrepareNamed(query)
		if err != nil {
			return err
		}
		id := map[string]interface{}{}
		err = stmt.QueryRow(model.Value).MapScan(id)
		if err != nil {
			if closeErr := stmt.Close(); closeErr != nil {
				return errors.Wrapf(err, "failed to close prepared statement: %s", closeErr)
			}
			return err
		}
		model.setID(id[model.IDField()])
		return errors.WithMessage(stmt.Close(), "failed to close statement")
	}
	return genericCreate(s, model, cols, p)
}

func (p *postgresql) Update(s store, model *Model, cols columns.Columns) error {
	return genericUpdate(s, model, cols, p)
}

func (p *postgresql) Destroy(s store, model *Model) error {
	stmt := p.TranslateSQL(fmt.Sprintf("DELETE FROM %s WHERE %s", p.Quote(model.TableName()), model.whereID()))
	_, err := genericExec(s, stmt, model.ID())
	if err != nil {
		return err
	}
	return nil
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

	db, err := openPotentiallyInstrumentedConnection(p, p.urlWithoutDb())
	if err != nil {
		return errors.Wrapf(err, "error creating PostgreSQL database %s", deets.Database)
	}
	defer db.Close()
	query := fmt.Sprintf("CREATE DATABASE %s", p.Quote(deets.Database))
	log(logging.SQL, query)

	_, err = db.Exec(query)
	if err != nil {
		return errors.Wrapf(err, "error creating PostgreSQL database %s", deets.Database)
	}

	log(logging.Info, "created database %s", deets.Database)
	return nil
}

func (p *postgresql) DropDB() error {
	deets := p.ConnectionDetails

	db, err := openPotentiallyInstrumentedConnection(p, p.urlWithoutDb())
	if err != nil {
		return errors.Wrapf(err, "error dropping PostgreSQL database %s", deets.Database)
	}
	defer db.Close()
	query := fmt.Sprintf("DROP DATABASE %s", p.Quote(deets.Database))
	log(logging.SQL, query)

	_, err = db.Exec(query)
	if err != nil {
		return errors.Wrapf(err, "error dropping PostgreSQL database %s", deets.Database)
	}

	log(logging.Info, "dropped database %s", deets.Database)
	return nil
}

func (p *postgresql) URL() string {
	c := p.ConnectionDetails
	if c.URL != "" {
		return c.URL
	}
	s := "postgres://%s:%s@%s:%s/%s?%s"
	return fmt.Sprintf(s, c.User, c.Password, c.Host, c.Port, c.Database, c.OptionsString(""))
}

func (p *postgresql) urlWithoutDb() string {
	c := p.ConnectionDetails
	// https://github.com/gobuffalo/buffalo/issues/836
	// If the db is not precised, postgresql takes the username as the database to connect on.
	// To avoid a connection problem if the user db is not here, we use the default "postgres"
	// db, just like the other client tools do.
	s := "postgres://%s:%s@%s:%s/postgres?%s"
	return fmt.Sprintf(s, c.User, c.Password, c.Host, c.Port, c.OptionsString(""))
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

func (p *postgresql) DumpSchema(w io.Writer) error {
	cmd := exec.Command("pg_dump", "-s", fmt.Sprintf("--dbname=%s", p.URL()))
	return genericDumpSchema(p.Details(), cmd, w)
}

// LoadSchema executes a schema sql file against the configured database.
func (p *postgresql) LoadSchema(r io.Reader) error {
	return genericLoadSchema(p, r)
}

// TruncateAll truncates all tables for the given connection.
func (p *postgresql) TruncateAll(tx *Connection) error {
	return tx.RawQuery(fmt.Sprintf(pgTruncate, tx.MigrationTableName())).Exec()
}

func newPostgreSQL(deets *ConnectionDetails) (dialect, error) {
	cd := &postgresql{
		commonDialect:  commonDialect{ConnectionDetails: deets},
		translateCache: map[string]string{},
		mu:             sync.Mutex{},
	}
	return cd, nil
}

// urlParserPostgreSQL parses the options the same way jackc/pgconn does:
// https://pkg.go.dev/github.com/jackc/pgconn?tab=doc#ParseConfig
// After parsed, they are set to ConnectionDetails instance
func urlParserPostgreSQL(cd *ConnectionDetails) error {
	conf, err := pgconn.ParseConfig(cd.URL)
	if err != nil {
		return err
	}

	cd.Database = conf.Database
	cd.Host = conf.Host
	cd.User = conf.User
	cd.Password = conf.Password
	cd.Port = fmt.Sprintf("%d", conf.Port)

	options := []string{"fallback_application_name"}
	for i := range options {
		if opt, ok := conf.RuntimeParams[options[i]]; ok {
			cd.Options[options[i]] = opt
		}
	}

	if conf.TLSConfig == nil {
		cd.Options["sslmode"] = "disable"
	}

	return nil
}

func finalizerPostgreSQL(cd *ConnectionDetails) {
	cd.Port = defaults.String(cd.Port, portPostgreSQL)
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
      WHERE  tablename <> '%s' AND schemaname NOT IN ('pg_catalog', 'information_schema') AND tableowner = current_user
   LOOP
      --RAISE ERROR '%%',
      EXECUTE  -- dangerous, test before you execute!
         format('TRUNCATE TABLE %%I.%%I CASCADE', _sch, _tbl);
   END LOOP;
END
$func$;`
