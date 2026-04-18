// Package schema provides a command for interacting with the database schema, including dumping and loading the schema.
package schema

import (
	"errors"
	"io"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/gobuffalo/pop/v6"
)

var dumpOptions = struct {
	env    string
	output string
}{}

// DumpCmd dumps out the schema of the selected database.
var DumpCmd = &cobra.Command{
	Use:   "dump",
	Short: "Dumps out the schema of the selected database",
	RunE: func(cmd *cobra.Command, _ []string) error {
		env := cmd.Flag("env")
		if env == nil {
			return errors.New("env is required")
		}
		dumpOptions.env = env.Value.String()
		c, err := pop.Connect(dumpOptions.env)
		if err != nil {
			return err
		}
		var out io.Writer
		rollback := func() {}
		if dumpOptions.output == "-" {
			out = os.Stdout
		} else {
			err = os.MkdirAll(filepath.Dir(dumpOptions.output), 0o755)
			if err != nil {
				return err
			}
			out, err = os.Create(dumpOptions.output)
			if err != nil {
				return err
			}
			rollback = func() {
				_ = os.RemoveAll(dumpOptions.output)
			}
		}
		if err := c.Dialect.DumpSchema(out); err != nil {
			rollback()
			return err
		}
		return nil
	},
}

func init() {
	DumpCmd.Flags().
		StringVarP(&dumpOptions.output, "output", "o", "./migrations/schema.sql", "The path to dump the schema to.")
}
