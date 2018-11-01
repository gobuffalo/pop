package pop

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_ConnectionDetails_Finalize(t *testing.T) {
	r := require.New(t)

	cd := &ConnectionDetails{
		URL: "postgres://user:pass@host:port/database",
	}
	err := cd.Finalize()
	r.NoError(err)

	r.Equal("database", cd.Database)
	r.Equal("postgres", cd.Dialect)
	r.Equal("host", cd.Host)
	r.Equal("pass", cd.Password)
	r.Equal("port", cd.Port)
	r.Equal("user", cd.User)
}

func Test_ConnectionDetails_Finalize_MySQL_DSN(t *testing.T) {
	r := require.New(t)

	url := "mysql://user:pass@(host:port)/database?param1=value1&param2=value2"
	cd := &ConnectionDetails{
		URL: url,
	}
	err := cd.Finalize()
	r.NoError(err)

	r.Equal(url, cd.URL)
	r.Equal("mysql", cd.Dialect)
	r.Equal("user", cd.User)
	r.Equal("pass", cd.Password)
	r.Equal("host", cd.Host)
	r.Equal("port", cd.Port)
	r.Equal("database", cd.Database)
}

func Test_ConnectionDetails_Finalize_Cockroach(t *testing.T) {
	r := require.New(t)
	cd := &ConnectionDetails{
		Dialect: "cockroach",
		URL:     "postgres://user:pass@host:port/database?sslmode=require&sslrootcert=certs/ca.crt&sslkey=certs/client.key&sslcert=certs/client.crt",
	}
	err := cd.Finalize()
	r.NoError(err)
	r.Equal("cockroach", cd.Dialect)
	r.Equal("database", cd.Database)
	r.Equal("host", cd.Host)
	r.Equal("port", cd.Port)
	r.Equal("user", cd.User)
	r.Equal("pass", cd.Password)
}

func Test_ConnectionDetails_Finalize_MySQL_DSN_collation(t *testing.T) {
	r := require.New(t)

	urls := []string{
		"mysql://user:pass@(host:port)/database?collation=utf8mb4_general_ci",
		"mysql://user:pass@(host:port)/database?collation=utf8mb4_general_ci&readTimeout=10s",
		"mysql://user:pass@(host:port)/database?readTimeout=10s&collation=utf8mb4_general_ci",
	}

	for _, url := range urls {
		cd := &ConnectionDetails{
			URL: url,
		}
		err := cd.Finalize()
		r.NoError(err)

		r.Equal(url, cd.URL)
		r.Equal("mysql", cd.Dialect)
		r.Equal("user", cd.User)
		r.Equal("pass", cd.Password)
		r.Equal("host", cd.Host)
		r.Equal("port", cd.Port)
		r.Equal("database", cd.Database)
		r.Equal("utf8mb4_general_ci", cd.Options["collation"])
	}
}

func Test_ConnectionDetails_Finalize_MySQL_DSN_Protocol(t *testing.T) {
	r := require.New(t)

	url := "mysql://user:pass@tcp(host:port)/protocol"
	cd := &ConnectionDetails{
		URL: url,
	}
	err := cd.Finalize()
	r.NoError(err)

	r.Equal(url, cd.URL)
	r.Equal("mysql", cd.Dialect)
	r.Equal("user", cd.User)
	r.Equal("pass", cd.Password)
	r.Equal("host", cd.Host)
	r.Equal("port", cd.Port)
	r.Equal("protocol", cd.Database)
}

func Test_ConnectionDetails_Finalize_MySQL_DSN_Socket(t *testing.T) {
	r := require.New(t)

	url := "mysql://user:pass@unix(/path/to/socket)/socket"
	cd := &ConnectionDetails{
		URL: url,
	}
	err := cd.Finalize()
	r.NoError(err)

	r.Equal(url, cd.URL)
	r.Equal("mysql", cd.Dialect)
	r.Equal("user", cd.User)
	r.Equal("pass", cd.Password)
	r.Equal("/path/to/socket", cd.Host)
	r.Equal("socket", cd.Port)
	r.Equal("socket", cd.Database)
}

func Test_ConnectionDetails_Finalize_UnknownDialect(t *testing.T) {
	r := require.New(t)
	cd := &ConnectionDetails{
		URL: "unknown://user:pass@host:port/database",
	}
	err := cd.Finalize()
	r.Error(err)
}

