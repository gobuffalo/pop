package pop

import (
	"os"

	"github.com/markbates/pkger"
	"github.com/pkg/errors"
)

// PkgerMigrator is a migrator for SQL and Fizz
// files packaged by markbates/pkger
type PkgerMigrator struct {
	Migrator
	Path string
}

// NewPkgerMigrator for a path and a Connection
func NewPkgerMigrator(path string, c *Connection) (PkgerMigrator, error) {
	fm := PkgerMigrator{
		Migrator: NewMigrator(c),
		Path:     path,
	}
	fm.SchemaPath = path

	runner := func(mf Migration, tx *Connection) error {
		f, err := pkger.Open(mf.Path)
		if err != nil {
			return err
		}
		defer f.Close()
		content, err := MigrationContent(mf, tx, f, true)
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
	}

	err := fm.findMigrations(runner)
	if err != nil {
		return fm, err
	}

	return fm, nil
}

func (fm *PkgerMigrator) findMigrations(runner func(mf Migration, tx *Connection) error) error {
	dir := fm.Path
	if fi, err := pkger.Stat(dir); err != nil || !fi.IsDir() {
		// directory doesn't exist
		return nil
	}
	return pkger.Walk(dir, func(p string, info os.FileInfo, err error) error {
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
				Runner:    runner,
			}
			fm.Migrations[mf.Direction] = append(fm.Migrations[mf.Direction], mf)
		}
		return nil
	})
}
