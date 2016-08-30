package pop

import (
	"bytes"
	"html/template"
	"io/ioutil"
	"path"

	"github.com/markbates/pop/fizz"
	"github.com/pkg/errors"
)

type migrationFile struct {
	Path      string
	FileName  string
	Version   string
	Name      string
	Direction string
	FileType  string
}

type migrationFiles []migrationFile

func (mfs migrationFiles) Len() int {
	return len(mfs)
}

func (mfs migrationFiles) Less(i, j int) bool {
	return mfs[i].Version < mfs[j].Version
}

func (mfs migrationFiles) Swap(i, j int) {
	mfs[i], mfs[j] = mfs[j], mfs[i]
}

func (m migrationFile) Content(c *Connection) (string, error) {
	b, err := ioutil.ReadFile(m.Path)
	if err != nil {
		return "", nil
	}
	content := string(b)
	ext := path.Ext(m.FileName)

	t := template.Must(template.New("sql").Parse(content))
	var bb bytes.Buffer
	err = t.Execute(&bb, c.Dialect.Details())
	if err != nil {
		return "", errors.Wrapf(err, "could not execute migration template %s", m.Path)
	}
	content = bb.String()

	if ext == ".fizz" {
		content, err = fizz.AString(content, c.Dialect.FizzTranslator())
		if err != nil {
			return "", errors.Wrapf(err, "could not fizz the migration %s", m.Path)
		}
	}
	return content, nil
}

func (m migrationFile) Execute(c *Connection) error {
	content, err := m.Content(c)
	if err != nil {
		return errors.Wrapf(err, "error processing %s", m.FileName)
	}

	if content == "" {
		return nil
	}

	err = c.RawQuery(content).Exec()
	if err != nil {
		return errors.Wrapf(err, "error executing %s", m.FileName)
	}
	return nil
}
