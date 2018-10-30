package pop

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_Cockroach_URL_Raw(t *testing.T) {
	r := require.New(t)
	cd := &ConnectionDetails{
		Dialect: "cockroach",
		URL:     "scheme://user:pass@host:port/database?option1=value1",
	}
	err := cd.Finalize()
	r.NoError(err)
	m := &cockroach{ConnectionDetails: cd}
	r.Equal("scheme://user:pass@host:port/database?option1=value1", m.URL())
	r.Equal("postgres://user:pass@host:port/?option1=value1", m.urlWithoutDb())
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
	m := &cockroach{ConnectionDetails: cd}
	r.Equal("postgres://user:pass@host:port/database?application_name=cockroach&option1=value1", m.URL())
	r.Equal("postgres://user:pass@host:port/?application_name=cockroach&option1=value1", m.urlWithoutDb())
}
