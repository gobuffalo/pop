package pop

import (
	"bytes"
	"io"
	"io/ioutil"
	"text/template"

	"github.com/gobuffalo/fizz"
	"github.com/pkg/errors"
)

// MigrationContent returns the content of a migration.
func MigrationContent(mf Migration, c *Connection, r io.Reader, usingTemplate bool) (string, error) {
	b, err := ioutil.ReadAll(r)
	if err != nil {
		return "", nil
	}

	content := ""
	if usingTemplate {
		t := template.Must(template.New("migration").Parse(string(b)))
		var bb bytes.Buffer
		err = t.Execute(&bb, c.Dialect.Details())
		if err != nil {
			return "", errors.Wrapf(err, "could not execute migration template %s", mf.Path)
		}
		content = bb.String()
	} else {
		content = string(b)
	}

	if mf.Type == "fizz" {
		content, err = fizz.AString(content, c.Dialect.FizzTranslator())
		if err != nil {
			return "", errors.Wrapf(err, "could not fizz the migration %s", mf.Path)
		}
	}

	return content, nil
}
