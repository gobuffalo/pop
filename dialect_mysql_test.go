package pop

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_MySQL_URL_As_Is(t *testing.T) {
	r := require.New(t)

	cd := &ConnectionDetails{
		URL: "mysql://user:pass@(host:port)/dbase?opt=value",
	}
	err := cd.Finalize()
	r.NoError(err)

	m := &mysql{ConnectionDetails: cd}
	r.Equal("user:pass@(host:port)/dbase?opt=value", m.URL())
	r.Equal("user:pass@(host:port)/?opt=value", m.urlWithoutDb())
}

func Test_MySQL_withURL(t *testing.T) {
	r := require.New(t)

	cd := &ConnectionDetails{
		Database: "xx",
		Host:     "xx",
		Port:     "xx",
		User:     "xx",
		Password: "xx",
		URL:      "mysql://user:pass@(host:port)/dbase?opt=value",
	}
	err := cd.Finalize()
	r.NoError(err)

	m := &mysql{ConnectionDetails: cd}
	r.Equal("user:pass@(host:port)/dbase?opt=value", m.URL())
	r.Equal("user:pass@(host:port)/?opt=value", m.urlWithoutDb())
}

func Test_MySQL_Default_Options(t *testing.T) {
	r := require.New(t)

	cd := &ConnectionDetails{
		Dialect:  "mysql",
		Database: "base",
		Host:     "host",
		Port:     "port",
		User:     "user",
		Password: "pass",
	}
	err := cd.Finalize()
	r.NoError(err)

	m := &mysql{ConnectionDetails: cd}
	r.True(strings.HasPrefix(m.URL(), "user:pass@(host:port)/base?"))
	r.Contains(m.URL(), "multiStatements=true")
	r.Contains(m.URL(), "parseTime=true")
	r.Contains(m.URL(), "readTimeout=1s")
	r.Contains(m.URL(), "collation=utf8mb4_general_ci")
	r.True(strings.HasPrefix(m.urlWithoutDb(), "user:pass@(host:port)/?"))
	r.Contains(m.urlWithoutDb(), "multiStatements=true")
	r.Contains(m.urlWithoutDb(), "parseTime=true")
	r.Contains(m.urlWithoutDb(), "readTimeout=1s")
	r.Contains(m.urlWithoutDb(), "collation=utf8mb4_general_ci")
}

func Test_MySQL_User_Defined_Options(t *testing.T) {
	r := require.New(t)

	cd := &ConnectionDetails{
		Dialect:  "mysql",
		Database: "base",
		Host:     "host",
		Port:     "port",
		User:     "user",
		Password: "pass",
		Options: map[string]string{
			"multiStatements": "false",
			"parseTime":       "false",
			"readTimeout":     "1h",
			"collation":       "utf8",
		},
	}
	err := cd.Finalize()
	r.NoError(err)

	m := &mysql{ConnectionDetails: cd}
	r.True(strings.HasPrefix(m.URL(), "user:pass@(host:port)/base?"))
	r.Contains(m.URL(), "multiStatements=false")
	r.Contains(m.URL(), "parseTime=false")
	r.Contains(m.URL(), "readTimeout=1h")
	r.Contains(m.URL(), "collation=utf8")
	r.True(strings.HasPrefix(m.urlWithoutDb(), "user:pass@(host:port)/?"))
	r.Contains(m.urlWithoutDb(), "multiStatements=false")
	r.Contains(m.urlWithoutDb(), "parseTime=false")
	r.Contains(m.urlWithoutDb(), "readTimeout=1h")
	r.Contains(m.urlWithoutDb(), "collation=utf8")
}

// preserve this test case while deprecated code alives
func Test_MySQL_Deprecated(t *testing.T) {
	r := require.New(t)

	cd := &ConnectionDetails{
		Dialect:  "mysql",
		Database: "base",
		Host:     "host",
		Port:     "port",
		User:     "user",
		Password: "pass",
		Encoding: "myEncoding",
	}
	err := cd.Finalize()
	r.NoError(err)
	r.NotNil(cd.Options)
	r.Equal("myEncoding", cd.Options["collation"])
}
