// +build !appengine

package pop

import (
	"fmt"

	"github.com/gobuffalo/fizz"
)

func newSchemaMigrations(name string) fizz.Table {
	return fizz.Table{
		Name: name,
		Columns: []fizz.Column{
			{
				Name:    "version",
				ColType: "string",
				Options: map[string]interface{}{
					"size": 14, // len(YYYYMMDDhhmmss)
				},
			},
		},
		Indexes: []fizz.Index{
			{Name: fmt.Sprintf("%s_version_idx", name), Columns: []string{"version"}, Unique: true},
		},
	}
}
