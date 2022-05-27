//go:build sqlite
// +build sqlite

package pop

import (
	"fmt"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

var sqliteDefaultOptions = map[string]string{"_busy_timeout": "5000", "_fk": "true"}

func Test_ConnectionDetails_Finalize_SQLite_URL_Only(t *testing.T) {
	r := require.New(t)

	cd := &ConnectionDetails{
		URL: "sqlite3:///tmp/foo.db",
	}
	err := cd.Finalize() // calls withURL() and urlParserSQLite3(cd)
	r.NoError(err)
	r.Equal("sqlite3", cd.Dialect, "given dialect: N/A")
	r.Equal("/tmp/foo.db", cd.Database, "given url: sqlite3:///tmp/foo.db")
	r.EqualValues(sqliteDefaultOptions, cd.Options, "given url: sqlite3:///tmp/foo.db")
}

func Test_ConnectionDetails_Finalize_SQLite_OverrideOptions_URL_Only(t *testing.T) {
	r := require.New(t)

	cd := &ConnectionDetails{
		URL: "sqlite3:///tmp/foo.db?_fk=false&foo=bar",
	}
	err := cd.Finalize() // calls withURL() and urlParserSQLite3(cd)
	r.NoError(err)
	r.Equal("sqlite3", cd.Dialect, "given dialect: N/A")
	r.Equal("/tmp/foo.db", cd.Database, "given url: sqlite3:///tmp/foo.db?_fk=false&foo=bar")
	r.EqualValues(map[string]string{"_fk": "false", "foo": "bar", "_busy_timeout": "5000"}, cd.Options, "given url: sqlite3:///tmp/foo.db?_fk=false&foo=bar")
}

func Test_ConnectionDetails_Finalize_SQLite_SynURL_Only(t *testing.T) {
	r := require.New(t)

	cd := &ConnectionDetails{
		URL: "sqlite:///tmp/foo.db",
	}
	err := cd.Finalize() // calls withURL() and urlParserSQLite3(cd)
	r.NoError(err)
	r.Equal("sqlite3", cd.Dialect, "given dialect: N/A")
	r.Equal("/tmp/foo.db", cd.Database, "given url: sqlite:///tmp/foo.db")
	r.EqualValues(sqliteDefaultOptions, cd.Options, "given url: sqlite3:///tmp/foo.db")
}

func Test_ConnectionDetails_Finalize_SQLite_Dialect_URL(t *testing.T) {
	r := require.New(t)
	cd := &ConnectionDetails{
		Dialect: "sqlite3",
		URL:     "sqlite3:///tmp/foo.db",
	}
	err := cd.Finalize() // calls withURL() and urlParserSQLite3(cd)
	r.NoError(err)
	r.Equal("sqlite3", cd.Dialect, "given dialect: sqlite3")
	r.Equal("/tmp/foo.db", cd.Database, "given url: sqlite3:///tmp/foo.db")
	r.EqualValues(sqliteDefaultOptions, cd.Options, "given url: sqlite3:///tmp/foo.db")
}

func Test_ConnectionDetails_Finalize_SQLite_Dialect_SynURL(t *testing.T) {
	r := require.New(t)
	cd := &ConnectionDetails{
		Dialect: "sqlite3",
		URL:     "sqlite:///tmp/foo.db",
	}
	err := cd.Finalize() // calls withURL() and urlParserSQLite3(cd)
	r.NoError(err)
	r.Equal("sqlite3", cd.Dialect, "given dialect: sqlite3")
	r.Equal("/tmp/foo.db", cd.Database, "given url: sqlite:///tmp/foo.db")
	r.EqualValues(sqliteDefaultOptions, cd.Options, "given url: sqlite3:///tmp/foo.db")
}

func Test_ConnectionDetails_Finalize_SQLite_Synonym_URL(t *testing.T) {
	r := require.New(t)
	cd := &ConnectionDetails{
		Dialect: "sqlite",
		URL:     "sqlite3:///tmp/foo.db",
	}
	err := cd.Finalize() // calls withURL() and urlParserSQLite3(cd)
	r.NoError(err)
	r.Equal("sqlite3", cd.Dialect, "given dialect: sqlite")
	r.Equal("/tmp/foo.db", cd.Database, "given url: sqlite3:///tmp/foo.db")
	r.Equal(sqliteDefaultOptions, cd.Options, "given url: sqlite3:///tmp/foo.db")
}

func Test_ConnectionDetails_Finalize_SQLite_Synonym_SynURL(t *testing.T) {
	r := require.New(t)
	cd := &ConnectionDetails{
		Dialect: "sqlite",
		URL:     "sqlite:///tmp/foo.db",
	}
	err := cd.Finalize() // calls withURL() and urlParserSQLite3(cd)
	r.NoError(err)
	r.Equal("sqlite3", cd.Dialect, "given dialect: sqlite")
	r.Equal("/tmp/foo.db", cd.Database, "given url: sqlite:///tmp/foo.db")
	r.EqualValues(sqliteDefaultOptions, cd.Options, "given url: sqlite:///tmp/foo.db")
}

func Test_ConnectionDetails_Finalize_SQLite_Synonym_Path(t *testing.T) {
	r := require.New(t)
	cd := &ConnectionDetails{
		Dialect:  "sqlite",
		Database: "./foo.db",
	}
	err := cd.Finalize() // calls withURL() and urlParserSQLite3(cd)
	r.NoError(err)
	r.Equal("sqlite3", cd.Dialect, "given dialect: sqlite")
	r.Equal("./foo.db", cd.Database, "given database: ./foo.db")
	r.EqualValues(sqliteDefaultOptions, cd.Options, "given url: ./foo.db")
}

func Test_ConnectionDetails_Finalize_SQLite_OverrideOptions_Synonym_Path(t *testing.T) {
	r := require.New(t)

	cd := &ConnectionDetails{
		URL: "sqlite3:///tmp/foo.db?_fk=false&foo=bar",
	}
	err := cd.Finalize() // calls withURL() and urlParserSQLite3(cd)
	r.NoError(err)
	r.Equal("sqlite3", cd.Dialect, "given dialect: N/A")
	r.Equal("/tmp/foo.db", cd.Database, "given url: sqlite3:///tmp/foo.db")
	r.EqualValues(map[string]string{"_fk": "false", "foo": "bar", "_busy_timeout": "5000"}, cd.Options, "given url: sqlite3:///tmp/foo.db?_fk=false&foo=bar")
}

func Test_ConnectionDetails_FinalizeOSPath(t *testing.T) {
	r := require.New(t)
	d := t.TempDir()
	p := filepath.Join(d, "testdb.sqlite")
	cd := &ConnectionDetails{
		Dialect:  "sqlite",
		Database: p,
	}
	r.NoError(cd.Finalize())
	r.Equal("sqlite3", cd.Dialect)
	r.EqualValues(p, cd.Database)
}

func TestSqlite_CreateDB(t *testing.T) {
	r := require.New(t)

	cd := &ConnectionDetails{Dialect: "sqlite"}
	dialect, err := newSQLite(cd)
	r.NoError(err)

	t.Run("CreateFile", func(t *testing.T) {
		dir := t.TempDir()
		cd.Database = filepath.Join(dir, "testdb.sqlite")

		r.NoError(dialect.CreateDB())
		r.FileExists(cd.Database)
	})

	t.Run("MemoryDB_tag", func(t *testing.T) {
		dir := t.TempDir()
		cd.Database = filepath.Join(dir, "file::memory:?cache=shared")

		r.NoError(dialect.CreateDB())
		r.NoFileExists(cd.Database)
	})

	t.Run("MemoryDB_only", func(t *testing.T) {
		dir := t.TempDir()
		cd.Database = filepath.Join(dir, ":memory:")

		r.NoError(dialect.CreateDB())
		r.NoFileExists(cd.Database)
	})

	t.Run("MemoryDB_param", func(t *testing.T) {
		dir := t.TempDir()
		cd.Database = filepath.Join(dir, "file:foobar?mode=memory&cache=shared")

		r.NoError(dialect.CreateDB())
		r.NoFileExists(cd.Database)
	})

	t.Run("CreateFile_ExistingDB", func(t *testing.T) {
		dir := t.TempDir()
		cd.Database = filepath.Join(dir, "testdb.sqlite")

		r.NoError(dialect.CreateDB())
		r.EqualError(dialect.CreateDB(), fmt.Sprintf("could not create SQLite database '%s'; database exists", cd.Database))
	})

}

func TestSqlite_NewDriver(t *testing.T) {
	_, err := newSQLiteDriver()
	require.NoError(t, err)
}
