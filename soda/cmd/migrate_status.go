package cmd

import (
	"os"

	"github.com/spf13/cobra"

	"github.com/gobuffalo/pop/v6"
)

var migrateStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Displays the status of all migrations.",
	RunE: func(_ *cobra.Command, _ []string) error {
		mig, err := pop.NewFileMigrator(migrationPath, getConn())
		if err != nil {
			return err
		}
		return mig.Status(os.Stdout)
	},
}

func init() {
	migrateCmd.AddCommand(migrateStatusCmd)
}
