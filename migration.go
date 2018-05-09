package pop

import (
	"fmt"
	"path/filepath"
	"runtime"
	"time"

	"github.com/gobuffalo/makr"
	"github.com/gobuffalo/pop/log"
	"github.com/pkg/errors"
)

// MigrationCreate writes contents for a given migration in normalized files
func MigrationCreate(path, name, ext string, up, down []byte) error {
	g := makr.New()
	n := time.Now().UTC()
	s := n.Format("20060102150405")

	upf := filepath.Join(path, (fmt.Sprintf("%s_%s.up.%s", s, name, ext)))
	g.Add(makr.NewFile(upf, string(up)))

	downf := filepath.Join(path, (fmt.Sprintf("%s_%s.down.%s", s, name, ext)))
	g.Add(makr.NewFile(downf, string(down)))

	return g.Run(".", makr.Data{})
}

// MigrateUp is deprecated, and will be removed in a future version. Use FileMigrator#Up instead.
func (c *Connection) MigrateUp(path string) error {
	logctx := log.DefaultLogger
	_, file, no, ok := runtime.Caller(1)
	if ok {
		logctx = logctx.WithField("file", file).WithField("line", no)
	}
	logctx.Warn("Connection#MigrateUp is deprecated, and will be removed in a future version. Use FileMigrator#Up instead.")

	mig, err := NewFileMigrator(path, c)
	if err != nil {
		return errors.WithStack(err)
	}
	return mig.Up()
}

// MigrateDown is deprecated, and will be removed in a future version. Use FileMigrator#Down instead.
func (c *Connection) MigrateDown(path string, step int) error {
	logctx := log.DefaultLogger
	_, file, no, ok := runtime.Caller(1)
	if ok {
		logctx = logctx.WithField("file", file).WithField("line", no)
	}
	logctx.Warn("Connection#MigrateDown is deprecated, and will be removed in a future version. Use FileMigrator#Down instead.")

	mig, err := NewFileMigrator(path, c)
	if err != nil {
		return errors.WithStack(err)
	}
	return mig.Down(step)
}

// MigrateStatus is deprecated, and will be removed in a future version. Use FileMigrator#Status instead.
func (c *Connection) MigrateStatus(path string) error {
	logctx := log.DefaultLogger
	_, file, no, ok := runtime.Caller(1)
	if ok {
		logctx = logctx.WithField("file", file).WithField("line", no)
	}
	logctx.Warn("Connection#MigrateStatus is deprecated, and will be removed in a future version. Use FileMigrator#Status instead.")

	mig, err := NewFileMigrator(path, c)
	if err != nil {
		return errors.WithStack(err)
	}
	return mig.Status()
}

// MigrateReset is deprecated, and will be removed in a future version. Use FileMigrator#Reset instead.
func (c *Connection) MigrateReset(path string) error {
	logctx := log.DefaultLogger
	_, file, no, ok := runtime.Caller(1)
	if ok {
		logctx = logctx.WithField("file", file).WithField("line", no)
	}
	logctx.Warn("Connection#MigrateReset is deprecated, and will be removed in a future version. Use FileMigrator#Reset instead.")

	mig, err := NewFileMigrator(path, c)
	if err != nil {
		return errors.WithStack(err)
	}
	return mig.Reset()
}
