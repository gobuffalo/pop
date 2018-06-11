package schema

import (
	"os"

	"github.com/gobuffalo/pop"
	"github.com/spf13/cobra"
)

var loadOptions = struct {
	env   string
	input string
}{}

// LoadCmd loads a schema.sql file into a database.
var LoadCmd = &cobra.Command{
	Use:   "load",
	Short: "Load a schema.sql file into a database",
	RunE: func(cmd *cobra.Command, args []string) error {
		f, err := os.Open(loadOptions.input)
		if err != nil {
			return err
		}

		c, err := pop.Connect(loadOptions.env)
		if err != nil {
			return err
		}
		defer c.Close()

		return c.Dialect.LoadSchema(f)
	},
}

func init() {
	LoadCmd.Flags().StringVarP(&loadOptions.env, "env", "e", "development", "The environment you want to run schema against. Will use $GO_ENV if set.")
	LoadCmd.Flags().StringVarP(&loadOptions.input, "input", "i", "schema.sql", "The path to the schema file you want to load")
}
