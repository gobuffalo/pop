package pop

import (
	"bytes"
	"database/sql"
	"encoding/gob"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"

	"github.com/gobuffalo/pop/v6/columns"
	"github.com/gobuffalo/pop/v6/logging"
	"github.com/gofrs/uuid"
)

func init() {
	gob.Register(uuid.UUID{})
}

type commonDialect struct {
	ConnectionDetails *ConnectionDetails
}

func (commonDialect) Lock(fn func() error) error {
	return fn()
}

func (commonDialect) Quote(key string) string {
	parts := strings.Split(key, ".")

	for i, part := range parts {
		part = strings.Trim(part, `"`)
		part = strings.TrimSpace(part)

		parts[i] = fmt.Sprintf(`"%v"`, part)
	}

	return strings.Join(parts, ".")
}

func genericCreate(s store, model *Model, cols columns.Columns, quoter quotable) error {
	keyType, err := model.PrimaryKeyType()
	if err != nil {
		return err
	}
	switch keyType {
	case "int", "int64":
		var id int64
		cols.Remove(model.IDField())
		w := cols.Writeable()
		query := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)", quoter.Quote(model.TableName()), w.QuotedString(quoter), w.SymbolizedString())
		log(logging.SQL, query, model.Value)
		res, err := s.NamedExec(query, model.Value)
		if err != nil {
			return err
		}
		id, err = res.LastInsertId()
		if err == nil {
			model.setID(id)
		}
		if err != nil {
			return err
		}
		return nil
	case "UUID", "string":
		if keyType == "UUID" {
			if model.ID() == emptyUUID {
				u, err := uuid.NewV4()
				if err != nil {
					return err
				}
				model.setID(u)
			}
		} else if model.ID() == "" {
			return fmt.Errorf("missing ID value")
		}
		w := cols.Writeable()
		w.Add(model.IDField())
		query := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)", quoter.Quote(model.TableName()), w.QuotedString(quoter), w.SymbolizedString())
		log(logging.SQL, query, model.Value)
		stmt, err := s.PrepareNamed(query)
		if err != nil {
			return err
		}
		_, err = stmt.ExecContext(model.ctx, model.Value)
		if err != nil {
			if closeErr := stmt.Close(); closeErr != nil {
				return fmt.Errorf("failed to close prepared statement: %s: %w", closeErr, err)
			}
			return err
		}
		if err := stmt.Close(); err != nil {
			return fmt.Errorf("failed to close statement: %w", err)
		}
		return nil
	}
	return fmt.Errorf("can not use %s as a primary key type!", keyType)
}

func genericUpdate(s store, model *Model, cols columns.Columns, quoter quotable) error {
	stmt := fmt.Sprintf("UPDATE %s AS %s SET %s WHERE %s", quoter.Quote(model.TableName()), model.Alias(), cols.Writeable().QuotedUpdateString(quoter), model.WhereNamedID())
	log(logging.SQL, stmt, model.ID())
	_, err := s.NamedExec(stmt, model.Value)
	if err != nil {
		return err
	}
	return nil
}

func genericDestroy(s store, model *Model, quoter quotable) error {
	stmt := fmt.Sprintf("DELETE FROM %s AS %s WHERE %s", quoter.Quote(model.TableName()), model.Alias(), model.WhereID())
	_, err := genericExec(s, stmt, model.ID())
	if err != nil {
		return err
	}
	return nil
}

func genericDelete(s store, model *Model, query Query) error {
	sqlQuery, args := query.ToSQL(model)
	_, err := genericExec(s, sqlQuery, args...)
	return err
}

func genericExec(s store, stmt string, args ...interface{}) (sql.Result, error) {
	log(logging.SQL, stmt, args...)
	res, err := s.Exec(stmt, args...)
	return res, err
}

func genericSelectOne(s store, model *Model, query Query) error {
	sqlQuery, args := query.ToSQL(model)
	log(logging.SQL, sqlQuery, args...)
	err := s.Get(model.Value, sqlQuery, args...)
	if err != nil {
		return err
	}
	return nil
}

func genericSelectMany(s store, models *Model, query Query) error {
	sqlQuery, args := query.ToSQL(models)
	log(logging.SQL, sqlQuery, args...)
	err := s.Select(models.Value, sqlQuery, args...)
	if err != nil {
		return err
	}
	return nil
}

func genericLoadSchema(d dialect, r io.Reader) error {
	deets := d.Details()

	// Open DB connection on the target DB
	db, err := openPotentiallyInstrumentedConnection(d, d.MigrationURL())
	if err != nil {
		return fmt.Errorf("unable to load schema for %s: %w", deets.Database, err)
	}
	defer db.Close()

	// Get reader contents
	contents, err := ioutil.ReadAll(r)
	if err != nil {
		return err
	}

	if len(contents) == 0 {
		log(logging.Info, "schema is empty for %s, skipping", deets.Database)
		return nil
	}

	_, err = db.Exec(string(contents))
	if err != nil {
		return fmt.Errorf("unable to load schema for %s: %w", deets.Database, err)
	}

	log(logging.Info, "loaded schema for %s", deets.Database)
	return nil
}

func genericDumpSchema(deets *ConnectionDetails, cmd *exec.Cmd, w io.Writer) error {
	log(logging.SQL, strings.Join(cmd.Args, " "))

	bb := &bytes.Buffer{}
	mw := io.MultiWriter(w, bb)

	cmd.Stdout = mw
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		return err
	}

	x := bytes.TrimSpace(bb.Bytes())
	if len(x) == 0 {
		return fmt.Errorf("unable to dump schema for %s", deets.Database)
	}

	log(logging.Info, "dumped schema for %s", deets.Database)
	return nil
}
