package validate

import (
	"regexp"
	"sort"
	"go/ast"
	"os"
	"path/filepath"
	"go/token"
	"go/parser"
	"strings"
)

func getPackages(folder string)  map[string]*ast.Package  {
	path, _ := os.Getwd()
	path = filepath.Join(folder)
	fset := token.NewFileSet()

	pkgs, err := parser.ParseDir(fset, path, func(f os.FileInfo) bool {
		return !strings.HasSuffix(f.Name(), "_test.go")
	}, 0)

	if err != nil {
		panic(err)
	}

	return pkgs
}

func getTags(tagName string, packages map[string]*ast.Package) map[string][]string {

	var dbRegex = regexp.MustCompile(tagName + `[ ]*:[ ]*"(.+)"`)

	tagChans := []*chan []string{}
	tags := map[string][]string{}

	for _, pkg := range packages {
		for _, file := range pkg.Files {
			tagChan := make(chan []string)
			tagChans = append(tagChans, &tagChan)

			go func(tag chan []string, file *ast.File, dbRegex *regexp.Regexp) {
				structs := make(map[int]string)
				keys := []int{}

				ast.Inspect(file, func(node ast.Node) bool {
					switch x := node.(type) {
					case *ast.TypeSpec:
						pos := int(x.Name.Pos())
						structs[pos] = x.Name.Name
						keys = append(keys, pos)

					case *ast.StructType:

						sort.SliceStable(keys, func(i, j int) bool {
							return keys[i] > keys[j]
						})
						//Get the struct name based on the current position
						currentPos := int(x.Pos())
						structName, exist := structs[currentPos]

						if !exist {
							for _, pos := range keys {
								if  currentPos >= pos {
									structName = structs[pos]
									break
								}
							}
						}

						//Extract all db tags from the struct fields
						for _, field := range x.Fields.List {
							if field.Tag != nil {
								if  matches := dbRegex.FindStringSubmatch(field.Tag.Value); len(matches) > 0 && len(matches[1]) > 0 {
									res := []string{structName, matches[1]}
									tag <- res
								}
							}
						}
					}
					return true
				})

				//remove struct names from proccessed file
				structs = map[int]string{}
				keys = keys[:0]
				tag <- nil

			}(tagChan, file, dbRegex)
		}
	}

	for _, ch := range tagChans {
	Loop: for   {
		select {
		case tag := <-*ch:
			if tag == nil {
				close(*ch)
				break Loop
			}

			tags[tag[0]] = append(tags[tag[0]], tag[1])
		}
	}
	}

	return tags
}
