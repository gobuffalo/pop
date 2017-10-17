package cmd

import (
	"github.com/markbates/pop"
	"github.com/spf13/cobra"
)

var all bool

var dropCmd = &cobra.Command{
	Use:   "drop",
	Short: "Drops databases for you",
	RunE: func(cmd *cobra.Command, args []string) error {
		var err error
		if all {
			for _, conn := range pop.Connections {
				err = pop.DropDB(conn)
				if err != nil {
					return err
				}
			}
		} else {
			err = pop.DropDB(getConn())
		}
		return err
	},
}

func init() {
	dropCmd.Flags().BoolVarP(&all, "all", "a", false, "Drops all of the databases in the database.yml")
	RootCmd.AddCommand(dropCmd)
}
