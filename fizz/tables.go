package fizz

type Table struct {
	Name    string
	Columns []Column
	Indexes []Index
}

func (t *Table) Column(name string, colType string, options map[string]interface{}) {
	c := Column{
		Name:    name,
		ColType: colType,
		Options: options,
	}
	t.Columns = append(t.Columns, c)
}

func (f fizzer) CreateTable() interface{} {
	return func(name string, fn func(t *Table)) {
		t := Table{
			Name:    name,
			Columns: []Column{ID_COL, CREATED_COL, UPDATED_COL},
		}
		fn(&t)
		f.add(f.Bubbler.CreateTable(t))
	}
}

func (f fizzer) DropTable() interface{} {
	return func(name string) {
		f.add(f.Bubbler.DropTable(Table{Name: name}))
	}
}

func (f fizzer) RenameTable() interface{} {
	return func(old, new string) {
		f.add(f.Bubbler.RenameTable([]Table{
			{Name: old},
			{Name: new},
		}))
	}
}
