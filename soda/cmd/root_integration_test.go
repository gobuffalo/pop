package cmd

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_RootCmd_Environtment(t *testing.T) {
	oldEnv := os.Getenv("GO_ENV")
	defer os.Setenv("GO_ENV", oldEnv)

	r := require.New(t)
	c := RootCmd

	// Fallback on default env
	c.SetArgs([]string{"help"})
	err := c.Execute()
	r.NoError(err)
	r.Equal("development", env)

	// Override with GO_ENV
	c.SetArgs([]string{"help"})
	os.Setenv("GO_ENV", "test")
	err = c.Execute()
	r.NoError(err)
	r.Equal("test", env)

	// CLI flag priority
	c.SetArgs([]string{
		"help",
		"--env",
		"production",
	})
	os.Setenv("GO_ENV", "test")
	err = c.Execute()
	r.NoError(err)
	r.Equal("production", env)
}
