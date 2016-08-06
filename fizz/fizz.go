package fizz

import (
	"io/ioutil"
	"log"
	"os"
)

type Options map[string]interface{}

type fizzer struct {
	Bubbler *Bubbler
}

func (f fizzer) add(s string, err error) error {
	if err != nil {
		panic(err.Error())
	}
	f.Bubbler.data = append(f.Bubbler.data, s)
	return nil
}

func AFile(f *os.File, t Translator) (string, error) {
	b, err := ioutil.ReadAll(f)
	if err != nil {
		log.Fatal(err)
	}
	return AString(string(b), t)
}

func AString(s string, t Translator) (string, error) {
	b := NewBubbler(t)
	return b.Bubble(s)
}
