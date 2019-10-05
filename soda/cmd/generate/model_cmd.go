package generate

import (
	"context"
	"errors"

	"github.com/gobuffalo/attrs"
	"github.com/gobuffalo/fizz"
	"github.com/gobuffalo/genny"
	"github.com/gobuffalo/logger"
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
	ModelCmd.Flags().StringVarP(&modelCmdConfig.StructTag, "struct-tag", "", "json", "sets the struct tags for model (xml/json/jsonapi)")
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

		var (
			atts attrs.Attrs
			err  error
		)
		if len(args) > 1 {
			atts, err = attrs.ParseArgs(args[1:]...)
			if err != nil {
				return err
			}
		}

		run := genny.WetRunner(context.Background())

		// Ensure the generator is as verbose as the old one.
		lg := logger.New(logger.DebugLevel)
		run.Logger = lg

		// Mount models generator
		g, err := gmodel.New(&gmodel.Options{
			Name:                   name,
			Attrs:                  atts,
			Path:                   modelCmdConfig.ModelPath,
			Encoding:               modelCmdConfig.StructTag,
			ForceDefaultID:         true,
			ForceDefaultTimestamps: true,
		})
		if err != nil {
			return err
		}

		run.With(g)

		// Mount migrations generator
		if !modelCmdConfig.SkipMigration {
			p := cmd.Flag("path")
			path := ""
			if p != nil {
				path = p.Value.String()
			}
			e := cmd.Flag("env")
			var translator fizz.Translator
			if modelCmdConfig.MigrationType == "sql" {
				db, err := pop.Connect(e.Value.String())
				if err != nil {
					return err
				}
				translator = db.Dialect.FizzTranslator()
			}

			g, err = ctable.New(&ctable.Options{
				TableName:              name,
				Attrs:                  atts,
				Path:                   path,
				Type:                   modelCmdConfig.MigrationType,
				Translator:             translator,
				ForceDefaultID:         true,
				ForceDefaultTimestamps: true,
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

	mt, found := opts["marshalType"].(string)
	if !found {
		return errors.New("marshalType option is required")
	}

	pp, found := opts["modelPath"].(string)
	if !found {
		return errors.New("modelPath option is required")
	}

	atts, err := attrs.ParseArgs(attributes...)
	if err != nil {
		return err
	}

	run := genny.WetRunner(context.Background())

	// Mount models generator
	g, err := gmodel.New(&gmodel.Options{
		Name:                   name,
		Attrs:                  atts,
		Path:                   pp,
		Encoding:               mt,
		ForceDefaultID:         true,
		ForceDefaultTimestamps: true,
	})
	if err != nil {
		return err
	}
	run.With(g)

	sm, found := opts["skipMigration"].(bool)
	// Mount migrations generator
	if found && sm {
		p, found := opts["path"].(string)
		if !found {
			return errors.New("path option is required")
		}
		migrationT, found := opts["migrationType"].(string)
		if !found {
			return errors.New("migrationType option is required")
		}
		var translator fizz.Translator
		if migrationT == "sql" {
			env, found := opts["env"].(string)
			if !found {
				return errors.New("env option is required")
			}
			db, err := pop.Connect(env)
			if err != nil {
				return err
			}
			translator = db.Dialect.FizzTranslator()
		}

		g, err = ctable.New(&ctable.Options{
			TableName:  name,
			Attrs:      atts,
			Path:       p,
			Type:       migrationT,
			Translator: translator,
		})
		if err != nil {
			return err
		}
		run.With(g)
	}
	return run.Run()
}
