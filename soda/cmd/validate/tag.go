package validate

import (
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"github.com/pkg/errors"
)

type Tag struct {
	name string
	value string
	structName string
}

func (t *Tag) getName() string  {
	return t.name
}

func (t *Tag) getValue() string  {
	return t.value
}

func (t *Tag) getStructName() string  {
	return t.structName
}

func getPackages(folder string, models ...string) map[string]*ast.Package {
	path := os.Getenv("GOPATH")
	path = filepath.Join(path, "src")

	path = filepath.Join(path, folder)

	fset := token.NewFileSet()
	modelMap := make(map[string]bool, len(models))

	for _, model := range models {
		k := strings.Join([]string{
			strings.ToLower(model),
			"go",
		}, ".")

		modelMap[k] = true
	}

	pkgs, err := parser.ParseDir(fset, path, func(f os.FileInfo) bool {
		isNotTest := !strings.HasSuffix(f.Name(), "_test.go")

		if len(modelMap) > 0 {
			_, exists := modelMap[strings.ToLower(f.Name())]
			//I want XOR, != is not natural
			return isNotTest != !exists
		}

		return isNotTest
	}, 0)

	if err != nil {
		panic(err)
	}

	if len(pkgs) == 0 {
		panic(errors.New("Could not find models package"))
	}

	return pkgs
}

func getTags(tagNames []string, packages map[string]*ast.Package) map[string][]*Tag {

	var dbRegex = regexp.MustCompile(
		strings.Join([]string{
			"(",
			strings.Join(tagNames, "|"),
			")",
			`[ ]*:[ ]*"([^"]+)"`},
			"",
			),
		)

	tagChans := []*chan *Tag{}
	tags := map[string][]*Tag{}

	for _, pkg := range packages {
		for _, file := range pkg.Files {
			tagChan := make(chan *Tag)
			tagChans = append(tagChans, &tagChan)

			go collecFields(tagChan, file, dbRegex)
		}
	}

	for _, ch := range tagChans {
	Loop:
		for {
			select {
			case tag, _ := <-*ch:

				if tag == nil {
					close(*ch)
					break Loop
				}

				tags[tag.structName] = append(tags[tag.structName], tag)
			}
		}
	}

	return tags
}

func collecFields(tagChan chan *Tag, file *ast.File, dbRegex *regexp.Regexp) {
	var structName string

	ast.Inspect(file, func(node ast.Node) bool {
		switch x := node.(type) {
		case *ast.TypeSpec:
			//Get the struct name
			structName = x.Name.Name
		case *ast.StructType:
			//Extract all db tags from the struct fields
			for _, field := range x.Fields.List {
				if field.Tag != nil {
					matches := dbRegex.FindAllStringSubmatch(field.Tag.Value, -1);
					if len(matches) > 0 {
						for _, matchTags := range matches {
							tagChan <- &Tag{
								matchTags[1],
								matchTags[2],
								structName,
							}
						}
					}
				}
			}
		}
		return true
	})

	//Tell the listener that we are done sending
	//values to this channel
	tagChan <- nil
}
