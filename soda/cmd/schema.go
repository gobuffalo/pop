package cmd

import (
	"fmt"

	"github.com/markbates/pop/soda/cmd/schema"
	"github.com/spf13/cobra"
)

// schemaCmd represents the schema command
var schemaCmd = &cobra.Command{
	Use:   "schema",
	Short: "A brief description of your command",
	Run: func(cmd *cobra.Command, args []string) {
		// TODO: Work your own magic here
		fmt.Println("schema called")
	},
}

func init() {
	schemaCmd.AddCommand(schema.LoadCmd)
	schemaCmd.AddCommand(schema.DumpCmd)
	RootCmd.AddCommand(schemaCmd)
}
