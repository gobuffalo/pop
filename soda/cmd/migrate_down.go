package cmd

import "github.com/spf13/cobra"

var migrateDownCmd = &cobra.Command{
	Use:   "down",
	Short: "Apply all of the 'down' migrations.",
	RunE: func(cmd *cobra.Command, args []string) error {
		c := getConn()
		return c.MigrateDown(migrationPath)
	},
}

func init() {
	migrateCmd.AddCommand(migrateDownCmd)
}
