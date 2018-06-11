package generate

import (
	"regexp"

	"github.com/markbates/inflect"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var skipMigration bool
var structTag string
var migrationType string

var nrx = regexp.MustCompile(`^nulls\.(.+)`)

func init() {
	ModelCmd.Flags().StringVarP(&structTag, "struct-tag", "", "json", "sets the struct tags for model (xml or json)")
	ModelCmd.Flags().StringVarP(&migrationType, "migration-type", "", "fizz", "sets the type of migration files for model (sql or fizz)")
	ModelCmd.Flags().BoolVarP(&skipMigration, "skip-migration", "s", false, "Skip creating a new fizz migration for this model.")

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

		model := newModel(args[0])

		switch structTag {
		case "json":
			model.Imports = append(model.Imports, "encoding/json")
		case "xml":
			model.Imports = append(model.Imports, "encoding/xml")
		default:
			return errors.New("Invalid struct tags (use xml or json)")
		}

		for _, def := range args[1:] {
			a := newAttribute(def, &model)
			model.addAttribute(a)
		}

		// Add a default UUID, if no custom ID is provided
		model.addID()

		err := model.generateModelFile()
		if err != nil {
			return err
		}

		if skipMigration {
			return nil
		}

		switch migrationType {
		case "sql":
			err = model.generateSQL(cmd.Flag("path"), cmd.Flag("env"))
		default:
			err = model.generateFizz(cmd.Flag("path"))
		}

		return err
	},
}
