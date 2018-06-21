package cmd

import (
	"os"

	"github.com/gobuffalo/pop"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var resetOptions = struct {
	all   bool
	input string
}{}

var resetCmd = &cobra.Command{
	Use:   "reset",
	Short: "Drop, then recreate databases",
	RunE: func(cmd *cobra.Command, args []string) error {
		var schema *os.File
		useMigrations := false

		if _, err := os.Stat(resetOptions.input); err == nil {
			schema, err = os.Open(resetOptions.input)
			if err != nil {
				return err
			}
			defer schema.Close()
		} else {
			// Fallback to migrations
			useMigrations = true
		}

		if all {
			for _, conn := range pop.Connections {
				if err := doReset(conn, schema, useMigrations); err != nil {
					return err
				}
			}
		} else {
			c := getConn()
			if err := doReset(c, schema, useMigrations); err != nil {
				return err
			}
		}

		return nil
	},
}

func init() {
	resetCmd.Flags().BoolVarP(&resetOptions.all, "all", "a", false, "Drops and recreate all of the databases in the database.yml")
	resetCmd.Flags().StringVarP(&resetOptions.input, "input", "i", "schema.sql", "The path to the schema file you want to load")
	RootCmd.AddCommand(resetCmd)
}

func doReset(c *pop.Connection, f *os.File, useMigrations bool) error {
	if err := pop.DropDB(c); err != nil {
		return err
	}
	if err := pop.CreateDB(c); err != nil {
		return err
	}
	mig, err := pop.NewFileMigrator(migrationPath, getConn())
	if err != nil {
		return errors.WithStack(err)
	}
	if useMigrations {
		// Apply the migrations directly
		return mig.Up()
	}
	// Otherwise, use schema instead
	if err := c.Dialect.LoadSchema(f); err != nil {
		return err
	}
	// Then load migrations entries, without applying them
	return mig.UpLogOnly()
}
