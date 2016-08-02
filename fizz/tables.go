package fizz

type Table struct {
	Name    string
	Columns []Column
}

func (t *Table) Column(name string, colType string, options map[string]interface{}) {
	c := Column{
		Name:    name,
		ColType: colType,
		Options: options,
	}
	t.Columns = append(t.Columns, c)
}

func init() {
	fizzers["create_table"] = CreateTable
	fizzers["drop_table"] = DropTable
	fizzers["rename_table"] = RenameTable
}

func CreateTable(ch chan *Bubble) interface{} {
	return func(name string, fn func(t *Table)) {
		t := &Table{
			Name:    name,
			Columns: []Column{ID_COL, CREATED_COL, UPDATED_COL},
		}
		fn(t)
		ch <- &Bubble{BubbleType: E_CREATE_TABLE, Data: t}
	}
}

func DropTable(ch chan *Bubble) interface{} {
	return func(name string) {
		ch <- &Bubble{BubbleType: E_DROP_TABLE, Data: &Table{Name: name}}
	}
}

func RenameTable(ch chan *Bubble) interface{} {
	return func(old, new string) {
		ch <- &Bubble{
			BubbleType: E_RENAME_TABLE,
			Data: []*Table{
				{Name: old},
				{Name: new},
			},
		}
	}
}
