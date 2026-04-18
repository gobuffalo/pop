package fix

import (
	"io"
	"io/fs"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_Anko(t *testing.T) {
	r := require.New(t)
	fsys := os.DirFS("fixtures/anko")
	err := fs.WalkDir(fsys, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			return nil
		}

		f, err := fsys.Open(path)
		if err != nil {
			return err
		}

		if strings.HasPrefix(path, "pass") {
			t.Run(path, func(t *testing.T) {
				r := require.New(t)
				b, err := io.ReadAll(f)
				r.NoError(err)

				body := string(b)
				fixed, err := Anko(body)
				r.NoError(err)

				if strings.Contains(path, "anko") {
					r.NotEqual(body, fixed)
				} else {
					r.Equal(body, fixed)
				}
			})
			return nil
		}
		t.Run(path, func(t *testing.T) {
			r := require.New(t)
			b, err := io.ReadAll(f)
			r.NoError(err)

			body := string(b)
			_, err = Anko(body)
			r.Error(err)
		})
		return nil
	})
	r.NoError(err)
}
