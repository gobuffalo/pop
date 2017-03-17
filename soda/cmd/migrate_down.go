package cmd

import "github.com/spf13/cobra"

var migrationStep int

var migrateDownCmd = &cobra.Command{
	Use:   "down",
	Short: "Apply one or more of the 'down' migrations.",
	RunE: func(cmd *cobra.Command, args []string) error {
		c := getConn()
		return c.MigrateDown(migrationPath, migrationStep)
	},
}

func init() {
	migrateCmd.AddCommand(migrateDownCmd)
	migrateDownCmd.Flags().IntVarP(&migrationStep, "step", "s", 1, "Number of migration to down")
}
