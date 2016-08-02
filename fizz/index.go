package fizz

import (
	"fmt"
	"strings"
)

func init() {
	fizzers["add_index"] = AddIndex
	fizzers["drop_index"] = DropIndex
	fizzers["rename_index"] = RenameIndex
}

type Index struct {
	TableName string
	Name      string
	Columns   []string
	Unique    bool
	Options   Options
}

func AddIndex(ch chan *Bubble) interface{} {
	return func(table string, columns interface{}, options Options) {
		i := &Index{
			TableName: table,
		}
		switch t := columns.(type) {
		default:
			fmt.Printf("unexpected type %T\n", t) // %T prints whatever type t has
		case string:
			i.Columns = []string{t}
		case []interface{}:
			cl := make([]string, len(t))
			for i, c := range t {
				cl[i] = c.(string)
			}
			i.Columns = cl
		}

		if options["name"] != nil {
			i.Name = options["name"].(string)
		} else {
			i.Name = fmt.Sprintf("%s_%s_idx", i.TableName, strings.Join(i.Columns, "_"))
		}
		i.Unique = options["unique"] != nil
		ch <- &Bubble{BubbleType: E_ADD_INDEX, Data: i}
	}
}

func DropIndex(ch chan *Bubble) interface{} {
	return func(table, name string) {
		ch <- &Bubble{BubbleType: E_DROP_INDEX, Data: &Index{TableName: table, Name: name}}
	}
}

func RenameIndex(ch chan *Bubble) interface{} {
	return func(table, old, new string) {
		ch <- &Bubble{
			BubbleType: E_RENAME_INDEX,
			Data: []*Index{
				{TableName: table, Name: old},
				{TableName: table, Name: new},
			},
		}
	}
}
