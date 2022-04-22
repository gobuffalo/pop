package cmd

import (
	"github.com/gobuffalo/pop/v6"
	"github.com/spf13/cobra"
)

var migrationStepDown int

var migrateDownCmd = &cobra.Command{
	Use:   "down",
	Short: "Apply one or more of the 'down' migrations.",
	RunE: func(cmd *cobra.Command, args []string) error {
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
