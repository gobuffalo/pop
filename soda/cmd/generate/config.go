package generate

import (
	"fmt"
	"html/template"
	"os"
	"path"
	"strings"

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
		return GenerateConfig(dialect, cfgFile)
	},
}

func GenerateConfig(dialect string, cfgFile string) error {
	dialect = strings.ToLower(dialect)
	if t, ok := configTemplates[dialect]; ok {
		dir, err := os.Getwd()
		if err != nil {
			return errors.Wrap(err, "couldn't get the current directory")
		}
		err = os.MkdirAll(path.Dir(cfgFile), 0766)
		f, err := os.Create(cfgFile)
		if err != nil {
			return errors.Wrapf(err, "couldn't create the config file %s", cfgFile)
		}
		tp := template.Must(template.New("database.yml").Parse(t))

		dir = path.Base(dir)
		err = tp.Execute(f, dir)
		if err != nil {
			return errors.Wrap(err, "couldn't execute template")
		}
		fmt.Printf("Generated %s using the %s template.\n", cfgFile, dialect)
		return nil
	}
	return errors.Errorf("Could not initialize %s!", dialect)
}
