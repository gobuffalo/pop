package cmd

import (
	"github.com/gobuffalo/pop"
	"github.com/spf13/cobra"
)

var migrateStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Displays the status of all migrations.",
	RunE: func(cmd *cobra.Command, args []string) error {
		mig, err := pop.NewFileMigrator(migrationPath, getConn())
		if err != nil {
			return err
		}
		return mig.Status()
	},
}

func init() {
	migrateCmd.AddCommand(migrateStatusCmd)
}
