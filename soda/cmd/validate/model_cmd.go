package validate

import (
	"strings"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"path/filepath"
	"os"
)

var modelPath string

func init() {
	ModelCmd.Flags().StringVarP(&modelPath, "model-path", "", "", "sets the path for your models")
}

//ModelCmd is the cmd to validate models
var ModelCmd = &cobra.Command{
	Use:     "model [model]",
	Aliases: []string{"m"},
	Short:   "Validates model db fields",
	RunE: func(cmd *cobra.Command, args []string) error {
		path := modelPath

		if len(modelPath) == 0 {
			base := os.Getenv("GOPATH")
			base = filepath.Join(base, "src")

			rel, _ := os.Getwd()

			relPath, err := filepath.Rel(base, rel)

			if err != nil {
				panic(err)
			}

			path = relPath
		}


		v := NewValidator(path)
		v.AddDefaultProcessors()
		errs, err := v.Run()

		if err != nil {
			panic(err)
		}

		if len(errs) > 0 {
			msgs := []string{}

			for _, tagErrs := range errs {
				for _, err := range tagErrs {
					msgs = append(msgs, err.Error())
				}
			}

			return errors.New(strings.Join(msgs, "\n"))
		}
		return nil
	},
}
