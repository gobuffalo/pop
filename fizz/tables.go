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

func CreateTable(ch chan Bubble) interface{} {
	return func(name string, fn func(t *Table)) {
		t := &Table{
			Name:    name,
			Columns: []Column{ID_COL, CREATED_COL, UPDATED_COL},
		}
		fn(t)
		ch <- Bubble{Type: E_CREATE_TABLE, Data: t}
	}
}

func DropTable(ch chan Bubble) interface{} {
	return func(name string) {
		ch <- Bubble{Type: E_DROP_TABLE, Data: &Table{Name: name}}
	}
}
