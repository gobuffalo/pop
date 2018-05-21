package validate

import (
	"github.com/spf13/cobra"
	"path/filepath"
	"os"
	"strings"
	"github.com/pkg/errors"
)

var modelPath string
var models *[]string
var tags *[]string

func init() {
	ModelCmd.Flags().StringVarP(&modelPath, "model-path", "", "", "sets the path for your models")
	ModelCmd.Flags().StringArrayVarP(models, "models", "", []string{}, "sets models to be validated")
	ModelCmd.Flags().StringArrayVarP(tags, "tags", "", []string{}, "sets tags to be validated")
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
		v.AddDefaultProcessors(*tags...)
		errs, err := v.Run(*models...)

		if err != nil {
			panic(err)
		}

		if len(errs.Errors) > 0 {
			msgs := []string{}

			for _, structErrs := range errs.Errors {
				for _, err := range structErrs {
					msgs = append(msgs, err)
				}
			}

			return errors.New(strings.Join(msgs, "\n"))
		}
		return nil
	},
}
