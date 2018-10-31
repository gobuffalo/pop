package cmd

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_RootCmd_NoArg(t *testing.T) {
	// Fallback on default env
	r := require.New(t)
	c := RootCmd
	c.SetArgs([]string{})
	err := c.Execute()
	r.NoError(err)
	r.Equal("development", env)

	// Override with GO_ENV
	c.SetArgs([]string{})
	oldEnv := os.Getenv("GO_ENV")
	defer os.Setenv("GO_ENV", oldEnv)
	os.Setenv("GO_ENV", "test")
	err = c.Execute()
	r.NoError(err)
	r.Equal("test", env)

	// CLI flag priority
	c.SetArgs([]string{
		"--env",
		"production",
	})
	oldEnv = os.Getenv("GO_ENV")
	defer os.Setenv("GO_ENV", oldEnv)
	os.Setenv("GO_ENV", "test")
	err = c.Execute()
	r.NoError(err)
	r.Equal("production", env)
}
