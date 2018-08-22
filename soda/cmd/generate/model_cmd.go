package generate

import (
	"fmt"
	"strings"

	"github.com/markbates/inflect"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var skipMigration bool
var structTag string
var migrationType string

func init() {
	ModelCmd.Flags().StringVarP(&structTag, "struct-tag", "", "json", "sets the struct tags for model (xml or json)")
	ModelCmd.Flags().StringVarP(&migrationType, "migration-type", "", "fizz", "sets the type of migration files for model (sql or fizz)")
	ModelCmd.Flags().BoolVarP(&skipMigration, "skip-migration", "s", false, "Skip creating a new fizz migration for this model.")

	inflect.AddAcronym("ID")
	inflect.AddAcronym("UUID")
}

// ModelCmd is the cmd to generate a model
var ModelCmd = &cobra.Command{
	Use:     "model [name]",
	Aliases: []string{"m"},
	Short:   "Generates a model for your database",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return errors.New("you must supply a name for your model")
		}

		p := cmd.Flag("path")
		e := cmd.Flag("env")
		data := map[string]interface{}{
			"marshalType": structTag,
			"path":        p.Value.String(),
			"env":         e.Value.String(),
		}
		return Model(args[0], data, args[1:])
	},
}

// Model generates new model files to work with pop.
func Model(name string, opts map[string]interface{}, attributes []string) error {
	if strings.TrimSpace(name) == "" {
		return errors.New("model name can't be empty")
	}
	model := newModel(name)

	mt, found := opts["marshalType"].(string)
	if !found {
		return errors.New("marshalType option is required")
	}
	switch mt {
	case "json":
		model.Imports = append(model.Imports, "encoding/json")
	case "xml":
		model.Imports = append(model.Imports, "encoding/xml")
	default:
		return errors.New("invalid struct tags (use xml or json)")
	}

	attrs := make(map[inflect.Name]struct{})
	for _, def := range attributes {
		a := newAttribute(def, &model)
		if _, found := attrs[a.Name]; found {
			return fmt.Errorf("duplicated field \"%s\"", a.Name.String())
		}
		attrs[a.Name] = struct{}{}
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

	p, found := opts["path"].(string)
	if !found {
		return errors.New("path option is required")
	}
	switch migrationType {
	case "sql":
		env, found := opts["env"].(string)
		if !found {
			return errors.New("env option is required")
		}
		err = model.generateSQL(p, env)
	default:
		err = model.generateFizz(p)
	}
	return err
}
