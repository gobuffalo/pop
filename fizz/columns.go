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
	fizzers["drop_column"] = DropColumn
	fizzers["rename_column"] = RenameColumn
}

func AddColumn(ch chan *Bubble) interface{} {
	return func(table, name, ctype string, options Options) {
		t := &Table{
			Name: table,
			Columns: []Column{
				{Name: name, ColType: ctype, Options: options},
			},
		}
		ch <- &Bubble{BubbleType: E_ADD_COLUMN, Data: t}
	}
}

func DropColumn(ch chan *Bubble) interface{} {
	return func(table, name string) {
		t := &Table{
			Name: table,
			Columns: []Column{
				{Name: name},
			},
		}
		ch <- &Bubble{BubbleType: E_DROP_COLUMN, Data: t}
	}
}

func RenameColumn(ch chan *Bubble) interface{} {
	return func(table, old, new string) {
		t := &Table{
			Name: table,
			Columns: []Column{
				{Name: old},
				{Name: new},
			},
		}
		ch <- &Bubble{BubbleType: E_RENAME_COLUMN, Data: t}
	}
}
