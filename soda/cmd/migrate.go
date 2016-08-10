package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var migrationPath string

var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Runs migrations against your database.",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		RootCmd.PersistentPreRun(cmd, args)
		return os.MkdirAll(migrationPath, 0766)
	},
	Run: func(cmd *cobra.Command, args []string) {
		c := getConn()
		err := c.MigrateUp(migrationPath)
		if err != nil {
			fmt.Println(err)
			os.Exit(-1)
		}
	},
}

func init() {
	RootCmd.AddCommand(migrateCmd)
	migrateCmd.PersistentFlags().StringVarP(&migrationPath, "path", "p", "./migrations", "Path to the migrations folder")
}
