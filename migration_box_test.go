package pop

import (
	"os"
	"testing"

	"github.com/gobuffalo/pop/v6/logging"
	"github.com/stretchr/testify/require"
)

type logEntry struct {
	lvl  logging.Level
	s    string
	args []interface{}
}

func setNewTestLogger() *[]logEntry {
	var logs []logEntry
	log = func(lvl logging.Level, s string, args ...interface{}) {
		logs = append(logs, logEntry{lvl, s, args})
	}
	return &logs
}

func Test_MigrationBox(t *testing.T) {
	if PDB == nil {
		t.Skip("skipping integration tests")
	}

	t.Run("finds testdata", func(t *testing.T) {
		r := require.New(t)

		b, err := NewMigrationBox(os.DirFS("testdata/migrations/multiple"), PDB)
		r.NoError(err)
		r.Equal(4, len(b.UpMigrations.Migrations))
		r.Equal("mysql", b.UpMigrations.Migrations[0].DBType)
		r.Equal("postgres", b.UpMigrations.Migrations[1].DBType)
		r.Equal("sqlite3", b.UpMigrations.Migrations[2].DBType)
		r.Equal("all", b.UpMigrations.Migrations[3].DBType)
	})

	t.Run("ignores clutter files", func(t *testing.T) {
		logs := setNewTestLogger()
		r := require.New(t)

		b, err := NewMigrationBox(os.DirFS("testdata/migrations/cluttered"), PDB)
		r.NoError(err)
		r.Equal(1, len(b.UpMigrations.Migrations))
		r.Equal(1, len(*logs))
		r.Equal(logging.Warn, (*logs)[0].lvl)
		r.Contains((*logs)[0].s, "ignoring file")
		r.Equal([]interface{}{"clutter.txt"}, (*logs)[0].args)
	})

	t.Run("ignores unsupported files", func(t *testing.T) {
		logs := setNewTestLogger()
		r := require.New(t)

		b, err := NewMigrationBox(os.DirFS("testdata/migrations/unsupported_dialect"), PDB)
		r.NoError(err)
		r.Equal(0, len(b.UpMigrations.Migrations))
		r.Equal(1, len(*logs))
		r.Equal(logging.Warn, (*logs)[0].lvl)
		r.Contains((*logs)[0].s, "ignoring migration")
		r.Equal([]interface{}{"unsupported dialect unsupported"}, (*logs)[0].args)
	})
}
