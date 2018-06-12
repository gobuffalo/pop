package cmd

import (
	"github.com/gobuffalo/pop"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var migrationVersion string

var migrateUpCmd = &cobra.Command{
	Use:   "up",
	Short: "Apply all of the 'up' migrations.",
	RunE: func(cmd *cobra.Command, args []string) error {
		mig, err := pop.NewFileMigrator(migrationPath, getConn())
		if err != nil {
			return errors.WithStack(err)
		}
		if migrationVersion != "" {
			return mig.Up(migrationVersion)
		}
		return mig.Up()
	},
}

func init() {
	migrateCmd.AddCommand(migrateUpCmd)
	migrateUpCmd.Flags().StringVarP(&migrationVersion, "version", "v", "", "Apply a specific migration")
}
