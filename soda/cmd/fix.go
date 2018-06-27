package cmd

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/gobuffalo/pop/fix"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var fixCmd = &cobra.Command{
	Use:     "fix",
	Aliases: []string{"f", "update"},
	Short:   "Brings pop, soda, and fizz files in line with the latest APIs",
	RunE: func(cmd *cobra.Command, args []string) error {
		return filepath.Walk(migrationPath, func(path string, info os.FileInfo, _ error) error {
			if info == nil {
				return nil
			}
			ext := strings.ToLower(filepath.Ext(path))
			if ext != ".fizz" {
				return nil
			}

			b, err := ioutil.ReadFile(path)
			if err != nil {
				return errors.WithStack(err)
			}

			content := string(b)

			fixed, err := fix.Anko(content)
			if err != nil {
				return errors.WithStack(err)
			}
			if strings.TrimSpace(fixed) != strings.TrimSpace(content) {
				f, err := os.Create(path)
				if err != nil {
					return errors.WithStack(err)
				}
				if _, err := f.WriteString(fixed); err != nil {
					return errors.WithStack(err)
				}
				if err := f.Close(); err != nil {
					return errors.WithStack(err)
				}
			}

			return nil
		})
	},
}

func init() {
	RootCmd.AddCommand(fixCmd)
}
