// Package generate provides commands for generating various files related to the pop package.
package generate

import (
	"cmp"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/gobuffalo/genny/v2"
	"github.com/spf13/cobra"

	"github.com/gobuffalo/pop/v6"
	"github.com/gobuffalo/pop/v6/genny/config"
)

func init() {
	ConfigCmd.Flags().StringVarP(&dialect, "type", "t", "postgres",
		fmt.Sprintf("The type of database you want to use (%s)", strings.Join(pop.AvailableDialects, ", ")),
	)
}

var dialect string

// ConfigCmd is the command to generate pop config files
var ConfigCmd = &cobra.Command{
	Use:   "config",
	Short: "Generates a database.yml file for your project.",
	RunE: func(cmd *cobra.Command, _ []string) error {
		cflag := cmd.Flag("config")
		cflagVal := ""
		if cflag != nil {
			cflagVal = cflag.Value.String()
		}
		cfgFile := cmp.Or(cflagVal, "database.yml")

		run := genny.WetRunner(context.Background())

		pwd, _ := os.Getwd()
		g, err := config.New(&config.Options{
			Root:     pwd,
			Prefix:   filepath.Base(pwd),
			FileName: cfgFile,
			Dialect:  dialect,
		})
		if err != nil {
			return err
		}
		if err := run.With(g); err != nil {
			return err
		}
		return run.Run()
	},
}
