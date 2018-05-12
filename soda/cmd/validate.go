package cmd

import (
	"github.com/spf13/cobra"
	"github.com/petar/pop/soda/cmd/validate"
)

var validateCmd = &cobra.Command{
	Use:     "validate",
	Aliases: []string{"v"},
}

func init() {
	validateCmd.AddCommand(validate.ModelCmd)
	RootCmd.AddCommand(validateCmd)
}
