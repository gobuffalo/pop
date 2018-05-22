package cmd


import (
	"github.com/gobuffalo/pop/soda/cmd/validate"
	"github.com/spf13/cobra"
)

var validateCmd = &cobra.Command{
	Use:     "validate",
	Aliases: []string{"v"},
}

func init() {
	validateCmd.AddCommand(validate.ModelCmd)
	RootCmd.AddCommand(validateCmd)
}
