package pop

import "fmt"

// Migration handles the data for a given database migration
type Migration struct {
	// Path to the migration (./migrations/123_create_widgets.up.sql)
	Path string
	// Version of the migration (123)
	Version string
	// Name of the migration (create_widgets)
	Name string
	// Direction of the migration (up)
	Direction string
	// Type of migration (sql)
	Type string
	// DB type (all|postgres|mysql...)
	DBType string
	// Runner function to run/execute the migration
	Runner func(Migration, *Connection) error
}

// Run the migration. Returns an error if there is
// no mf.Runner defined.
func (mf Migration) Run(c *Connection) error {
	if mf.Runner == nil {
		return fmt.Errorf("no runner defined for %s", mf.Path)
	}
	return mf.Runner(mf, c)
}

// Migrations is a collection of Migration
type Migrations []Migration

// Len returns the number of migrations in the collection.
func (mfs Migrations) Len() int {
	return len(mfs)
}

// Swap swaps the position of two migrations in the collection.
func (mfs Migrations) Swap(i, j int) {
	mfs[i], mfs[j] = mfs[j], mfs[i]
}

// Filter filters the migrations in place based on the provided function.
func (mfs *Migrations) Filter(f func(mf Migration) bool) {
	vsf := make(Migrations, 0)
	for _, v := range *mfs {
		if f(v) {
			vsf = append(vsf, v)
		}
	}
	*mfs = vsf
}

type (
	// UpMigrations is used to sort up migrations in order.
	UpMigrations struct {
		Migrations
	}

	// DownMigrations is used to sort down migrations in reverse order.
	DownMigrations struct {
		Migrations
	}
)

// Less is used to sort up migrations.
func (mfs UpMigrations) Less(i, j int) bool {
	if mfs.Migrations[i].Version == mfs.Migrations[j].Version {
		// force "all" to the back
		return mfs.Migrations[i].DBType != "all"
	}
	return mfs.Migrations[i].Version < mfs.Migrations[j].Version
}

// Less is used to sort down migrations in reverse order.
func (mfs DownMigrations) Less(i, j int) bool {
	if mfs.Migrations[i].Version == mfs.Migrations[j].Version {
		// force "all" to the back
		return mfs.Migrations[i].DBType != "all"
	}
	return mfs.Migrations[i].Version > mfs.Migrations[j].Version
}
