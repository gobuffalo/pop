package validate

import (
	"github.com/spf13/cobra"
	"github.com/pkg/errors"
	"strings"
)

//ModelCmd is the cmd to generate a model
var ModelCmd = &cobra.Command{
	Use:     "model",
	Aliases: []string{"m"},
	Short:   "Validates model db fields",
	RunE: func(cmd *cobra.Command, args []string) error {
		m := NewModel()
		errs := m.Validate()

		if len(errs) > 0 {
			msgs := []string{}

			for _, err := range errs {
				msgs = append(msgs, err.Error())
			}

			return errors.New(strings.Join(msgs, "\n"))
		}
		return nil
	},
}

func init() {

}
