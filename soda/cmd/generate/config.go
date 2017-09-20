package generate

import (
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/gobuffalo/envy"
	"github.com/gobuffalo/makr"
	"github.com/markbates/going/defaults"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

func init() {
	ConfigCmd.Flags().StringVarP(&dialect, "type", "t", "postgres", "What type of database do you want to use? (postgres, mysql, sqlite3)")
}

var dialect string

var ConfigCmd = &cobra.Command{
	Use:   "config",
	Short: "Generates a database.yml file for your project.",
	RunE: func(cmd *cobra.Command, args []string) error {
		cflag := cmd.Flag("config")
		cfgFile := defaults.String(cflag.Value.String(), "database.yml")
		dir, err := os.Getwd()
		if err != nil {
			return errors.Wrap(err, "couldn't get the current directory")
		}
		pwd, _ := os.Getwd()
		data := map[string]interface{}{
			"dialect": dialect,
			"name":    path.Base(dir),
			"appPath": pwd,
		}
		return GenerateConfig(cfgFile, data)
	},
}

func goPath(root string) string {
	gpMultiple := envy.GoPaths()
	path := ""

	for i := 0; i < len(gpMultiple); i++ {
		if strings.HasPrefix(root, filepath.Join(gpMultiple[i], "src")) {
			path = gpMultiple[i]
			break
		}
	}
	return path
}

func packagePath(rootPath string) string {
	gosrcpath := strings.Replace(filepath.Join(goPath(rootPath), "src"), "\\", "/", -1)
	rootPath = strings.Replace(rootPath, "\\", "/", -1)
	return strings.Replace(rootPath, gosrcpath+"/", "", 2)
}

func GenerateConfig(cfgFile string, data map[string]interface{}) error {
	if data["appPath"] == nil {
		pwd, _ := os.Getwd()
		data["appPath"] = pwd
	}

	dialect = strings.ToLower(data["dialect"].(string))
	if t, ok := configTemplates[dialect]; ok {
		g := makr.New()
		g.Add(makr.NewFile(cfgFile, t))
		return g.Run(".", data)
	}
	return errors.Errorf("Could not initialize %s!", dialect)
}
