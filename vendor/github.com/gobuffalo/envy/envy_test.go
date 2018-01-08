package envy_test

import (
	"os"
	"testing"

	"github.com/gobuffalo/envy"
	"github.com/stretchr/testify/require"
)

func Test_Get(t *testing.T) {
	r := require.New(t)
	r.NotZero(os.Getenv("GOPATH"))
	r.Equal(os.Getenv("GOPATH"), envy.Get("GOPATH", "foo"))
	r.Equal("bar", envy.Get("IDONTEXIST", "bar"))
}

func Test_MustGet(t *testing.T) {
	r := require.New(t)
	r.NotZero(os.Getenv("GOPATH"))
	v, err := envy.MustGet("GOPATH")
	r.NoError(err)
	r.Equal(os.Getenv("GOPATH"), v)

	_, err = envy.MustGet("IDONTEXIST")
	r.Error(err)
}

func Test_Set(t *testing.T) {
	r := require.New(t)
	_, err := envy.MustGet("FOO")
	r.Error(err)

	envy.Set("FOO", "foo")
	r.Equal("foo", envy.Get("FOO", "bar"))
}

func Test_MustSet(t *testing.T) {
	r := require.New(t)

	r.Zero(os.Getenv("FOO"))

	err := envy.MustSet("FOO", "BAR")
	r.NoError(err)

	r.Equal("BAR", os.Getenv("FOO"))
}

func Test_Temp(t *testing.T) {
	r := require.New(t)

	_, err := envy.MustGet("BAR")
	r.Error(err)

	envy.Temp(func() {
		envy.Set("BAR", "foo")
		r.Equal("foo", envy.Get("BAR", "bar"))
		_, err = envy.MustGet("BAR")
		r.NoError(err)
	})

	_, err = envy.MustGet("BAR")
	r.Error(err)
}

func Test_GoPath(t *testing.T) {
	r := require.New(t)
	envy.Temp(func() {
		envy.Set("GOPATH", "/foo")
		r.Equal("/foo", envy.GoPath())
	})
}

func Test_GoPaths(t *testing.T) {
	r := require.New(t)
	envy.Temp(func() {
		envy.Set("GOPATH", "/foo:/bar")
		r.Equal([]string{"/foo", "/bar"}, envy.GoPaths())
	})
}

func Test_CurrentPackage(t *testing.T) {
	r := require.New(t)
	r.Equal("github.com/gobuffalo/envy", envy.CurrentPackage())
}

// Env files loading

func Test_LoadDefaultEnvWhenNoArgsPassed(t *testing.T) {
	r := require.New(t)
	envy.Temp(func() {
		err := envy.Load()
		r.NoError(err)

		r.Equal("root", envy.Get("DIR", ""))
		r.Equal("none", envy.Get("FLAVOUR", ""))
		r.Equal("false", envy.Get("INSIDE_FOLDER", ""))
	})
}

func Test_DoNotLoadDefaultEnvWhenArgsPassed(t *testing.T) {
	r := require.New(t)
	envy.Temp(func() {
		err := envy.Load("test_env/.env")
		r.NoError(err)

		r.Equal("test_env", envy.Get("DIR", ""))
		r.Equal("none", envy.Get("FLAVOUR", ""))
		r.Equal("true", envy.Get("INSIDE_FOLDER", ""))
	})
}

func Test_OverloadParams(t *testing.T) {
	r := require.New(t)
	envy.Temp(func() {
		err := envy.Load("test_env/.env.test", "test_env/.env.prod")
		r.NoError(err)

		r.Equal("production", envy.Get("FLAVOUR", ""))
	})
}

func Test_ErrorWhenSingleFileLoadDoesNotExist(t *testing.T) {
	r := require.New(t)
	envy.Temp(func() {
		err := envy.Load(".env.fake")
		r.Error(err)

		r.Equal("FAILED", envy.Get("FLAVOUR", "FAILED"))
	})
}

func Test_KeepEnvWhenFileInListFails(t *testing.T) {
	r := require.New(t)
	envy.Temp(func() {
		err := envy.Load(".env", ".env.FAKE")
		r.Error(err)
		r.Equal("none", envy.Get("FLAVOUR", "FAILED"))
		r.Equal("root", envy.Get("DIR", "FAILED"))
	})
}

func Test_KeepEnvWhenSecondLoadFails(t *testing.T) {
	r := require.New(t)
	envy.Temp(func() {
		err := envy.Load(".env")
		r.NoError(err)
		r.Equal("none", envy.Get("FLAVOUR", "FAILED"))
		r.Equal("root", envy.Get("DIR", "FAILED"))

		err = envy.Load(".env.FAKE")

		r.Equal("none", envy.Get("FLAVOUR", "FAILED"))
		r.Equal("root", envy.Get("DIR", "FAILED"))
	})
}

func Test_StopLoadingWhenFileInListFails(t *testing.T) {
	r := require.New(t)
	envy.Temp(func() {
		err := envy.Load(".env", ".env.FAKE", "test_env/.env.prod")
		r.Error(err)
		r.Equal("none", envy.Get("FLAVOUR", "FAILED"))
		r.Equal("root", envy.Get("DIR", "FAILED"))
	})
}
