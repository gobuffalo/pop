package pop

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_Cockroach_URL_Raw(t *testing.T) {
	r := require.New(t)
	cd := &ConnectionDetails{
		Dialect: "cockroach",
		URL:     "scheme://user:pass@host:1337/database?option1=value1",
	}
	err := cd.Finalize()
	r.NoError(err)
	m := &cockroach{commonDialect: commonDialect{ConnectionDetails: cd}}
	r.Equal("scheme://user:pass@host:1337/database?option1=value1", m.URL())
	r.Equal("postgres://user:pass@host:1337/?option1=value1", m.urlWithoutDb())
}

func Test_Cockroach_URL_Build(t *testing.T) {
	r := require.New(t)

	cd := &ConnectionDetails{
		Dialect:  "cockroach",
		Database: "database",
		Host:     "host",
		Port:     "port",
		User:     "user",
		Password: "pass",
		Options: map[string]string{
			"option1": "value1",
		},
	}
	err := cd.Finalize()
	r.NoError(err)
	m := &cockroach{commonDialect: commonDialect{ConnectionDetails: cd}}
	r.True(strings.HasPrefix(m.URL(), "postgres://user:pass@host:port/database?"), "URL() returns %v", m.URL())
	r.Contains(m.URL(), "option1=value1")
	r.Contains(m.URL(), "application_name=pop.test")
	r.True(strings.HasPrefix(m.urlWithoutDb(), "postgres://user:pass@host:port/?"), "urlWithoutDb() returns %v", m.urlWithoutDb())
	r.Contains(m.urlWithoutDb(), "option1=value1")
	r.Contains(m.urlWithoutDb(), "application_name=pop.test")
	r.True(strings.HasPrefix(m.MigrationURL(), "postgres://user:pass@host:port/database?"), "MigrationURL() returns %v", m.MigrationURL())
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
