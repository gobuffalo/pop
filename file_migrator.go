package pop

import (
	"bytes"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"text/template"

	"github.com/gobuffalo/fizz"
	"github.com/pkg/errors"
)

// FileMigrator is a migrator for SQL and Fizz
// files on disk at a specified path.
type FileMigrator struct {
	Migrator
	Path string
}

// NewFileMigrator for a path and a Connection
func NewFileMigrator(path string, c *Connection) (FileMigrator, error) {
	fm := FileMigrator{
		Migrator: NewMigrator(c),
		Path:     path,
	}
	fm.SchemaPath = path

	err := fm.findMigrations()
	if err != nil {
		return fm, err
	}

	return fm, nil
}

func (fm *FileMigrator) findMigrations() error {
	dir := fm.Path
	if fi, err := os.Stat(dir); err != nil || !fi.IsDir() {
		// directory doesn't exist
		return nil
	}
	filepath.Walk(dir, func(p string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			match, err := ParseMigrationFilename(info.Name())
			if err != nil {
				return err
			}
			if match == nil {
				return nil
			}
			mf := Migration{
				Path:      p,
				Version:   match.Version,
				Name:      match.Name,
				DBType:    match.DBType,
				Direction: match.Direction,
				Type:      match.Type,
				Runner: func(mf Migration, tx *Connection) error {
					f, err := os.Open(p)
					if err != nil {
						return err
					}
					content, err := migrationContent(mf, tx, f)
					if err != nil {
						return errors.Wrapf(err, "error processing %s", mf.Path)
					}

					if content == "" {
						return nil
					}

					err = tx.RawQuery(content).Exec()
					if err != nil {
						return errors.Wrapf(err, "error executing %s, sql: %s", mf.Path, content)
					}
					return nil
				},
			}
			fm.Migrations[mf.Direction] = append(fm.Migrations[mf.Direction], mf)
		}
		return nil
	})
	return nil
}

func migrationContent(mf Migration, c *Connection, r io.Reader) (string, error) {
	b, err := ioutil.ReadAll(r)
	if err != nil {
		return "", nil
	}

	t := template.Must(template.New("sql").Parse(string(b)))
	var bb bytes.Buffer
	err = t.Execute(&bb, c.Dialect.Details())
	if err != nil {
		return "", errors.Wrapf(err, "could not execute migration template %s", mf.Path)
	}
	content := bb.String()

	if mf.Type == "fizz" {
		content, err = fizz.AString(content, c.Dialect.FizzTranslator())
		if err != nil {
			return "", errors.Wrapf(err, "could not fizz the migration %s", mf.Path)
		}
	}
	return content, nil
}
