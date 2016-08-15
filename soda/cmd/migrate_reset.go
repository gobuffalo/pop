package cmd

import "github.com/spf13/cobra"

var migrateResetCmd = &cobra.Command{
	Use:   "reset",
	Short: "The equivalent of running `migrate down` and then `migrate up`",
	RunE: func(cmd *cobra.Command, args []string) error {
		c := getConn()
		return c.MigrateReset(migrationPath)
	},
}

func init() {
	migrateCmd.AddCommand(migrateResetCmd)
}
