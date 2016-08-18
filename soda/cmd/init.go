package cmd

import (
	"github.com/markbates/pop/soda/cmd/generate"
	"github.com/spf13/cobra"
)

var dialect string

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Sets up a new Pop/Soda Project",
	RunE: func(cmd *cobra.Command, args []string) error {
		err := generate.ConfigCmd.RunE(cmd, args)
		return err
	},
}

func init() {
	RootCmd.AddCommand(initCmd)
	initCmd.Flags().StringVarP(&dialect, "type", "t", "postgres", "What type of database do you want to use? (postgres, mysql, sqlite3)")
}
