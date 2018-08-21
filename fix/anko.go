package fix

import (
	"bytes"
	"strings"
)

// Anko converts old anko-form migrations to new plush ones.
func Anko(content string) (string, error) {
	bb := &bytes.Buffer{}

	lines := strings.Split(content, "\n")

	// fix create_table
	in_create_table := false
	for i := 0; i < len(lines); i++ {
		line := lines[i]
		tl := strings.TrimSpace(line)
		if strings.HasPrefix(tl, "create_table") {
			line = strings.Replace(line, ", func(t) {", ") {", -1)
			in_create_table = true
		}
		if strings.HasPrefix(tl, "}") && in_create_table {
			in_create_table = false
		}
		if strings.HasPrefix(tl, "})") && in_create_table {
			line = "}"
			in_create_table = false
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

	bb.WriteString(strings.Join(lines, "\n"))

	return bb.String(), nil
}
