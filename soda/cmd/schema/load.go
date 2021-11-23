package schema

import (
	"errors"
	"fmt"
	"os"

	"github.com/gobuffalo/pop/v6"
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
		env := cmd.Flag("env")
		if env == nil {
			return errors.New("env is required")
		}
		loadOptions.env = env.Value.String()

		f, err := os.Open(loadOptions.input)
		if err != nil {
			return fmt.Errorf("unable to load schema file: %w", err)
		}
		defer f.Close()

		c, err := pop.Connect(loadOptions.env)
		if err != nil {
			return fmt.Errorf("unable to connect to database: %w", err)
		}
		defer c.Close()

		return c.Dialect.LoadSchema(f)
	},
}

func init() {
	LoadCmd.Flags().StringVarP(&loadOptions.input, "input", "i", "./migrations/schema.sql", "The path to the schema file you want to load")
}
