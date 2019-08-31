package pop

import (
	"context"
	"fmt"
	"path/filepath"
	"time"

	"github.com/gobuffalo/genny"

	"github.com/gobuffalo/pop/internal/oncer"
)

// MigrationCreate writes contents for a given migration in normalized files
func MigrationCreate(path, name, ext string, up, down []byte) error {
	run := genny.WetRunner(context.Background())
	g := genny.New()

	n := time.Now().UTC()
	s := n.Format("20060102150405")

	upf := filepath.Join(path, fmt.Sprintf("%s_%s.up.%s", s, name, ext))
	g.File(genny.NewFileB(upf, up))

	downf := filepath.Join(path, fmt.Sprintf("%s_%s.down.%s", s, name, ext))
	g.File(genny.NewFileB(downf, down))

	run.With(g)

	return run.Run()
}

// MigrateUp is deprecated, and will be removed in a future version. Use FileMigrator#Up instead.
func (c *Connection) MigrateUp(path string) error {
	oncer.Deprecate(0, "pop.Connection#MigrateUp", "Use pop.FileMigrator#Up instead.")

	mig, err := NewFileMigrator(path, c)
	if err != nil {
		return err
	}
	return mig.Up()
}

// MigrateDown is deprecated, and will be removed in a future version. Use FileMigrator#Down instead.
func (c *Connection) MigrateDown(path string, step int) error {
	oncer.Deprecate(0, "pop.Connection#MigrateDown", "Use pop.FileMigrator#Down instead.")

	mig, err := NewFileMigrator(path, c)
	if err != nil {
		return err
	}
	return mig.Down(step)
}

// MigrateStatus is deprecated, and will be removed in a future version. Use FileMigrator#Status instead.
func (c *Connection) MigrateStatus(path string) error {
	oncer.Deprecate(0, "pop.Connection#MigrateStatus", "Use pop.FileMigrator#Status instead.")

	mig, err := NewFileMigrator(path, c)
	if err != nil {
		return err
	}
	return mig.Status()
}

// MigrateReset is deprecated, and will be removed in a future version. Use FileMigrator#Reset instead.
func (c *Connection) MigrateReset(path string) error {
	oncer.Deprecate(0, "pop.Connection#MigrateReset", "Use pop.FileMigrator#Reset instead.")

	mig, err := NewFileMigrator(path, c)
	if err != nil {
		return err
	}
	return mig.Reset()
}
