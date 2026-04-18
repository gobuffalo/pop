package cmd

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_RootCmd_Environment(t *testing.T) {
	r := require.New(t)
	c := RootCmd

	// Fallback on default env
	c.SetArgs([]string{"help"})
	err := c.Execute()
	r.NoError(err)
	r.Equal("development", env)

	// Override with GO_ENV
	c.SetArgs([]string{"help"})
	t.Setenv("GO_ENV", "test")
	err = c.Execute()
	r.NoError(err)
	r.Equal("test", env)

	// CLI flag priority: the preferred order of flags and commands
	c.SetArgs([]string{
		"--env",
		"production",
		"help",
	})
	t.Setenv("GO_ENV", "test")
	err = c.Execute()
	r.NoError(err)
	r.Equal("production", env)

	// the following order works fine now but need to be considered again
	// CLI flag priority
	c.SetArgs([]string{
		"help",
		"--env",
		"production",
	})
	t.Setenv("GO_ENV", "test")
	err = c.Execute()
	r.NoError(err)
	r.Equal("production", env)
}
