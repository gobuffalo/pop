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

func init() {
	fizzers["add_column"] = AddColumn
}

func AddColumn(ch chan Bubble) interface{} {
	return func(table, name, ctype string, options Options) {
		t := &Table{
			Name: table,
			Columns: []Column{
				{Name: name, ColType: ctype, Options: options},
			},
		}
		ch <- Bubble{Type: E_ADD_COLUMN, Data: t}
	}
}
