package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/gobuffalo/pop/v6"
	"github.com/gobuffalo/pop/v6/internal/defaults"
	"github.com/spf13/cobra"
)

var cfgFile string
var env string
var version bool

// RootCmd is the entry point of soda CLI.
var RootCmd = &cobra.Command{
	SilenceUsage: true,
	Short:        "A tasty treat for all your database needs",
	PersistentPreRun: func(c *cobra.Command, args []string) {
		fmt.Printf("pop %s\n\n", Version)

		/* NOTE: Do not use c.PersistentFlags. `c` is not always the
		RootCmd. The naming is confusing. The meaning of "persistent"
		in the `PersistentPreRun` is something like "this function will
		be 'sticky' to all subcommands and will run for them. So `c`
		can be any subcommands.
		However, the meaning of "persistent" in the `PersistentFlags`
		is, as the function comment said, "persistent FlagSet
		specifically set in the **current command**" so it is sticky
		to specific command!

		Use c.Flags() or c.Root().Flags() here.
		*/

		// CLI flag has priority
		if !c.Flags().Changed("env") {
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
