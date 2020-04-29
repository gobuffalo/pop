package fix

import (
	"io/ioutil"
	"regexp"
	"testing"

	"github.com/gobuffalo/packr/v2"
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

			re := regexp.MustCompile(`(?m)([\n\r])+$`)

			cleaned := re.ReplaceAllString(expected, "")
			cleanedPatched := re.ReplaceAllString(patched, "")

			rr.Equal(cleaned, cleanedPatched)

		})
		return nil
	})
	r.NoError(err)
}
