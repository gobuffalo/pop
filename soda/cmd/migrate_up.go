package cmd

import (
	"github.com/spf13/cobra"

	"github.com/gobuffalo/pop/v6"
)

var migrationStepUp int

var migrateUpCmd = &cobra.Command{
	Use:   "up",
	Short: "Apply one or more of the 'up' migrations.",
	RunE: func(_ *cobra.Command, _ []string) error {
		mig, err := pop.NewFileMigrator(migrationPath, getConn())
		if err != nil {
			return err
		}
		_, err = mig.UpTo(migrationStepUp)
		return err
	},
}

func init() {
	migrateCmd.AddCommand(migrateUpCmd)
	migrateUpCmd.Flags().
		IntVarP(&migrationStepUp, "step", "s", 0, "Number of migrations to apply. Use 0 to apply all pending.")
}
