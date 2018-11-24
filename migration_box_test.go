package pop

import (
	"testing"

	"github.com/gobuffalo/packr/v2"
	"github.com/stretchr/testify/require"
)

func Test_MigrationBox(t *testing.T) {
	r := require.New(t)

	b, err := NewMigrationBox(packr.New("./migrations/multiple", "./migrations/multiple"), PDB)
	r.NoError(err)
	r.Equal(4, len(b.Migrations["up"]))
	r.Equal("mysql", b.Migrations["up"][0].DBType)
	r.Equal("postgres", b.Migrations["up"][1].DBType)
	r.Equal("sqlite3", b.Migrations["up"][2].DBType)
	r.Equal("all", b.Migrations["up"][3].DBType)
}
