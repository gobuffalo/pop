package fizz

import (
	"fmt"

	"github.com/mattn/anko/vm"

	core "github.com/mattn/anko/builtins"
)

type BubbleType int

const (
	E_CREATE_TABLE BubbleType = iota
	E_DROP_TABLE
	E_RENAME_TABLE
	E_RAW_SQL
	E_ADD_COLUMN
	E_DROP_COLUMN
	E_RENAME_COLUMN
	E_ADD_INDEX
	E_DROP_INDEX
	E_RENAME_INDEX
)

type Bubble struct {
	BubbleType BubbleType
	Data       interface{}
}

type Bubbler struct {
	Bubbles []Bubble
}

func NewBubbler() *Bubbler {
	return &Bubbler{
		Bubbles: []Bubble{},
	}
}

func (b *Bubbler) Bubble(s string) error {
	env := core.Import(vm.NewEnv())

	f := fizzer{b}

	// columns:
	env.Define("add_column", f.AddColumn())
	env.Define("drop_column", f.DropColumn())
	env.Define("rename_column", f.RenameColumn())

	env.Define("raw", f.RawSql())

	// indexes:
	env.Define("add_index", f.AddIndex())
	env.Define("drop_index", f.DropIndex())
	env.Define("rename_index", f.RenameIndex())

	// tables:
	env.Define("create_table", f.CreateTable())
	env.Define("drop_table", f.DropTable())
	env.Define("rename_table", f.RenameTable())

	_, err := env.Execute(s)
	fmt.Printf("### err -> %#v\n", err)
	return err
}
