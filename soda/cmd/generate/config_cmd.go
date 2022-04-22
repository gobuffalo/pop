package generate

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/gobuffalo/genny/v2"
	"github.com/gobuffalo/pop/v6"
	"github.com/gobuffalo/pop/v6/genny/config"
	"github.com/gobuffalo/pop/v6/internal/defaults"
	"github.com/spf13/cobra"
)

func init() {
	ConfigCmd.Flags().StringVarP(&dialect, "type", "t", "postgres", fmt.Sprintf("The type of database you want to use (%s)", strings.Join(pop.AvailableDialects, ", ")))
}

var dialect string

// ConfigCmd is the command to generate pop config files
var ConfigCmd = &cobra.Command{
	Use:              "config",
	Short:            "Generates a database.yml file for your project.",
	PersistentPreRun: func(c *cobra.Command, args []string) {},
	RunE: func(cmd *cobra.Command, args []string) error {
		cflag := cmd.Flag("config")
		cflagVal := ""
		if cflag != nil {
			cflagVal = cflag.Value.String()
		}
		cfgFile := defaults.String(cflagVal, "database.yml")

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
		run.With(g)

		return run.Run()
	},
}
