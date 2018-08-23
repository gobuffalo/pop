package fix

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/gobuffalo/plush"
)

// Anko converts old anko-form migrations to new plush ones.
func Anko(content string) (string, error) {
	bb := &bytes.Buffer{}

	lines := strings.Split(content, "\n")
	l := len(lines)

	for i := 0; i < l; i++ {
		line := lines[i]
		tl := strings.TrimSpace(line)
		if strings.HasPrefix(tl, "create_table") {
			// skip already converted create_table
			if strings.Contains(line, ", func(t) {") {
				// fix create_table
				line = strings.Replace(line, ", func(t) {", ") {", -1)
				ll := i
				lines[i] = line
				for {
					if strings.HasPrefix(tl, "})") {
						line = "}"
						break
					} else if strings.HasPrefix(tl, "}") {
						break
					}
					i++
					line = lines[i]
					tl = strings.TrimSpace(line)
					if l == i {
						return "", fmt.Errorf("unclosed create_table statement line %d", ll+1)
					}
				}
			}
		} else if strings.HasPrefix(tl, "raw(") {
			// fix raw
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
