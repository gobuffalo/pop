package cmd

import "github.com/spf13/cobra"

var migrateUpCmd = &cobra.Command{
	Use:   "up",
	Short: "Apply all of the 'up' migrations.",
	RunE: func(cmd *cobra.Command, args []string) error {
		c := getConn()
		return c.MigrateUp(migrationPath)
	},
}

func init() {
	migrateCmd.AddCommand(migrateUpCmd)
}
