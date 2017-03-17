package cmd

import (
	"os"

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
		c := getConn()
		return c.MigrateUp(migrationPath)
	},
}

func init() {
	RootCmd.AddCommand(migrateCmd)
	RootCmd.PersistentFlags().StringVarP(&migrationPath, "path", "p", "./migrations", "Path to the migrations folder")
}
