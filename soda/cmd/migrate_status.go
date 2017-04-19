package cmd

import "github.com/spf13/cobra"

var migrateStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Displays the status of all migrations.",
	RunE: func(cmd *cobra.Command, args []string) error {
		c := getConn()
		return c.MigrateStatus(migrationPath)
	},
}

func init() {
	migrateCmd.AddCommand(migrateStatusCmd)
}
