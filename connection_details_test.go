package pop_test

import (
	"testing"

	"github.com/gobuffalo/pop"
	"github.com/stretchr/testify/require"
)

func Test_ConnectionDetails_Finalize(t *testing.T) {
	r := require.New(t)

	cd := &pop.ConnectionDetails{
		URL: "postgres://user:pass@host:port/database",
	}
	err := cd.Finalize()
	r.NoError(err)

	r.Equal(cd.Database, "database")
	r.Equal(cd.Dialect, "postgres")
	r.Equal(cd.Host, "host")
	r.Equal(cd.Password, "pass")
	r.Equal(cd.Port, "port")
	r.Equal(cd.User, "user")
}

func Test_ConnectionDetails_Finalize_MySQL_DSN(t *testing.T) {
	r := require.New(t)

	url := "mysql://user:pass@(host:port)/database?param1=value1&param2=value2"
	cd := &pop.ConnectionDetails{
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

func Test_ConnectionDetails_Finalize_MySQL_DSN_Protocol(t *testing.T) {
	r := require.New(t)

	url := "mysql://user:pass@tcp(host:port)/protocol"
	cd := &pop.ConnectionDetails{
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
	cd := &pop.ConnectionDetails{
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
	cd := &pop.ConnectionDetails{
		URL: "unknown://user:pass@host:port/database",
	}
	err := cd.Finalize()
	r.Error(err)
}

func Test_ConnectionDetails_Finalize_SQLite(t *testing.T) {
	r := require.New(t)

	cd := &pop.ConnectionDetails{
		URL: "sqlite3:///tmp/foo.db",
	}
	err := cd.Finalize()
	r.NoError(err)

	r.Equal(cd.Database, "/tmp/foo.db")
	r.Equal(cd.Dialect, "sqlite3")
	r.Equal(cd.Host, "")
	r.Equal(cd.Password, "")
	r.Equal(cd.Port, "")
	r.Equal(cd.User, "")
}
