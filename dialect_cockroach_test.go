package pop

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_Cockroach_ConnectionDetails_URL_Finalize(t *testing.T) {
	r := require.New(t)

	cd := &ConnectionDetails{
		Dialect: "cockroach",
		URL:     "cockroach://user:pass%23@host:1234/database",
	}
	err := cd.Finalize()
	r.NoError(err)

	r.Equal("database", cd.Database)
	r.Equal("cockroach", cd.Dialect)
	r.Equal("host", cd.Host)
	r.Equal("pass#", cd.Password)
	r.Equal("1234", cd.Port)
	r.Equal("user", cd.User)
}

func Test_Cockroach_ConnectionDetails_Values_Finalize(t *testing.T) {
	r := require.New(t)

	cd := &ConnectionDetails{
		Dialect:  "cockroach",
		Database: "database",
		Host:     "host",
		Port:     "1234",
		User:     "user",
		Password: "pass#",
		Options:  map[string]string{"application_name": "testing"},
	}
	err := cd.Finalize()
	r.NoError(err)

	p := &cockroach{commonDialect: commonDialect{ConnectionDetails: cd}}

	r.Equal("postgres://user:pass%23@host:1234/database?application_name=testing", p.URL())
}

func Test_Cockroach_URL_Raw(t *testing.T) {
	r := require.New(t)
	cd := &ConnectionDetails{
		Dialect: "cockroach",
		URL:     "cockroach://user:pass@host:1234/database?option1=value1",
	}
	err := cd.Finalize()
	r.NoError(err)
	m := &cockroach{commonDialect: commonDialect{ConnectionDetails: cd}}
	r.Equal("postgres://user:pass@host:1234/database?option1=value1", m.URL())
	r.Equal("postgres://user:pass@host:1234/?option1=value1", m.urlWithoutDb())
}

func Test_Cockroach_URL_Build(t *testing.T) {
	r := require.New(t)

	cd := &ConnectionDetails{
		Dialect:  "cockroach",
		Database: "database",
		Host:     "host",
		Port:     "port",
		User:     "user",
		Password: "pass#",
		Options: map[string]string{
			"option1": "value1",
		},
	}
	err := cd.Finalize()
	r.NoError(err)

	m := &cockroach{commonDialect: commonDialect{ConnectionDetails: cd}}
	r.True(strings.HasPrefix(m.URL(), "postgres://user:pass%23@host:port/database?"), "URL() returns %v", m.URL())
	r.Contains(m.URL(), "option1=value1")
	r.Contains(m.URL(), "application_name=pop.test")

	r.True(strings.HasPrefix(m.urlWithoutDb(), "postgres://user:pass%23@host:port/?"), "urlWithoutDb() returns %v", m.urlWithoutDb())
	r.Contains(m.urlWithoutDb(), "option1=value1")
	r.Contains(m.urlWithoutDb(), "application_name=pop.test")

	r.True(strings.HasPrefix(m.MigrationURL(), "postgres://user:pass%23@host:port/database?"), "MigrationURL() returns %v", m.MigrationURL())
}

func Test_Cockroach_URL_UserDefinedAppName(t *testing.T) {
	r := require.New(t)

	cd := &ConnectionDetails{
		Dialect:  "cockroach",
		Database: "database",
		Options: map[string]string{
			"application_name": "myapp",
		},
	}
	err := cd.Finalize()
	r.NoError(err)
	m := &cockroach{commonDialect: commonDialect{ConnectionDetails: cd}}
	r.Contains(m.URL(), "database?application_name=myapp")
	r.Contains(m.urlWithoutDb(), "/?application_name=myapp")
}

func Test_Cockroach_tableQuery(t *testing.T) {
	r := require.New(t)
	cr := cockroach{}

	cr.info.version = "v1.0.7"
	r.Equal(selectTablesQueryCockroachV1, cr.tablesQuery())

	cr.info.version = "v1.1.9"
	r.Equal(selectTablesQueryCockroachV1, cr.tablesQuery())

	cr.info.version = "v2.0.7"
	r.Equal(selectTablesQueryCockroach, cr.tablesQuery())

	cr.info.version = "v2.1.7"
	r.Equal(selectTablesQueryCockroach, cr.tablesQuery())

	cr.info.version = "v19.1.1"
	r.Equal(selectTablesQueryCockroach, cr.tablesQuery())

	cr.info.version = "v20.1.1"
	r.Equal(selectTablesQueryCockroach, cr.tablesQuery())
}

func Test_Cockroach_URL_Only(t *testing.T) {
	r := require.New(t)
	cd := &ConnectionDetails{
		URL: "cockroach://user:pass@host:1337/database?option1=value1",
	}
	err := cd.Finalize()
	r.NoError(err)
	m := &cockroach{commonDialect: commonDialect{ConnectionDetails: cd}}
	r.Equal("postgres://user:pass@host:1337/database?option1=value1", m.URL())
	r.Equal("postgres://user:pass@host:1337/?option1=value1", m.urlWithoutDb())
}
