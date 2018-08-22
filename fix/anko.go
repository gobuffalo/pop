package fix

import (
	"bytes"
	"strings"

	"github.com/gobuffalo/plush"
)

// Anko converts old anko-form migrations to new plush ones.
func Anko(content string) (string, error) {
	bb := &bytes.Buffer{}

	lines := strings.Split(content, "\n")

	// fix create_table
	inCreateTable := false
	for i := 0; i < len(lines); i++ {
		line := lines[i]
		tl := strings.TrimSpace(line)
		if strings.HasPrefix(tl, "create_table") {
			line = strings.Replace(line, ", func(t) {", ") {", -1)
			inCreateTable = true
		}
		if strings.HasPrefix(tl, "}") && inCreateTable {
			inCreateTable = false
		}
		if strings.HasPrefix(tl, "})") && inCreateTable {
			line = "}"
			inCreateTable = false
		}
		lines[i] = line
	}

	// fix raw
	for i, line := range lines {
		tl := strings.TrimSpace(line)
		if strings.HasPrefix(tl, "raw(") {
			line = strings.Replace(line, "raw(", "sql(", -1)
		}
		lines[i] = line
	}

	body := strings.Join(lines, "\n")

	if _, err := plush.Parse(body); err != nil {
		return "", err
	}

	bb.WriteString(body)

	return bb.String(), nil
}
