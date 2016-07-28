package fizz

import (
	"io/ioutil"
	"log"

	"github.com/mattn/anko/builtins"
	"github.com/mattn/anko/vm"
)

func AFile(p string) chan Bubble {
	b, err := ioutil.ReadFile(p)
	if err != nil {
		log.Fatal(err)
	}
	return AString(string(b))
}

func AString(s string) chan Bubble {
	ch := make(chan Bubble)
	go func() {
		env := core.Import(vm.NewEnv())

		env.Define("raw", RawSQL(ch))
		env.Define("create_table", CreateTable(ch))
		env.Define("drop_table", DropTable(ch))
		env.Define("add_column", AddColumn(ch))

		_, err := env.Execute(s)
		if err != nil {
			log.Fatal(err)
		}
		close(ch)
	}()
	return ch
}

func RawSQL(ch chan Bubble) interface{} {
	return func(sql string) {
		ch <- Bubble{Type: E_RAW_SQL, Data: sql}
	}
}
