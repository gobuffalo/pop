package generate

import (
	"context"

	"github.com/gobuffalo/attrs"
	"github.com/gobuffalo/genny/v2"
	"github.com/gobuffalo/logger"
	"github.com/gobuffalo/pop/v6/genny/fizz/cempty"
	"github.com/gobuffalo/pop/v6/genny/fizz/ctable"
	"github.com/spf13/cobra"
)

// FizzCmd generates a new fizz migration
var FizzCmd = &cobra.Command{
	Use:     "fizz [name]",
	Aliases: []string{"migration"},
	Short:   "Generates Up/Down migrations for your database using fizz.",
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

		p := cmd.Flag("path")
		path := ""
		if p != nil {
			path = p.Value.String()
		}

		if len(atts) == 0 {
			g, err := cempty.New(&cempty.Options{
				Name: name,
				Path: path,
				Type: "fizz",
			})
			if err != nil {
				return err
			}
			run.With(g)
		} else {
			g, err := ctable.New(&ctable.Options{
				TableName: name,
				Path:      path,
				Type:      "fizz",
				Attrs:     atts,
			})
			if err != nil {
				return err
			}
			run.With(g)
		}

		return run.Run()
	},
}
