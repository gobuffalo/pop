package pop

import (
	"strings"
	"testing"

	"github.com/gobuffalo/envy"

	"github.com/stretchr/testify/require"
)

func Test_LoadsConnectionsFromConfig(t *testing.T) {
	r := require.New(t)

	conns := Connections
	r.Equal(6, len(conns))
}

func Test_AddLookupPaths(t *testing.T) {
	r := require.New(t)
	AddLookupPaths("./foo")
	r.Contains(LookupPaths(), "./foo")
}

func Test_ParseConfig(t *testing.T) {
	r := require.New(t)
	config := strings.NewReader(`
mysql:
  dialect: "mysql"
  database: "pop_test"
  host: {{ envOr "MYSQL_HOST" "127.0.0.1"  }}
  port: {{ envOr "MYSQL_PORT" "3306"  }}
  user: {{ envOr "MYSQL_USER"  "root"  }}
  password: {{ envOr "MYSQL_PASSWORD"  "root"  }}
  options:
    readTimeout: 5s`)
	conns, err := ParseConfig(config)
	r.NoError(err)
	r.Equal(1, len(conns))
	r.NotNil(conns["mysql"])
	r.Equal("mysql", conns["mysql"].Dialect)
	r.Equal("pop_test", conns["mysql"].Database)
	r.Equal("127.0.0.1", conns["mysql"].Host)
	r.Equal("3306", conns["mysql"].Port)
	r.Equal(envy.Get("MYSQL_USER", "root"), conns["mysql"].User)
	r.Equal(envy.Get("MYSQL_PASSWORD", "root"), conns["mysql"].Password)
	r.Equal("5s", conns["mysql"].Options["readTimeout"])
}
