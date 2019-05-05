package cmd

import (
	"github.com/gobuffalo/pop"
	"github.com/markbates/oncer"
	"github.com/spf13/cobra"
)

// Deprecated: use resetCmd instead.
var migrateResetCmd = &cobra.Command{
	Use:   "reset",
	Short: "The equivalent of running `migrate down` and then `migrate up`",
	RunE: func(cmd *cobra.Command, args []string) error {
		oncer.Deprecate(0, "command `migrate reset`", "Use command `reset` instead.")
		mig, err := pop.NewFileMigrator(migrationPath, getConn())
		if err != nil {
			return err
		}
		return mig.Reset()
	},
}

func init() {
	migrateCmd.AddCommand(migrateResetCmd)
}
