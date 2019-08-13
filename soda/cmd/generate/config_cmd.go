package generate

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/gobuffalo/genny"
	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/pop/genny/config"
	"github.com/markbates/going/defaults"
	"github.com/markbates/oncer"
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

// GenerateConfig generates pop configuration files.
//
// Deprecated: use github.com/gobuffalo/pop/genny/config instead.
func GenerateConfig(cfgFile string, data map[string]interface{}) error {
	oncer.Deprecate(0, "generate.GenerateConfig", "Use github.com/gobuffalo/pop/genny/config instead.")
	return Config(cfgFile, data)
}

// Config generates pop configuration files.
// Deprecated: use github.com/gobuffalo/pop/genny/config instead.
func Config(cfgFile string, data map[string]interface{}) error {
	oncer.Deprecate(0, "generate.Config", "Use github.com/gobuffalo/pop/genny/config instead.")
	pwd, _ := os.Getwd()

	run := genny.WetRunner(context.Background())

	d, _ := data["dialect"].(string)
	g, err := config.New(&config.Options{
		Root:     pwd,
		Prefix:   filepath.Base(pwd),
		FileName: cfgFile,
		Dialect:  d,
	})

	if err != nil {
		return err
	}
	run.With(g)

	return run.Run()
}
