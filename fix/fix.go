package fix

import (
	"io"
	"io/ioutil"
	"strings"
)

// Fizz fixes a fizz file to use the most up to date format.
// It takes the original contents from the Reader, and writes the fixed contents in the Writer.
func Fizz(r io.Reader, w io.Writer) error {
	b, err := ioutil.ReadAll(r)
	if err != nil {
		return err
	}

	content := string(b)

	// Old anko format
	fixed, err := Anko(content)
	if err != nil {
		return err
	}
	if strings.TrimSpace(fixed) != strings.TrimSpace(content) {
		content = fixed
	}

	// Rewrite migrations to use t.Timestamps() if necessary
	fixed, err = AutoTimestampsOff(content)
	if err != nil {
		return err
	}

	if strings.TrimSpace(fixed) != strings.TrimSpace(content) {
		if _, err := w.Write([]byte(fixed)); err != nil {
			return err
		}
	}
	return nil
}
