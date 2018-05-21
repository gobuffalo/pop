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
	"github.com/gobuffalo/envy"
	"sync"
)

// Tag represents a model struct tag.
type Tag struct {
	name string
	value string
	structName string
}

// GetName returns the name of the tag.
func (t *Tag) GetName() string  {
	return t.name
}

// GetValue returns the value of the tag.
func (t *Tag) GetValue() string  {
	return t.value
}

// GetStructName returns the struct name the tag belongs to.
func (t *Tag) GetStructName() string  {
	return t.structName
}

func getPackages(folder string, models ...string) map[string]*ast.Package {
	var path string

	for _, gopath := range envy.GoPaths() {
		path = filepath.Join(gopath, "src")
		path = filepath.Join(path, folder)

		if _, err := os.Stat(path); os.IsExist(err) {
			break
		}
	}

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

	concatNames := strings.Join(tagNames, "|")

	for _, name := range tagNames {
		if name == AllTags {
			concatNames = "[a-z0-9]"
			break
		}
	}

	var dbRegex = regexp.MustCompile(
		strings.Join([]string{
			"(",
			concatNames,
			")",
			`[ ]*:[ ]*"([^"]+)"`},
			"",
			),
		)

	tagChans := []<-chan *Tag{}
	tags := map[string][]*Tag{}

	for _, pkg := range packages {
		for _, file := range pkg.Files {
			tagChan := collecFields(file, dbRegex)
			tagChans = append(tagChans, tagChan)
		}
	}

	t := multiplex(tagChans...)

	 Loop: for {
		select {
			case tag, ok := <-t:
				if !ok {
					break Loop
				}

				tags[tag.structName] = append(tags[tag.structName], tag)
		}
	}

	return tags
}

func multiplex(cs ...<-chan *Tag) <-chan *Tag {
	var wg sync.WaitGroup
	out :=  make(chan *Tag, 50 * len(cs))

	output := func(c <- chan *Tag) {
		defer wg.Done()

		 for {
			select {
				case tag := <- c:
					if tag == nil {
						return
					}

					out <- tag
			}
		}
	}

	wg.Add(len(cs))

	for _, c := range cs {
		go output(c)
	}

	go func() {
		wg.Wait()
		close(out)
	}()

	return out
}

func collecFields(file *ast.File, dbRegex *regexp.Regexp) <-chan *Tag {

	tagChan := make(chan *Tag, 50)
	var structName string
	var nodesCnt int = 0
	var curNode int = 1

		go func() {
			ast.Inspect(file, func(node ast.Node) bool {
				nodesCnt++
				return true
			})


			ast.Inspect(file, func(node ast.Node) bool {
				switch x := node.(type) {
				case *ast.TypeSpec:
					//Get the struct name and end pos
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

				if  (curNode == nodesCnt) {
					tagChan <- nil
				}

				curNode++
				return true
			})
		}()

	return tagChan
}
