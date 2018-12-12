package pop

import (
	"database/sql"
	"fmt"
	"io"
	"os/exec"
	"sync"

	"github.com/gobuffalo/fizz"
	"github.com/gobuffalo/fizz/translators"
	"github.com/gobuffalo/pop/columns"
	"github.com/gobuffalo/pop/logging"
	"github.com/jmoiron/sqlx"
	"github.com/markbates/going/defaults"
	"github.com/pkg/errors"
)

const namePostgreSQL = "postgres"
const portPostgreSQL = "5432"

func init() {
	AvailableDialects = append(AvailableDialects, namePostgreSQL)
	dialectSynonyms["postgresql"] = namePostgreSQL
	dialectSynonyms["pg"] = namePostgreSQL
	finalizer[namePostgreSQL] = finalizerPostgreSQL
	newConnection[namePostgreSQL] = newPostgreSQL
}

var _ dialect = &postgresql{}

type postgresql struct {
	translateCache    map[string]string
	mu                sync.Mutex
	ConnectionDetails *ConnectionDetails
}

func (p *postgresql) Name() string {
	return namePostgreSQL
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

func (p *postgresql) Update(s store, model *Model, cols columns.Columns) error {
	return genericUpdate(s, model, cols)
}

func (p *postgresql) Destroy(s store, model *Model) error {
	stmt := p.TranslateSQL(fmt.Sprintf("DELETE FROM %s WHERE %s", model.TableName(), model.whereID()))
	err := genericExec(s, stmt, model.ID())
	if err != nil {
		return errors.WithStack(err)
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
	db, err := sql.Open(deets.Dialect, p.urlWithoutDb())
	if err != nil {
		return errors.Wrapf(err, "error creating PostgreSQL database %s", deets.Database)
	}
	defer db.Close()
	query := fmt.Sprintf("CREATE DATABASE \"%s\"", deets.Database)
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
	db, err := sql.Open(deets.Dialect, p.urlWithoutDb())
	if err != nil {
		return errors.Wrapf(err, "error dropping PostgreSQL database %s", deets.Database)
	}
	defer db.Close()
	query := fmt.Sprintf("DROP DATABASE \"%s\"", deets.Database)
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

func (p *postgresql) Lock(fn func() error) error {
	return fn()
}

func (p *postgresql) DumpSchema(w io.Writer) error {
	cmd := exec.Command("pg_dump", "-s", fmt.Sprintf("--dbname=%s", p.URL()))
	return genericDumpSchema(p.Details(), cmd, w)
}

// LoadSchema executes a schema sql file against the configured database.
func (p *postgresql) LoadSchema(r io.Reader) error {
	return genericLoadSchema(p.ConnectionDetails, p.MigrationURL(), r)
}

// TruncateAll truncates all tables for the given connection.
func (p *postgresql) TruncateAll(tx *Connection) error {
	return tx.RawQuery(fmt.Sprintf(pgTruncate, tx.MigrationTableName())).Exec()
}

func (p *postgresql) afterOpen(c *Connection) error {
	return nil
}

func newPostgreSQL(deets *ConnectionDetails) (dialect, error) {
	cd := &postgresql{
		ConnectionDetails: deets,
		translateCache:    map[string]string{},
		mu:                sync.Mutex{},
	}
	return cd, nil
}

func finalizerPostgreSQL(cd *ConnectionDetails) {
	cd.Options["sslmode"] = defaults.String(cd.Options["sslmode"], "disable")
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
