package pop

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func Test_ParseMigrationFilenameUp(t *testing.T) {
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

func Test_ParseMigrationFilenameDown(t *testing.T) {
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

func Test_ParseMigrationFilenameUpPostgres(t *testing.T) {
	r := require.New(t)

	m, err := ParseMigrationFilename("20190611004000_create_providers.pg.up.fizz")
	r.NoError(err)
	r.NotNil(m)
	r.Equal(m.Version, "20190611004000")
	r.Equal(m.Name, "create_providers")
	r.Equal(m.DBType, "postgres")
	r.Equal(m.Direction, "up")
	r.Equal(m.Type, "fizz")
}
