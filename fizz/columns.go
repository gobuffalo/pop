package fizz

var ID_COL = Column{
	Name:    "id",
	Primary: true,
	ColType: "integer",
	Options: map[string]interface{}{},
}
var CREATED_COL = Column{Name: "created_at", ColType: "timestamp", Options: map[string]interface{}{}}
var UPDATED_COL = Column{Name: "updated_at", ColType: "timestamp", Options: map[string]interface{}{}}

type Column struct {
	Name    string
	ColType string
	Primary bool
	Options map[string]interface{}
}

func AddColumn(ch chan Bubble) interface{} {
	return func(table, name, ctype string, options map[string]interface{}) {
		t := &Table{
			Name: table,
			Columns: []Column{
				{Name: name, ColType: ctype, Options: options},
			},
		}
		ch <- Bubble{Type: E_ADD_COLUMN, Data: t}
	}
}
