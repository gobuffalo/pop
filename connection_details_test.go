package pop

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_ConnectionDetails_Finalize(t *testing.T) {
	r := require.New(t)

	cd := &ConnectionDetails{
		URL: "postgres://user:pass@host:1234/database",
	}
	err := cd.Finalize()
	r.NoError(err)

	r.Equal("database", cd.Database)
	r.Equal("postgres", cd.Dialect)
	r.Equal("host", cd.Host)
	r.Equal("pass", cd.Password)
	r.Equal("1234", cd.Port)
	r.Equal("user", cd.User)
}

func Test_ConnectionDetails_Finalize_MySQL_Standard(t *testing.T) {
	r := require.New(t)

	url := "mysql://user:pass@(host:1337)/database?param1=value1&param2=value2"
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
	r.Equal("1337", cd.Port)
	r.Equal("database", cd.Database)
}

func Test_ConnectionDetails_Finalize_Cockroach(t *testing.T) {
	r := require.New(t)
	cd := &ConnectionDetails{
		Dialect: "cockroach",
		URL:     "postgres://user:pass@host:1234/database?sslmode=require&sslrootcert=certs/ca.crt&sslkey=certs/client.key&sslcert=certs/client.crt",
	}
	err := cd.Finalize()
	r.NoError(err)
	r.Equal("cockroach", cd.Dialect)
	r.Equal("database", cd.Database)
	r.Equal("host", cd.Host)
	r.Equal("1234", cd.Port)
	r.Equal("user", cd.User)
	r.Equal("pass", cd.Password)
}

func Test_ConnectionDetails_Finalize_UnknownSchemeURL(t *testing.T) {
	r := require.New(t)
	cd := &ConnectionDetails{
		URL: "unknown://user:pass@host:1337/database",
	}
	err := cd.Finalize()
	r.Error(err)
}

func Test_ConnectionDetails_Finalize_UnknownDialect(t *testing.T) {
	r := require.New(t)
	cd := &ConnectionDetails{
		Dialect: "unknown",
	}
	err := cd.Finalize()
	r.Error(err)
}

func Test_ConnectionDetails_Finalize_NoDB_NoURL(t *testing.T) {
	r := require.New(t)
	cd := &ConnectionDetails{
		Dialect: "sqlite3",
	}
	err := cd.Finalize()
	r.Error(err)
}

func Test_ConnectionDetails_OptionsString_Postgres(t *testing.T) {
	r := require.New(t)
	cd := &ConnectionDetails{
		Dialect:  "postgres",
		Database: "database",
		Host:     "host",
		Port:     "1234",
		User:     "user",
		Password: "pass",
		Options: map[string]string{
			"migration_table_name": "migrations",
			"sslmode":              "require",
		},
	}

	r.Equal("sslmode=require", cd.OptionsString(""))
	r.Equal("migrations", cd.MigrationTableName())
}