func Test_ConnectionDetails_Finalize_SQLite_URL_Only(t *testing.T) {
	r := require.New(t)

	cd := &ConnectionDetails{
		URL: "sqlite3:///tmp/foo.db",
	}
	err := cd.Finalize()
	r.NoError(err)
	r.Equal("sqlite3", cd.Dialect, "given dialect: N/A")
	r.Equal("/tmp/foo.db", cd.Database, "given url: sqlite3:///tmp/foo.db")
	r.Equal("", cd.Host)
	r.Equal("", cd.Password)
	r.Equal("", cd.Port)
	r.Equal("", cd.User)
}

func Test_ConnectionDetails_Finalize_SQLite_SynURL_Only(t *testing.T) {
	r := require.New(t)

	cd := &ConnectionDetails{
		URL: "sqlite:///tmp/foo.db",
	}
	err := cd.Finalize()
	r.NoError(err)
	r.Equal("sqlite3", cd.Dialect, "given dialect: N/A")
	r.Equal("/tmp/foo.db", cd.Database, "given url: sqlite:///tmp/foo.db")
	r.Equal("", cd.Host)
	r.Equal("", cd.Password)
	r.Equal("", cd.Port)
	r.Equal("", cd.User)
}

func Test_ConnectionDetails_Finalize_SQLite_Dialect_URL(t *testing.T) {
	r := require.New(t)
	cd := &ConnectionDetails{
		Dialect: "sqlite3",
		URL:     "sqlite3:///tmp/foo.db",
	}
	err := cd.Finalize()
	r.NoError(err)
	r.Equal("sqlite3", cd.Dialect, "given dialect: sqlite3")
	r.Equal("/tmp/foo.db", cd.Database, "given url: sqlite3:///tmp/foo.db")
	r.Equal("", cd.Host)
	r.Equal("", cd.Password)
	r.Equal("", cd.Port)
	r.Equal("", cd.User)
}

func Test_ConnectionDetails_Finalize_SQLite_Dialect_SynURL(t *testing.T) {
	r := require.New(t)
	cd := &ConnectionDetails{
		Dialect: "sqlite3",
		URL:     "sqlite:///tmp/foo.db",
	}
	err := cd.Finalize()
	r.NoError(err)
	r.Equal("sqlite3", cd.Dialect, "given dialect: sqlite3")
	r.Equal("/tmp/foo.db", cd.Database, "given url: sqlite:///tmp/foo.db")
	r.Equal("", cd.Host)
	r.Equal("", cd.Password)
	r.Equal("", cd.Port)
	r.Equal("", cd.User)
}

func Test_ConnectionDetails_Finalize_SQLite_Synonym_URL(t *testing.T) {
	r := require.New(t)
	cd := &ConnectionDetails{
		Dialect: "sqlite",
		URL:     "sqlite3:///tmp/foo.db",
	}
	err := cd.Finalize()
	r.NoError(err)
	r.Equal("sqlite3", cd.Dialect, "given dialect: sqlite")
	r.Equal("/tmp/foo.db", cd.Database, "given url: sqlite3:///tmp/foo.db")
	r.Equal("", cd.Host)
	r.Equal("", cd.Password)
	r.Equal("", cd.Port)
	r.Equal("", cd.User)
}

func Test_ConnectionDetails_Finalize_SQLite_Synonym_SynURL(t *testing.T) {
	r := require.New(t)
	cd := &ConnectionDetails{
		Dialect: "sqlite",
		URL:     "sqlite:///tmp/foo.db",
	}
	err := cd.Finalize()
	r.NoError(err)
	r.Equal("sqlite3", cd.Dialect, "given dialect: sqlite")
	r.Equal("/tmp/foo.db", cd.Database, "given url: sqlite:///tmp/foo.db")
	r.Equal("", cd.Host)
	r.Equal("", cd.Password)
	r.Equal("", cd.Port)
	r.Equal("", cd.User)
}

func Test_ConnectionDetails_Finalize_SQLite_Synonym_Path(t *testing.T) {
	r := require.New(t)
	cd := &ConnectionDetails{
		Dialect:  "sqlite",
		Database: "./foo.db",
	}
	err := cd.Finalize()
	r.NoError(err)
	r.Equal("sqlite3", cd.Dialect, "given dialect: sqlite")
	r.Equal("./foo.db", cd.Database, "given database: ./foo.db")
	r.Equal("", cd.Host)
	r.Equal("", cd.Password)
	r.Equal("", cd.Port)
	r.Equal("", cd.User)
}
