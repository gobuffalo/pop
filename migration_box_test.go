package pop

import (
	"testing"

	"github.com/gobuffalo/packr/v2"
	"github.com/gobuffalo/pop/v5/logging"
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

		b, err := NewMigrationBox(packr.New("./testdata/migrations/multiple", "./testdata/migrations/multiple"), PDB)
		r.NoError(err)
		r.Equal(4, len(b.Migrations["up"]))
		r.Equal("mysql", b.Migrations["up"][0].DBType)
		r.Equal("postgres", b.Migrations["up"][1].DBType)
		r.Equal("sqlite3", b.Migrations["up"][2].DBType)
		r.Equal("all", b.Migrations["up"][3].DBType)
	})

	t.Run("ignores clutter files", func(t *testing.T) {
		logs := setNewTestLogger()
		r := require.New(t)

		b, err := NewMigrationBox(packr.New("./testdata/migrations/cluttered", "./testdata/migrations/cluttered"), PDB)
		r.NoError(err)
		r.Equal(1, len(b.Migrations["up"]))
		r.Equal(1, len(*logs))
		r.Equal(logging.Warn, (*logs)[0].lvl)
		r.Contains((*logs)[0].s, "ignoring file")
		r.Equal([]interface{}{"clutter.txt"}, (*logs)[0].args)
	})

	t.Run("ignores unsupported files", func(t *testing.T) {
		logs := setNewTestLogger()
		r := require.New(t)

		b, err := NewMigrationBox(packr.New("./testdata/migrations/unsupported_dialect", "./testdata/migrations/unsupported_dialect"), PDB)
		r.NoError(err)
		r.Equal(0, len(b.Migrations["up"]))
		r.Equal(1, len(*logs))
		r.Equal(logging.Warn, (*logs)[0].lvl)
		r.Contains((*logs)[0].s, "ignoring migration")
		r.Equal([]interface{}{"unsupported dialect unsupported"}, (*logs)[0].args)
	})
}
