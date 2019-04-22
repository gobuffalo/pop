package cmd

import (
	"github.com/gobuffalo/pop"
	"github.com/spf13/cobra"
)

var migrateUpCmd = &cobra.Command{
	Use:   "up",
	Short: "Apply all of the 'up' migrations.",
	RunE: func(cmd *cobra.Command, args []string) error {
		mig, err := pop.NewFileMigrator(migrationPath, getConn())
		if err != nil {
			return err
		}
		return mig.Up()
	},
}

func init() {
	migrateCmd.AddCommand(migrateUpCmd)
}
