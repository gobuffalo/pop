package generate

import (
	"context"
	"errors"

	"github.com/gobuffalo/attrs"
	"github.com/gobuffalo/genny/v2"
	"github.com/gobuffalo/logger"
	"github.com/gobuffalo/pop/v6"
	"github.com/gobuffalo/pop/v6/genny/fizz/cempty"
	"github.com/gobuffalo/pop/v6/genny/fizz/ctable"
	"github.com/spf13/cobra"
)

// SQLCmd generates a SQL migration
var SQLCmd = &cobra.Command{
	Use:   "sql [name]",
	Short: "Generates Up/Down migrations for your database using SQL.",
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
		e := cmd.Flag("env")
		type nameable interface {
			Name() string
		}
		var translator nameable
		db, err := pop.Connect(e.Value.String())
		if err != nil {
			return err
		}
		t := db.Dialect.FizzTranslator()
		if tn, ok := t.(nameable); ok {
			translator = tn
		} else {
			return errors.New("invalid fizz translator")
		}

		if len(atts) == 0 {
			g, err := cempty.New(&cempty.Options{
				Name:       name,
				Path:       path,
				Type:       "sql",
				Translator: translator,
			})
			if err != nil {
				return err
			}
			run.With(g)
		} else {
			g, err := ctable.New(&ctable.Options{
				TableName:  name,
				Path:       path,
				Type:       "sql",
				Attrs:      atts,
				Translator: t,
			})
			if err != nil {
				return err
			}
			run.With(g)
		}

		return run.Run()
	},
}
