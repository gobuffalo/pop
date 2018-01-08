package grift

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_Add(t *testing.T) {
	r := require.New(t)
	r.Empty(griftList)

	Add("hello", func(c *Context) error {
		return nil
	})
	r.Len(griftList, 1)

	reset()
}

func Test_Add_Multiple(t *testing.T) {
	r := require.New(t)
	r.Empty(griftList)

	log := []string{}
	Add("hello", func(c *Context) error {
		log = append(log, "hello1")
		return nil
	})
	r.Len(griftList, 1)

	Add("hello", func(c *Context) error {
		log = append(log, "hello2")
		return nil
	})
	r.Len(griftList, 1)

	err := Run("hello", NewContext("hello"))
	r.NoError(err)
	r.Len(log, 2)

	reset()
}

func Test_Set(t *testing.T) {
	r := require.New(t)
	r.Empty(griftList)

	log := []string{}
	Set("hello", func(c *Context) error {
		log = append(log, "hello1")
		return nil
	})
	r.Len(griftList, 1)

	Set("hello", func(c *Context) error {
		log = append(log, "hello2")
		return nil
	})
	r.Len(griftList, 1)

	err := Run("hello", NewContext("hello"))
	r.NoError(err)
	r.Len(log, 1)
	r.Equal("hello2", log[0])

	reset()
}

func Test_Rename(t *testing.T) {
	r := require.New(t)

	Add("hello", func(c *Context) error {
		return nil
	})
	Rename("hello", "hi")

	name := List()[0]
	r.Equal("hi", name)

	reset()
}

func Test_Remove(t *testing.T) {
	r := require.New(t)

	Add("hello", func(c *Context) error {
		return nil
	})
	r.Len(griftList, 1)

	Remove("hello")
	r.Len(griftList, 0)
}

func Test_Desc(t *testing.T) {
	r := require.New(t)

	Desc("hello", "Hello!!")
	r.Equal("Hello!!", descriptions["hello"])

	reset()
}

func Test_Run(t *testing.T) {
	r := require.New(t)

	var msg string

	Add("hello", func(c *Context) error {
		msg = "Hello, World!"
		return nil
	})

	err := Run("hello", NewContext("hello"))
	r.NoError(err)
	r.Equal("Hello, World!", msg)

	reset()
}

func Test_List(t *testing.T) {
	r := require.New(t)

	Add("b", func(c *Context) error {
		return nil
	})
	Add("a", func(c *Context) error {
		return nil
	})
	Add("c", func(c *Context) error {
		return nil
	})

	r.Equal([]string{"a", "b", "c"}, List())

	reset()
}

func Test_Namespace(t *testing.T) {
	r := require.New(t)
	Add("b", func(c *Context) error {
		return nil
	})
	Add("c", func(c *Context) error {
		return nil
	})

	Namespace("a", func() {
		Add("a", func(c *Context) error {
			return nil
		})
		Add("d", func(c *Context) error {
			return nil
		})

		Remove("b")
		Remove(":c")

		Namespace("e", func() {
			Add("f", func(c *Context) error {
				return nil
			})

			Rename("f", "g")

			Remove(":d")

		})

	})

	r.Equal([]string{"a:a", "a:d", "a:e:g", "b"}, List())

	reset()
}

func Test_PrintGrifts(t *testing.T) {
	r := require.New(t)

	Add("b", func(c *Context) error {
		return nil
	})
	Desc("a", "AH!")
	Add("a", func(c *Context) error {
		return nil
	})

	bb := &bytes.Buffer{}
	PrintGrifts(bb)
	r.Equal("Available grifts\n================\ngrift a    # AH!\ngrift b    # \n", bb.String())
	reset()
}

func Test_Exec(t *testing.T) {
	r := require.New(t)

	var name string
	var args []string

	Add("hello", func(c *Context) error {
		name = c.Name
		args = c.Args
		return nil
	})
	t.Run("list", func(st *testing.T) {
		Exec([]string{}, false)
		r.Equal("", name)
		r.Empty(args)
	})
	t.Run("name grift", func(st *testing.T) {
		Exec([]string{"hello"}, false)
		r.Equal("hello", name)
		r.Empty(args)
	})
	t.Run("name grift with args", func(st *testing.T) {
		Exec([]string{"hello", "a", "b"}, false)
		r.Equal("hello", name)
		r.Equal([]string{"a", "b"}, args)
	})
	reset()
}

func reset() {
	griftList = map[string]Grift{}
	descriptions = map[string]string{}
	namespace = ""
}
