package generate

import (
	"context"
	"strings"

	"github.com/gobuffalo/fizz"

	"errors"

	"github.com/gobuffalo/attrs"
	"github.com/gobuffalo/genny"
	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/pop/genny/fizz/ctable"
	gmodel "github.com/gobuffalo/pop/genny/model"
	"github.com/gobuffalo/pop/internal/oncer"
	"github.com/spf13/cobra"
)

var modelCmdConfig struct {
	SkipMigration bool
	StructTag     string
	MigrationType string
	ModelPath     string
}

func init() {
	ModelCmd.Flags().StringVarP(&modelCmdConfig.StructTag, "struct-tag", "", "json", "sets the struct tags for model (xml or json)")
	ModelCmd.Flags().StringVarP(&modelCmdConfig.MigrationType, "migration-type", "", "fizz", "sets the type of migration files for model (sql or fizz)")
	ModelCmd.Flags().BoolVarP(&modelCmdConfig.SkipMigration, "skip-migration", "s", false, "Skip creating a new fizz migration for this model.")
	ModelCmd.Flags().StringVarP(&modelCmdConfig.ModelPath, "models-path", "", "models", "the path the model will be created in")
}

// ModelCmd is the cmd to generate a model
var ModelCmd = &cobra.Command{
	Use:     "model [name]",
	Aliases: []string{"m"},
	Short:   "Generates a model for your database",
	RunE: func(cmd *cobra.Command, args []string) error {
		name := ""
		if len(args) > 0 {
			name = args[0]
		}

		atts, err := attrs.ParseArgs(args[1:]...)
		if err != nil {
			return err
		}

		run := genny.WetRunner(context.Background())

		// Mount models generator
		g, err := gmodel.New(gmodel.Options{
			Name:           name,
			Attrs:          atts,
			Path:           modelCmdConfig.ModelPath,
			Encoding:       modelCmdConfig.StructTag,
			ForceDefaultID: true,
		})
		if err != nil {
			return err
		}

		run.With(g)

		// Mount migrations generator
		if !modelCmdConfig.SkipMigration {
			p := cmd.Flag("path")
			e := cmd.Flag("env")
			var translator fizz.Translator
			if modelCmdConfig.MigrationType == "sql" {
				db, err := pop.Connect(e.Value.String())
				if err != nil {
					return err
				}
				translator = db.Dialect.FizzTranslator()
			}

			g, err = ctable.New(ctable.Options{
				TableName:      name,
				Attrs:          atts,
				Path:           p.Value.String(),
				Type:           modelCmdConfig.MigrationType,
				FizzTranslator: translator,
			})
			if err != nil {
				return err
			}
			run.With(g)
		}

		return run.Run()
	},
}

// Model generates new model files to work with pop.
func Model(name string, opts map[string]interface{}, attributes []string) error {
	oncer.Deprecate(0, "generate.Model", "Use github.com/gobuffalo/pop/genny/model instead.")
	if strings.TrimSpace(name) == "" {
		return errors.New("model name can't be empty")
	}
	mt, found := opts["marshalType"].(string)
	if !found {
		return errors.New("marshalType option is required")
	}

	pp, found := opts["modelPath"].(string)
	if !found {
		return errors.New("modelPath option is required")
	}

	model, err := newModel(name, mt, pp)
	if err != nil {
		return err
	}

	for _, def := range attributes {
		a, err := newAttribute(def, &model)
		if err != nil {
			return err
		}
		if err := model.addAttribute(a); err != nil {
			return err
		}
	}

	// Add a default UUID, if no custom ID is provided
	model.addID()

	if err := model.generateModelFile(); err != nil {
		return err
	}

	sm, found := opts["skipMigration"].(bool)
	if found && sm {
		return nil
	}

	p, found := opts["path"].(string)
	if !found {
		return errors.New("path option is required")
	}

	migrationT, found := opts["migrationType"].(string)
	if !found {
		return errors.New("migrationType option is required")
	}
	switch migrationT {
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
