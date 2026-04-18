package cmd

import (
	"github.com/spf13/cobra"

	"github.com/gobuffalo/pop/v6"
)

var migrationStepDown int

var migrateDownCmd = &cobra.Command{
	Use:   "down",
	Short: "Apply one or more of the 'down' migrations.",
	RunE: func(_ *cobra.Command, _ []string) error {
		mig, err := pop.NewFileMigrator(migrationPath, getConn())
		if err != nil {
			return err
		}
		return mig.Down(migrationStepDown)
	},
}

func init() {
	migrateCmd.AddCommand(migrateDownCmd)
	migrateDownCmd.Flags().IntVarP(&migrationStepDown, "step", "s", 1, "Number of migration to down")
}
