package fizz

var ID_COL = Column{
	Name:    "id",
	Primary: true,
	ColType: "integer",
	Options: Options{},
}
var CREATED_COL = Column{Name: "created_at", ColType: "timestamp", Options: Options{}}
var UPDATED_COL = Column{Name: "updated_at", ColType: "timestamp", Options: Options{}}

type Column struct {
	Name    string
	ColType string
	Primary bool
	Options map[string]interface{}
}

func (f fizzer) AddColumn() interface{} {
	return func(table, name, ctype string, options Options) {
		t := Table{
			Name: table,
			Columns: []Column{
				{Name: name, ColType: ctype, Options: options},
			},
		}
		f.add(Bubble{BubbleType: E_ADD_COLUMN, Data: t})
	}
}

func (f fizzer) DropColumn() interface{} {
	return func(table, name string) {
		t := Table{
			Name: table,
			Columns: []Column{
				{Name: name},
			},
		}
		f.add(Bubble{BubbleType: E_DROP_COLUMN, Data: t})
	}
}

func (f fizzer) RenameColumn() interface{} {
	return func(table, old, new string) {
		t := Table{
			Name: table,
			Columns: []Column{
				{Name: old},
				{Name: new},
			},
		}
		f.add(Bubble{BubbleType: E_RENAME_COLUMN, Data: t})
	}
}
