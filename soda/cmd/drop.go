package cmd

import (
	"fmt"
	"os"

	"github.com/gobuffalo/pop/v6"
	"github.com/spf13/cobra"
)

var all bool

var dropCmd = &cobra.Command{
	Use:   "drop",
	Short: "Drops databases for you",
	Run: func(cmd *cobra.Command, args []string) {
		var err error
		if len(args) > 0 {
			err = fmt.Errorf("no arguments allowed with the drop database command")
			fmt.Println(err)
			os.Exit(1)
		}

		if all {
			for _, conn := range pop.Connections {
				err = pop.DropDB(conn)
				if err != nil {
					fmt.Println(err)
				}
			}
		} else {
			if err := pop.DropDB(getConn()); err != nil {
				fmt.Println(err)
			}
		}
	},
}

func init() {
	dropCmd.Flags().BoolVarP(&all, "all", "a", false, "Drops all of the databases in the database.yml")
	RootCmd.AddCommand(dropCmd)
}
