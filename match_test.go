package pop

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_ParseMigrationFilenameFizzDown(t *testing.T) {
	r := require.New(t)

	m, err := ParseMigrationFilename("20190611004000_create_providers.down.fizz")
	r.NoError(err)
	r.NotNil(m)
	r.Equal(m.Version, "20190611004000")
	r.Equal(m.Name, "create_providers")
	r.Equal(m.DBType, "all")
	r.Equal(m.Direction, "down")
	r.Equal(m.Type, "fizz")
}

func Test_ParseMigrationFilenameFizzUp(t *testing.T) {
	r := require.New(t)

	m, err := ParseMigrationFilename("20190611004000_create_providers.up.fizz")
	r.NoError(err)
	r.NotNil(m)
	r.Equal(m.Version, "20190611004000")
	r.Equal(m.Name, "create_providers")
	r.Equal(m.DBType, "all")
	r.Equal(m.Direction, "up")
	r.Equal(m.Type, "fizz")
}

func Test_ParseMigrationFilenameFizzUpPostgres(t *testing.T) {
	r := require.New(t)

	m, err := ParseMigrationFilename("20190611004000_create_providers.pg.up.fizz")
	r.NotNil(err)
	r.Equal(err.Error(), "invalid database type \"postgres\", expected \"all\" because fizz is database type independent")
	r.Nil(m)
}

func Test_ParseMigrationFilenameFizzDownPostgres(t *testing.T) {
	r := require.New(t)

	m, err := ParseMigrationFilename("20190611004000_create_providers.pg.down.fizz")
	r.NotNil(err)
	r.Equal(err.Error(), "invalid database type \"postgres\", expected \"all\" because fizz is database type independent")
	r.Nil(m)
}

func Test_ParseMigrationFilenameSQLUp(t *testing.T) {
	r := require.New(t)

	m, err := ParseMigrationFilename("20190611004000_create_providers.up.sql")
	r.NoError(err)
	r.NotNil(m)
	r.Equal(m.Version, "20190611004000")
	r.Equal(m.Name, "create_providers")
	r.Equal(m.DBType, "all")
	r.Equal(m.Direction, "up")
	r.Equal(m.Type, "sql")
}

func Test_ParseMigrationFilenameSQLUpPostgres(t *testing.T) {
	r := require.New(t)

	m, err := ParseMigrationFilename("20190611004000_create_providers.pg.up.sql")
	r.NoError(err)
	r.NotNil(m)
	r.Equal(m.Version, "20190611004000")
	r.Equal(m.Name, "create_providers")
	r.Equal(m.DBType, "postgres")
	r.Equal(m.Direction, "up")
	r.Equal(m.Type, "sql")
}
