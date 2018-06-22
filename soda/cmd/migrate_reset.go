package cmd

import (
	"fmt"

	"github.com/gobuffalo/pop"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

// Deprecated: use resetCmd instead.
var migrateResetCmd = &cobra.Command{
	Use:   "reset",
	Short: "The equivalent of running `migrate down` and then `migrate up`",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("[POP] The `migrate reset` command has been deprecated! Use `reset` instead.")
		mig, err := pop.NewFileMigrator(migrationPath, getConn())
		if err != nil {
			return errors.WithStack(err)
		}
		return mig.Reset()
	},
}

func init() {
	migrateCmd.AddCommand(migrateResetCmd)
}
