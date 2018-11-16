package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/gobuffalo/pop"
	"github.com/markbates/going/defaults"
	"github.com/spf13/cobra"
)

var cfgFile string
var env string
var version bool

// RootCmd is the entry point of soda CLI.
var RootCmd = &cobra.Command{
	Short: "A tasty treat for all your database needs",
	PersistentPreRun: func(c *cobra.Command, args []string) {
		fmt.Printf("%s\n\n", Version)
		// CLI flag has priority
		if !c.PersistentFlags().Changed("env") {
			env = defaults.String(os.Getenv("GO_ENV"), env)
		}
		// TODO! Only do this when the command needs it.
		setConfigLocation()
		pop.LoadConfigFile()
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		if !version {
			return cmd.Help()
		}
		return nil
	},
}

// Execute runs RunCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		os.Exit(-1)
	}
}

func init() {
	RootCmd.Flags().BoolVarP(&version, "version", "v", false, "Show version information")
	RootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "The configuration file you would like to use.")
	RootCmd.PersistentFlags().StringVarP(&env, "env", "e", "development", "The environment you want to run migrations against. Will use $GO_ENV if set.")
	RootCmd.PersistentFlags().BoolVarP(&pop.Debug, "debug", "d", false, "Use debug/verbose mode")
}

func setConfigLocation() {
	if cfgFile != "" {
		abs, err := filepath.Abs(cfgFile)
		if err != nil {
			return
		}
		dir, file := filepath.Split(abs)
		pop.AddLookupPaths(dir)
		pop.ConfigName = file
	}
}

func getConn() *pop.Connection {
	conn := pop.Connections[env]
	if conn == nil {
		fmt.Printf("There is no connection named %s defined!\n", env)
		os.Exit(1)
	}
	return conn
}
