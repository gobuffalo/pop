package fizz

import (
	"io/ioutil"
	"log"

	"github.com/mattn/anko/builtins"
	"github.com/mattn/anko/vm"
)

type Options map[string]interface{}

type fizzer func(chan *Bubble) interface{}

var fizzers = map[string]fizzer{}

func AFile(p string) Bubbler {
	b, err := ioutil.ReadFile(p)
	if err != nil {
		log.Fatal(err)
	}
	return AString(string(b))
}

func AString(s string) Bubbler {
	ch := make(chan *Bubble)
	b := Bubbler{Bubbles: ch}
	go func() {
		env := core.Import(vm.NewEnv())

		for k, v := range fizzers {
			env.Define(k, v(ch))
		}

		_, err := env.Execute(s)
		if err != nil {
			log.Fatal(err)
		}
		close(ch)
	}()
	return b
}
