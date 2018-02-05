package generate

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/gobuffalo/makr"
	"github.com/markbates/going/defaults"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

func init() {
	ConfigCmd.Flags().StringVarP(&dialect, "type", "t", "postgres", "What type of database do you want to use? (postgres, mysql, sqlite3, cockroach)")
}

var dialect string

//ConfigCmd is the command to generate pop config files
var ConfigCmd = &cobra.Command{
	Use:   "config",
	Short: "Generates a database.yml file for your project.",
	RunE: func(cmd *cobra.Command, args []string) error {
		cflag := cmd.Flag("config")
		cfgFile := defaults.String(cflag.Value.String(), "database.yml")
		pwd, err := os.Getwd()
		if err != nil {
			return errors.Wrap(err, "couldn't get the current directory")
		}
		data := map[string]interface{}{
			"dialect": dialect,
			"name":    filepath.Base(pwd),
		}
		return GenerateConfig(cfgFile, data)
	},
}

//GenerateConfig generates pop configuration files
func GenerateConfig(cfgFile string, data map[string]interface{}) error {
	pwd, _ := os.Getwd()
	if data["appPath"] == nil {
		data["appPath"] = pwd
	}
	if data["sqlitePath"] == nil {
		data["sqlitePath"] = pwd
	}

	dialect = strings.ToLower(data["dialect"].(string))
	if t, ok := configTemplates[dialect]; ok {
		g := makr.New()
		g.Add(makr.NewFile(cfgFile, t))
		return g.Run(".", data)
	}
	return errors.Errorf("Could not initialize %s!", dialect)
}
