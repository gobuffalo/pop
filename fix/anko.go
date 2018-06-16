package fix

import (
	"bytes"
	"strings"
)

func Anko(content string) (string, error) {
	if !strings.Contains(content, "create_table") {
		return content, nil
	}
	if !strings.Contains(content, "func(t) {") {
		return content, nil
	}

	bb := &bytes.Buffer{}

	for _, line := range strings.Split(content, "\n") {
		tl := strings.TrimSpace(line)
		if strings.HasPrefix(tl, "create_table") {
			line = strings.Replace(line, ", func(t) {", ") {", -1)
		}
		if strings.HasPrefix(tl, "})") {
			line = "}"
		}
		if tl == "" {
			continue
		}
		bb.WriteString(line + "\n")
	}

	return bb.String(), nil
}
