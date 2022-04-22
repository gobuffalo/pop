package cmd

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/gobuffalo/pop/v6/fix"
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
			return fixFizz(path)
		})
	},
}

func fixFizz(path string) error {
	ext := strings.ToLower(filepath.Ext(path))
	if ext != ".fizz" {
		return nil
	}
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()
	return fix.Fizz(f, f)
}

func init() {
	RootCmd.AddCommand(fixCmd)
}
