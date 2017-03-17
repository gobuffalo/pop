package schema

import (
	"fmt"

	"github.com/spf13/cobra"
)

// schema/loadCmd represents the schema/load command
var LoadCmd = &cobra.Command{
	Use:   "load",
	Short: "A brief description of your command",
	Run: func(cmd *cobra.Command, args []string) {
		// TODO: Work your own magic here
		fmt.Println("schema/load called")
	},
}
