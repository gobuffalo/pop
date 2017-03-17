package schema

import (
	"os"

	"github.com/markbates/pop"
	"github.com/spf13/cobra"
)

var env string
var DumpCmd = &cobra.Command{
	Use:   "dump",
	Short: "A brief description of your command",
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := pop.Connect(env)
		if err != nil {
			return err
		}
		err = c.Dialect.DumpSchema(os.Stdout)
		if err != nil {
			return err
		}
		return nil
	},
}

func init() {
	DumpCmd.Flags().StringVarP(&env, "env", "e", "development", "The environment you want to run schema against. Will use $GO_ENV if set.")
}
