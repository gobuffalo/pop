package fizz

import (
	"os"
	"strings"

	"github.com/mattn/anko/core"
	"github.com/mattn/anko/vm"
	"github.com/pkg/errors"
)

type BubbleType int

// Bubbler is used to transform fizz script to SQL, using a Translator.
type Bubbler struct {
	Translator
	data []string
}

// NewBubbler constructs a new Bubbler.
func NewBubbler(t Translator) *Bubbler {
	return &Bubbler{
		Translator: t,
		data:       []string{},
	}
}

func (b *Bubbler) String() string {
	return strings.Join(b.data, "\n")
}

// Bubble transforms a fizz script to SQL.
func (b *Bubbler) Bubble(s string) (string, error) {
	env := core.Import(vm.NewEnv())

	f := fizzer{b}

	// columns:
	env.Define("change_column", f.ChangeColumn())
	env.Define("add_column", f.AddColumn())
	env.Define("drop_column", f.DropColumn())
	env.Define("rename_column", f.RenameColumn())

	env.Define("raw", f.RawSQL())
	env.Define("exec", f.Exec(os.Stdout))

	// indexes:
	env.Define("add_index", f.AddIndex())
	env.Define("drop_index", f.DropIndex())
	env.Define("rename_index", f.RenameIndex())

	// foreign keys
	env.Define("add_foreign_key", f.AddForeignKey())
	env.Define("drop_foreign_key", f.DropForeignKey())

	// tables:
	env.Define("create_table", f.CreateTable())
	env.Define("drop_table", f.DropTable())
	env.Define("rename_table", f.RenameTable())

	_, err := env.Execute(s)
	return b.String(), errors.Wrap(err, "parse error")
}
