package cmd

import (
	"os"

	"github.com/gobuffalo/pop"
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
		f, err := os.Open(resetOptions.input)
		if err != nil {
			return err
		}
		defer f.Close()

		doReset := func(c *pop.Connection, f *os.File) error {
			if err := pop.DropDB(c); err != nil {
				return err
			}
			if err := pop.CreateDB(c); err != nil {
				return err
			}
			if err := c.Dialect.LoadSchema(f); err != nil {
				return err
			}
			return nil
		}

		if all {
			for _, conn := range pop.Connections {
				if err := doReset(conn, f); err != nil {
					return err
				}
			}
		} else {
			c := getConn()
			if err := doReset(c, f); err != nil {
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
