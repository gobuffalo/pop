package cmd

import (
	"fmt"
	"strings"

	"bufio"
	"os"

	"github.com/gobuffalo/pop"
	"github.com/spf13/cobra"
)

var all bool
var confirmed bool

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

		if !confirmed {
			reader := bufio.NewReader(os.Stdin)
			fmt.Print("Do you really want to drop the database [y/N]? ")
			r, _ := reader.ReadString('\n')
			r = strings.TrimSpace(r)
			if r != "y" && r != "Y" {
				fmt.Println("Aborting due to lack of user confirmation.")
				os.Exit(0)
			}

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
	dropCmd.Flags().BoolVarP(&confirmed, "yes", "y", false, "Runs without asking the user for confirmation")
	RootCmd.AddCommand(dropCmd)
}
