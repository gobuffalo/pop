package fix

import (
	"io/fs"
	"os"
	"regexp"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_AutoTimestampsOff(t *testing.T) {
	r := require.New(t)
	box := os.DirFS("fixtures/auto_timestamps_off/raw")
	boxPatched := os.DirFS("fixtures/auto_timestamps_off/patched")

	err := fs.WalkDir(box, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			return nil
		}

		t.Run(path, func(tt *testing.T) {
			rr := require.New(tt)
			b, err := fs.ReadFile(box, path)
			rr.NoError(err)

			body := string(b)
			patched, err := AutoTimestampsOff(body)
			rr.NoError(err)
			expected, err := fs.ReadFile(boxPatched, path)
			rr.NoError(err)

			re := regexp.MustCompile(`(?m)([\n\r])+$`)

			cleaned := re.ReplaceAllString(string(expected), "")
			cleanedPatched := re.ReplaceAllString(patched, "")

			rr.Equal(cleaned, cleanedPatched)
		})
		return nil
	})

	r.NoError(err)
}
