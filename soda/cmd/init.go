package cmd

import (
	"fmt"
	"html/template"
	"os"
	"path"

	"github.com/spf13/cobra"
)

var dialect string

var initTemplates = map[string]string{
	"postgres": `development:
  dialect: postgres
  database: {{.}}_development
  username: postgres
  password: postgres
  host: 127.0.0.1

test:
  dialect: postgres
  database: {{.}}_test
  username: postgres
  password: postgres
  host: 127.0.0.1

production:
  dialect: postgres
  database: {{.}}_production
  username: postgres
  password: postgres
  host: 127.0.0.1
`,
	"mysql": `development:
  dialect: "mysql"
  database: "{{.}}_development"
  host: "localhost"
  port: "3306"
  user: "root"
  password: "root"

test:
  dialect: "mysql"
  database: "{{.}}_test"
  host: "localhost"
  port: "3306"
  user: "root"
  password: "root"

production:
  dialect: "mysql"
  database: "{{.}}_production"
  host: "localhost"
  port: "3306"
  user: "root"
  password: "root"
	`,
	"sqlite3": `development:
	dialect: "sqlite3"
	database: "./{{.}}_development.sqlite"

test:
	dialect: "sqlite3"
	database: "./{{.}}_test.sqlite"

production:
	dialect: "sqlite3"
	database: "./{{.}}_production.sqlite"
`,
}

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Sets up a new Pop/Soda Project",
	RunE: func(cmd *cobra.Command, args []string) error {
		if t, ok := initTemplates[dialect]; ok {
			dir, err := os.Getwd()
			if err != nil {
				return err
			}
			f, err := os.Create(cfgFile)
			if err != nil {
				return err
			}
			tp := template.Must(template.New("database.yml").Parse(t))
			return tp.Execute(f, path.Base(dir))
		}
		return fmt.Errorf("Could not initialize %s!", dialect)
	},
}

func init() {
	RootCmd.AddCommand(initCmd)
	initCmd.Flags().StringVarP(&dialect, "type", "t", "postgres", "What type of database do you want to use? (postgres, mysql, sqlite3)")
}
