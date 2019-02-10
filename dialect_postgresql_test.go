package pop

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_PostgreSQL_Connection_String(t *testing.T) {
	r := require.New(t)

	url := "host=host port=port dbname=database user=user password=pass"
	cd := &ConnectionDetails{
		Dialect: "postgres",
		URL:     url,
	}
	err := cd.Finalize()
	r.NoError(err)

	r.Equal(url, cd.URL)
	r.Equal("postgres", cd.Dialect)
	r.Equal("host", cd.Host)
	r.Equal("pass", cd.Password)
	r.Equal("port", cd.Port)
	r.Equal("user", cd.User)
	r.Equal("database", cd.Database)
}

func Test_PostgreSQL_Connection_String_Options(t *testing.T) {
	r := require.New(t)

	url := "host=host port=port dbname=database user=user password=pass sslmode=disable fallback_application_name=test_app connect_timeout=10 sslcert=/some/location sslkey=/some/other/location sslrootcert=/root/location"
	cd := &ConnectionDetails{
		Dialect: "postgres",
		URL:     url,
	}
	err := cd.Finalize()
	r.NoError(err)

	r.Equal(url, cd.URL)

	r.Equal("disable", cd.Options["sslmode"])
	r.Equal("test_app", cd.Options["fallback_application_name"])
	r.Equal("10", cd.Options["connect_timeout"])
	r.Equal("/some/location", cd.Options["sslcert"])
	r.Equal("/some/other/location", cd.Options["sslkey"])
	r.Equal("/root/location", cd.Options["sslrootcert"])
}

func Test_PostgreSQL_Connection_String_Without_User(t *testing.T) {
	r := require.New(t)

	url := "dbname=database"
	cd := &ConnectionDetails{
		Dialect: "postgres",
		URL:     url,
	}
	err := cd.Finalize()
	r.NoError(err)

	r.Equal(url, cd.URL)
	r.Equal("postgres", cd.Dialect)
	r.Equal("", cd.Host)
	r.Equal("", cd.Password)
	r.Equal(portPostgreSQL, cd.Port) // fallback
	r.Equal("", cd.User)
	r.Equal("database", cd.Database)
}

func Test_PostgreSQL_Connection_String_Failure(t *testing.T) {
	r := require.New(t)

	url := "abc"
	cd := &ConnectionDetails{
		Dialect: "postgres",
		URL:     url,
	}
	err := cd.Finalize()
	r.Error(err)
	r.Equal("postgres", cd.Dialect)
}
