package pop

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_Migrator_UpLogOnly(t *testing.T) {
	if PDB == nil {
		t.Skip("skipping integration tests")
	}

	t.Run("logs a version shared by two migrations once", func(t *testing.T) {
		r := require.New(t)

		const version = "29991231235959"
		mtn := PDB.MigrationTableName()
		t.Cleanup(func() {
			_ = PDB.RawQuery(fmt.Sprintf("delete from %s where version = ?", mtn), version).Exec()
		})

		m := NewMigrator(PDB)
		m.UpMigrations.Migrations = []Migration{
			{Version: version, Name: "create_users", DBType: PDB.Dialect.Name()},
			{Version: version, Name: "create_users", DBType: "all"},
		}
		r.NoError(m.UpLogOnly())

		count, err := PDB.Where("version = ?", version).Count(mtn)
		r.NoError(err)
		r.Equal(1, count)
	})
}
