package fix

import (
	"bytes"
	"fmt"
	"regexp"
	"strings"

	"github.com/gobuffalo/plush/v4"
)

// Anko converts old anko-form migrations to new plush ones.
func Anko(content string) (string, error) {
	bb := &bytes.Buffer{}

	lines := strings.Split(content, "\n")
	l := len(lines)
	fre := regexp.MustCompile(`,\s*func\(t\)\s*{`)

	for i := 0; i < l; i++ {
		line := lines[i]
		tl := strings.TrimSpace(line)
		switch {
		case strings.HasPrefix(tl, "create_table") && fre.MatchString(line):
			// fix create_table
			line = fre.ReplaceAllString(line, ") {")
			ll := i
			lines[i] = line
			waitParen := false
			for {
				if strings.HasPrefix(tl, "})") {
					line = "}" + tl[2:]
					break
				} else if strings.HasPrefix(tl, "}") {
					// Now, we have to make sure to match the missing ")"
					waitParen = true
				} else if waitParen && strings.HasPrefix(tl, ")") {
					line = tl[1:]
					break
				}
				i++
				if l == i {
					return "", fmt.Errorf("unclosed create_table statement line %d", ll+1)
				}
				line = lines[i]
				tl = strings.TrimSpace(line)
			}
		case strings.HasPrefix(tl, "raw("):
			// fix raw
			line = strings.ReplaceAll(line, "raw(", "sql(")
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
