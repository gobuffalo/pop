package fizz

import (
	"io/ioutil"
	"log"
)

type Options map[string]interface{}

type fizzer struct {
	Bubbler *Bubbler
}

func (f fizzer) add(b Bubble) {
	f.Bubbler.Bubbles = append(f.Bubbler.Bubbles, b)
}

func AFile(p string) (*Bubbler, error) {
	b, err := ioutil.ReadFile(p)
	if err != nil {
		log.Fatal(err)
	}
	return AString(string(b))
}

func AString(s string) (*Bubbler, error) {
	b := NewBubbler()
	err := b.Bubble(s)
	return b, err
}
