package cmd

import (
	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/pop/log"
	"github.com/spf13/cobra"
)

var all bool

var dropCmd = &cobra.Command{
	Use:   "drop",
	Short: "Drops databases for you",
	Run: func(cmd *cobra.Command, args []string) {
		var err error
		if all {
			for env, conn := range pop.Connections {
				err = pop.DropDB(conn)
				if err != nil {
					log.DefaultLogger.WithField("environment", env).WithField("error", err).Error("Failed to drop database")
				}
			}
		} else {
			if err := pop.DropDB(getConn()); err != nil {
				log.DefaultLogger.WithField("error", err).Error("Failed to drop database")
			}
		}
	},
}

func init() {
	dropCmd.Flags().BoolVarP(&all, "all", "a", false, "Drops all of the databases in the database.yml")
	RootCmd.AddCommand(dropCmd)
}
