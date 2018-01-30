package generate

import (
	"regexp"
	"strings"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var skipMigration bool
var structTag string

var nrx = regexp.MustCompile(`^nulls\.(.+)`)

func init() {
	ModelCmd.Flags().StringVarP(&structTag, "struct-tag", "", "json", "sets the struct tags for model (xml or json)")
	ModelCmd.Flags().BoolVarP(&skipMigration, "skip-migration", "s", false, "Skip creating a new fizz migration for this model.")
}

//ModelCmd is the cmd to generate a model
var ModelCmd = &cobra.Command{
	Use:     "model [name]",
	Aliases: []string{"m"},
	Short:   "Generates a model for your database",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return errors.New("you must supply a name for your model")
		}

		model := newModel(args[0])

		switch structTag {
		case "json":
			model.Imports = append(model.Imports, "encoding/json")
		case "xml":
			model.Imports = append(model.Imports, "encoding/xml")
		default:
			return errors.New("invalid struct tags (use xml or json)")
		}

		for _, def := range args[1:] {
			a := newAttribute(def, &model)
			model.AppendAttribute(a)
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

		err = model.generateFizz(cmd.Flag("path"))
		if err != nil {
			return err
		}

		return nil
	},
}

func colType(s string) string {
	switch strings.ToLower(s) {
	case "text":
		return "string"
	case "time", "timestamp", "datetime":
		return "time.Time"
	case "nulls.text":
		return "nulls.String"
	case "uuid":
		return "uuid.UUID"
	case "json", "jsonb":
		return "slices.Map"
	case "[]string":
		return "slices.String"
	case "[]int":
		return "slices.Int"
	case "slices.float", "[]float", "[]float32", "[]float64":
		return "slices.Float"
	default:
		return s
	}
}

func fizzColType(s string) string {
	switch strings.ToLower(s) {
	case "int":
		return "integer"
	case "time", "datetime":
		return "timestamp"
	case "uuid.uuid", "uuid":
		return "uuid"
	case "nulls.float32", "nulls.float64":
		return "float"
	case "slices.string", "slices.uuid", "[]string":
		return "varchar[]"
	case "slices.float", "[]float", "[]float32", "[]float64":
		return "numeric[]"
	case "slices.int":
		return "int[]"
	case "slices.map":
		return "jsonb"
	default:
		if nrx.MatchString(s) {
			return fizzColType(strings.Replace(s, "nulls.", "", -1))
		}
		return strings.ToLower(s)
	}
}
