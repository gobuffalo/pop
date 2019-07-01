package pop

import (
	"github.com/gobuffalo/packd"
	"github.com/pkg/errors"
)

// MigrationBox is a wrapper around packr.Box and Migrator.
// This will allow you to run migrations from a packed box
// inside of a compiled binary.
type MigrationBox struct {
	Migrator
	Box packd.Walkable
}

// NewMigrationBox from a packr.Box and a Connection.
func NewMigrationBox(box packd.Walkable, c *Connection) (MigrationBox, error) {
	fm := MigrationBox{
		Migrator: NewMigrator(c),
		Box:      box,
	}

	err := fm.findMigrations()
	if err != nil {
		return fm, err
	}

	return fm, nil
}

func (fm *MigrationBox) findMigrations() error {
	return fm.Box.Walk(func(p string, f packd.File) error {
		info, err := f.FileInfo()
		if err != nil {
			return err
		}
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
		return nil
	})
}
