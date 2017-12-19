package fizz

import (
	"fmt"
)

type ForeignKeyRef struct {
	Table   string
	Columns []string
}

type ForeignKey struct {
	Name       string
	Column     string
	References ForeignKeyRef
	Options    Options
}

func (f fizzer) AddForeignKey() interface{} {
	return func(table string, name string, column string, refs interface{}, options Options) {
		f.add(f.Bubbler.AddForeignKey(Table{
			Name: table,
			ForeignKeys: []ForeignKey{
				{
					Name:       name,
					Column:     column,
					References: parseForeignKeyRef(refs),
					Options:    options,
				},
			},
		}))
	}
}

func (f fizzer) DropForeignKey() interface{} {
	return func(table string, fk string, options Options) {
		f.add(f.Bubbler.DropForeignKey(Table{
			Name: table,
			ForeignKeys: []ForeignKey{
				{
					Name:    fk,
					Options: options,
				},
			},
		}))
	}
}

func parseForeignKeyRef(refs interface{}) (fkr ForeignKeyRef) {
	refMap, ok := refs.(map[string]interface{})
	if !ok {
		fmt.Printf(`invalid references format %s\nmust be "{"table": ["colum1", "column2"]}"`, refs)
		return
	}
	if len(refMap) > 1 {
		fmt.Printf("only one table is supported as Foreign key reference")
		return
	}
	for table, columns := range refMap {
		fkr.Table = table
		for _, c := range columns.([]interface{}) {
			fkr.Columns = append(fkr.Columns, fmt.Sprintf("%s", c))
		}
	}

	return
}
