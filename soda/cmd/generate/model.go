package generate

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/markbates/inflect"
	"github.com/spf13/cobra"
)

var ModelCmd = &cobra.Command{
	Use:     "model [name]",
	Aliases: []string{"m"},
	Short:   "Generates a model for your database",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return errors.New("You must supply a name for your model!")
		}
		name := args[0]
		table := inflect.Tableize(name)
		mname := inflect.Camelize(name)

		s := []string{"package models\n", "import \"time\"\n"}
		s = append(s, fmt.Sprintf("// %s models the database table '%s'", mname, table))
		s = append(s, fmt.Sprintf("type %s struct {", mname))
		s = append(s, fmt.Sprintf("\tID int `json:\"id\" db:\"id\"`"))

		nrx := regexp.MustCompile(`^nulls.(.+)`)
		hasNulls := false
		for _, def := range args[1:] {
			col := strings.Split(def, ":")
			if len(col) == 1 {
				col = append(col, "string")
			}
			if nrx.MatchString(col[1]) {
				hasNulls = true
			}
			s = append(s, fmt.Sprintf("\t%s %s `json:\"%s\" db:\"%s\"`", inflect.Camelize(col[0]), colType(col[1]), col[0], col[0]))
		}
		s = append(s, fmt.Sprintf("\tCreatedAt time.Time `json:\"created_at\" db:\"created_at\"`"))
		s = append(s, fmt.Sprintf("\tUpdatedAt time.Time `json:\"updated_at\" db:\"updated_at\"`"))

		pname := inflect.Pluralize(mname)
		s = append(s, "}")
		s = append(s, fmt.Sprintf("\ntype %s []%s", pname, mname))

		if hasNulls {
			s2 := []string{s[0], `import "github.com/markbates/going/nulls"`}
			s = append(s2, s[1:]...)
		}
		err := os.MkdirAll("models", 0766)
		if err != nil {
			return err
		}

		uname := inflect.Underscore(name)
		fname := filepath.Join("models", uname+".go")
		err = ioutil.WriteFile(fname, []byte(strings.Join(s, "\n")), 0766)
		if err != nil {
			return err
		}

		err = ioutil.WriteFile(filepath.Join("models", uname+"_test.go"), []byte(`package models_test`), 0766)
		if err != nil {
			return err
		}

		md, _ := filepath.Abs(fname)
		goi := exec.Command("gofmt", "-w", md)
		out, err := goi.CombinedOutput()
		if err != nil {
			fmt.Printf("Received an error when trying to run gofmt -> %#v\n", err)
			fmt.Println(out)
		}

		b, err := ioutil.ReadFile(fname)
		if err != nil {
			return err
		}

		fmt.Println(string(b))

		return nil
	},
}

func colType(s string) string {
	switch s {
	case "text":
		return "string"
	case "time", "timestamp":
		return "time.Time"
	default:
		return s
	}
	return s
}
