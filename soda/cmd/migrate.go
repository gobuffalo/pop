package cmd

import (
	"os"

	"errors"

	"github.com/gobuffalo/pop"
	"github.com/spf13/cobra"
)

var migrationPath string

var migrateCmd = &cobra.Command{
	Use:     "migrate",
	Aliases: []string{"m"},
	Short:   "Runs migrations against your database.",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		RootCmd.PersistentPreRun(cmd, args)
		return os.MkdirAll(migrationPath, 0766)
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) > 0 {
			return errors.New("migrate command does not accept any argument")
		}
		mig, err := pop.NewFileMigrator(migrationPath, getConn())
		if err != nil {
			return err
		}
		return mig.Up()
	},
}

func init() {
	RootCmd.AddCommand(migrateCmd)
	RootCmd.PersistentFlags().StringVarP(&migrationPath, "path", "p", "./migrations", "Path to the migrations folder")
}
