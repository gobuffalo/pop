package cmd

import (
	"errors"

	"github.com/spf13/cobra"
)

var migrationType string

var migrateCreateCmd = &cobra.Command{
	Use:   "create [name]",
	Short: "Generates Up/Down migrations for your database.",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return errors.New("You must supply a name for your migration!")
		}

		c := getConn()
		return c.MigrationCreate(migrationPath, args[0], migrationType)
	},
}

func init() {
	migrateCmd.AddCommand(migrateCreateCmd)
	migrateCreateCmd.Flags().StringVarP(&migrationType, "type", "t", "fizz", "Which type of migration file do you want? fizz or sql?")
}
