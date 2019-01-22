package pop

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_ConnectionDetails_Finalize_SQLite_URL_Only(t *testing.T) {
	r := require.New(t)

	cd := &ConnectionDetails{
		URL: "sqlite3:///tmp/foo.db",
	}
	err := cd.Finalize() // calls withURL() and urlParserSQLite3(cd)
	r.NoError(err)
	r.Equal("sqlite3", cd.Dialect, "given dialect: N/A")
	r.Equal("/tmp/foo.db", cd.Database, "given url: sqlite3:///tmp/foo.db")
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
}

func Test_ConnectionDetails_FinalizeOSPath(t *testing.T) {
	r := require.New(t)
	d, err := ioutil.TempDir("", "")
	r.NoError(err)
	p := filepath.Join(d, "testdb.sqlite")
	defer os.RemoveAll(p)
	cd := &ConnectionDetails{
		Dialect:  "sqlite",
		Database: p,
	}
	r.NoError(cd.Finalize())
	r.Equal("sqlite3", cd.Dialect)
	r.Equal(p, cd.Database)
}
