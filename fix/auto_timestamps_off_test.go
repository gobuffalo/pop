package fix

import (
	"io/ioutil"
	"strings"
	"testing"

	packr "github.com/gobuffalo/packr/v2"
	"github.com/stretchr/testify/require"
)

func Test_AutoTimestampsOff(t *testing.T) {
	r := require.New(t)
	box := packr.New("./fixtures/auto_timestamps_off/raw", "./fixtures/auto_timestamps_off/raw")
	boxPatched := packr.New("./fixtures/auto_timestamps_off/patched", "./fixtures/auto_timestamps_off/patched")

	err := box.Walk(func(path string, info packr.File) error {
		t.Run(path, func(tt *testing.T) {
			rr := require.New(tt)
			b, err := ioutil.ReadAll(info)
			rr.NoError(err)

			body := string(b)
			patched, err := AutoTimestampsOff(body)
			rr.NoError(err)
			expected, err := boxPatched.FindString(path)
			rr.NoError(err)
			rr.Equal(strings.Replace(expected, "\r", "", -1), patched)
		})
		return nil
	})
	r.NoError(err)
}
