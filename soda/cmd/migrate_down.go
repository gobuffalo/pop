package cmd

import (
	"github.com/gobuffalo/pop"
	"github.com/spf13/cobra"
)

var migrationStep int

var migrateDownCmd = &cobra.Command{
	Use:   "down",
	Short: "Apply one or more of the 'down' migrations.",
	RunE: func(cmd *cobra.Command, args []string) error {
		mig, err := pop.NewFileMigrator(migrationPath, getConn())
		if err != nil {
			return err
		}
		return mig.Down(migrationStep)
	},
}

func init() {
	migrateCmd.AddCommand(migrateDownCmd)
	migrateDownCmd.Flags().IntVarP(&migrationStep, "step", "s", 1, "Number of migration to down")
}
