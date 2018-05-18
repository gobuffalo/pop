package generate

import (
	"regexp"

	"github.com/markbates/inflect"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var skipMigration bool
var structTag = "json"
var migrationPath string
var migrationType string
var env string

var nrx = regexp.MustCompile(`^nulls\.(.+)`)

func init() {
	ModelCmd.Flags().StringVarP(&structTag, "struct-tag", "", "json", "sets the struct tags for model (xml or json)")
	ModelCmd.Flags().StringVarP(&migrationType, "migration-type", "", "fizz", "sets the type of migration files for model (sql or fizz)")
	ModelCmd.Flags().BoolVarP(&skipMigration, "skip-migration", "s", false, "Skip creating a new fizz migration for this model.")
	ModelCmd.Flags().StringVarP(&migrationPath, "path", "p", "./migrations", "location of migrations folder")
	ModelCmd.Flags().StringVarP(&env, "env", "env", "", "environment")

	inflect.AddAcronym("ID")
	inflect.AddAcronym("UUID")
}

//ModelCmd is the cmd to generate a model
var ModelCmd = &cobra.Command{
	Use:     "model [name]",
	Aliases: []string{"m"},
	Short:   "Generates a model for your database",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return errors.New("You must supply a name for your model")
		}

		model, err := newModelFromArgs(args)

		if err != nil {
			return err
		}

		// Add a default UUID, if no custom ID is provided
		model.addID()

		err = model.generateModelFile()
		if err != nil {
			return err
		}

		if skipMigration {
			return nil
		}

		switch migrationType {
		case "sql":
			err = model.generateSQL(migrationPath, env)
		default:
			err = model.generateFizz(migrationPath)
		}

		if err != nil {
			return err
		}

		err = model.generateFizz(migrationPath)
		return err
	},
}
