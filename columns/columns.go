package columns

import (
	"fmt"
	"sort"
	"strings"
	"sync"
)

type Columns struct {
	Cols      map[string]*Column
	lock      *sync.RWMutex
	TableName string
}

// Add a column to the list.
func (c *Columns) Add(names ...string) []*Column {
	ret := []*Column{}
	c.lock.Lock()
	for _, name := range names {

		var xs []string
		var col *Column
		ss := ""
		//support for distinct xx, or distinct on (field) table.fields
		if strings.Contains(name, "distinct") {
			i := strings.LastIndex(name, " ")
			if i > -1 {
				colName := name[i+1 : len(name)]
				xs = strings.Split(colName, ",")
				ss = name
				if len(xs) > 1 {
					ss = ss[0 : len(ss)-len(xs[1])-1]
				}
			} else {
				panic("distinct without name")
			}
		} else {
			xs = strings.Split(name, ",")
		}

		col = c.Cols[xs[0]]
		//fmt.Printf("column: name: %v, col: %v, xs: %v, ss: %v\n", name, col, xs, ss)
		if col == nil {
			if ss == "" {
				ss = xs[0]
				if c.TableName != "" {
					ss = fmt.Sprintf("%s.%s", c.TableName, ss)
				}
			}

			col = &Column{
				Name:      xs[0],
				SelectSQL: ss,
				Readable:  true,
				Writeable: true,
			}

			if len(xs) > 1 {
				if xs[1] == "r" {
					col.Writeable = false
				}
				if xs[1] == "w" {
					col.Readable = false
				}
			} else if col.Name == "id" {
				col.Writeable = false
			}

			c.Cols[col.Name] = col
		}
		ret = append(ret, col)
	}
	c.lock.Unlock()
	return ret
}

// Remove a column from the list.
func (c *Columns) Remove(names ...string) {
	for _, name := range names {
		xs := strings.Split(name, ",")
		name = xs[0]
		delete(c.Cols, name)
	}
}

func (c Columns) Writeable() *WriteableColumns {
	w := &WriteableColumns{NewColumns(c.TableName)}
	for _, col := range c.Cols {
		if col.Writeable {
			w.Cols[col.Name] = col
		}
	}
	return w
}

func (c Columns) Readable() *ReadableColumns {
	w := &ReadableColumns{NewColumns(c.TableName)}
	for _, col := range c.Cols {
		if col.Readable {
			w.Cols[col.Name] = col
		}
	}
	return w
}

func (c Columns) String() string {
	xs := []string{}
	for _, t := range c.Cols {
		xs = append(xs, t.Name)
	}
	sort.Strings(xs)
	return strings.Join(xs, ", ")
}

func (c Columns) SymbolizedString() string {
	xs := []string{}
	for _, t := range c.Cols {
		xs = append(xs, ":"+t.Name)
	}
	sort.Strings(xs)
	return strings.Join(xs, ", ")
}

func NewColumns(tableName string) Columns {
	return Columns{
		lock:      &sync.RWMutex{},
		Cols:      map[string]*Column{},
		TableName: tableName,
	}
}
