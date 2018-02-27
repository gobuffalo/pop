package generate

import (
	"github.com/pkg/errors"

	"github.com/gobuffalo/pop"
	"github.com/spf13/cobra"
)

//FizzCmd generates a new fizz migration
var FizzCmd = &cobra.Command{
	Use:     "fizz [name]",
	Aliases: []string{"migration"},
	Short:   "Generates Up/Down migrations for your database using fizz.",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return errors.New("You must supply a name for your migration")
		}
		//Just the migration name
		if len(args) == 1 {
			return pop.MigrationCreate(migrationPath, args[0], "fizz", nil, nil)
		}

		structTag = "json" //Of no use here, but required to be "json" or "xml" by func newModel()
		m, err := newModelFromArgs(args)

		if err != nil {
			return err
		}

		return m.generateFizz(migrationPath)
	},
}
