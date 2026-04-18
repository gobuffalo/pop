package cmd

import (
	"github.com/spf13/cobra"

	"github.com/gobuffalo/pop/v6"
)

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Creates databases for you",
	RunE: func(_ *cobra.Command, _ []string) error {
		if !all {
			return pop.CreateDB(getConn())
		}
		for _, conn := range pop.Connections {
			err := pop.CreateDB(conn)
			if err != nil {
				return err
			}
		}
		return nil
	},
}

func init() {
	createCmd.Flags().BoolVarP(&all, "all", "a", false, "Creates all of the databases in the database.yml")
	RootCmd.AddCommand(createCmd)
}
