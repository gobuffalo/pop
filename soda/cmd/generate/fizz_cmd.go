package generate

import (
	"context"

	"github.com/gobuffalo/genny"
	"github.com/gobuffalo/logger"
	"github.com/gobuffalo/pop/genny/fizz/cempty"
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
		run := genny.WetRunner(context.Background())

		// Ensure the generator is as verbose as the old one.
		lg := logger.New(logger.DebugLevel)
		run.Logger = lg

		p := cmd.Flag("path")
		path := ""
		if p != nil {
			path = p.Value.String()
		}

		g, err := cempty.New(&cempty.Options{
			TableName: name,
			Path:      path,
			Type:      "fizz",
		})
		if err != nil {
			return err
		}
		run.With(g)

		return run.Run()
	},
}
