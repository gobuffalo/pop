package model

import (
	"path"
	"sort"
	"strings"
)

func buildImports(opts *Options) []string {
	imps := map[string]bool{
		"github.com/gobuffalo/validate/v3": true,
		"github.com/gobuffalo/pop/v6":      true,
	}
	if opts.Encoding == "jsonapi" {
		imps["github.com/google/jsonapi"] = true
		imps["strings"] = true
	} else {
		imps[path.Join("encoding", strings.ToLower(opts.Encoding))] = true
	}
	ats := opts.Attrs
	for _, a := range ats {
		switch a.GoType() {
		case "uuid", "uuid.UUID":
			imps["github.com/gofrs/uuid"] = true
		case "time.Time":
			imps["time"] = true
		default:
			if strings.HasPrefix(a.GoType(), "nulls") {
				imps["github.com/gobuffalo/nulls"] = true
			}
			if strings.HasPrefix(a.GoType(), "slices") {
				imps["github.com/gobuffalo/pop/v6/slices"] = true
			}
		}
	}
	i := make([]string, 0, len(imps))
	for k := range imps {
		i = append(i, k)
	}
	sort.Strings(i)
	return i
}
