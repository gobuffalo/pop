//go:build !appengine

package pop

import (
	"github.com/gobuffalo/fizz"
)

func newSchemaMigrations(name string) fizz.Table {
	tab := fizz.Table{
		Name: name,
		Columns: []fizz.Column{
			{
				Name:    "version",
				ColType: "string",
				Options: map[string]any{
					"size": 14, // len(YYYYMMDDhhmmss)
				},
			},
		},
		Indexes: []fizz.Index{
			{Name: name + "_version_idx", Columns: []string{"version"}, Unique: true},
		},
	}
	// this is for https://github.com/gobuffalo/pop/issues/659.
	// primary key is not necessary for the migration table but it looks like
	// some database engine versions requires it for index.
	tab.PrimaryKey("version")
	return tab
}
