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

	var dbRegex = regexp.MustCompile(tagName + `:"([a-z_0-9, ]+)"`)
	var keys []int
	chanLimit := 5000
	tagChan := make(chan []string, chanLimit)
	tags := map[string][]string{}

	for _, pkg := range packages {
		for _, file := range pkg.Files {

			go func(tag chan []string) {
				structs := make(map[int]string)

				for _, obj := range file.Scope.Objects {
					structs[int(obj.Pos())] = obj.Name
					keys = append(keys, int(obj.Pos()))
				}

				sort.Ints(keys)

				ast.Inspect(file, func(node ast.Node) bool {
					switch x := node.(type) {
					case *ast.StructType:

						//Get the struct name based on the current position
						currentPos := int(x.Pos())
						structName, exist := structs[currentPos]

						if !exist {
							for _, pos := range keys {
								if  ok := x != nil; ok && currentPos >= pos {
									structName = structs[pos]
								}
							}
						}

						//Extract all db tags from the struct fields
						for _, field := range x.Fields.List {
							if field.Tag != nil {
								if  matches := dbRegex.FindStringSubmatch(field.Tag.Value); len(matches[1]) > 0 {
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

			}(tagChan)
		}
	}

	for i:= 0; i < chanLimit; i++ {
		tag := <- tagChan
		tags[tag[0]] = append(tags[tag[0]], tag[1])
	}

	return tags
}
