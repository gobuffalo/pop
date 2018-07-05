package fix

import (
	"io/ioutil"
	"strings"
	"testing"

	"github.com/gobuffalo/packr"
	"github.com/stretchr/testify/require"
)

func Test_Anko(t *testing.T) {
	r := require.New(t)
	box := packr.NewBox("./fixtures")
	err := box.Walk(func(path string, info packr.File) error {
		if strings.HasPrefix(path, "pass") {
			t.Run(path, testPass(path, info))
			return nil
		}
		t.Run(path, testFail(path, info))
		return nil
	})
	r.NoError(err)
}

func testPass(path string, info packr.File) func(*testing.T) {
	return func(t *testing.T) {
		r := require.New(t)
		b, err := ioutil.ReadAll(info)
		r.NoError(err)

		body := string(b)
		fixed, err := Anko(body)
		r.NoError(err)
		if strings.Contains(path, "anko") {
			r.NotEqual(body, fixed)
		} else {
			r.Equal(body, fixed)
		}
	}
}

func testFail(path string, info packr.File) func(*testing.T) {
	return func(t *testing.T) {
		r := require.New(t)
		b, err := ioutil.ReadAll(info)
		r.NoError(err)

		body := string(b)
		_, err = Anko(body)
		r.Error(err)
	}
}
