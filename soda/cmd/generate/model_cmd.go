package generate

import (
	"context"
	"os"
	"os/exec"

	"github.com/gobuffalo/attrs"
	"github.com/gobuffalo/fizz"
	"github.com/gobuffalo/genny/v2"
	"github.com/gobuffalo/genny/v2/gogen"
	"github.com/gobuffalo/logger"
	"github.com/gobuffalo/pop/v6"
	"github.com/gobuffalo/pop/v6/genny/fizz/ctable"
	gmodel "github.com/gobuffalo/pop/v6/genny/model"
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

		// format generated go files
		pwd, _ := os.Getwd()
		g, err = gogen.Fmt(pwd)
		if err != nil {
			return err
		}
		run.With(g)

		// generated modules may have new dependencies
		if _, err := os.Stat("go.mod"); err == nil {
			g = genny.New()
			g.Command(exec.Command("go", "mod", "tidy"))
			run.With(g)
		}

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
