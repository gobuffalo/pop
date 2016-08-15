package cmd

import (
	"github.com/markbates/pop"
	"github.com/spf13/cobra"
)

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Creates databases for you",
	Run: func(cmd *cobra.Command, args []string) {
		if all {
			for _, conn := range pop.Connections {
				pop.CreateDB(conn)
			}
		} else {
			pop.CreateDB(getConn())
		}
	},
}

func init() {
	createCmd.Flags().BoolVarP(&all, "all", "a", false, "Creates all of the databases in the database.yml")
	RootCmd.AddCommand(createCmd)
}
