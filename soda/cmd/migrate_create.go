package cmd

import (
	"fmt"

	"github.com/pkg/errors"

	"github.com/markbates/pop"
	"github.com/spf13/cobra"
)

var migrationType string

var depWarning = "[DEPRACTION WARNING] This command is deprecated. Please use `soda generate fizz` or `soda generate sql` instead."

var migrateCreateCmd = &cobra.Command{
	Use:   "create [name]",
	Short: depWarning,
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println(depWarning)
		if len(args) == 0 {
			return errors.New("You must supply a name for your migration!")
		}

		return pop.MigrationCreate(migrationPath, args[0], migrationType, nil, nil)
	},
}

func init() {
	migrateCmd.AddCommand(migrateCreateCmd)
	migrateCreateCmd.Flags().StringVarP(&migrationType, "type", "t", "fizz", "Which type of migration file do you want? fizz or sql?")
}
