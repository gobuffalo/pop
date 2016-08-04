package translators

import "github.com/markbates/pop/fizz"

type FizzTranslator interface {
	CreateTable(fizz.Table) (string, error)
	DropTable(fizz.Table) (string, error)
	RenameTable([]fizz.Table) (string, error)
	AddColumn(fizz.Table) (string, error)
	DropColumn(fizz.Table) (string, error)
	RenameColumn(fizz.Table) (string, error)
	AddIndex(fizz.Table) (string, error)
	DropIndex(fizz.Index) (string, error)
	RenameIndex([]fizz.Index) (string, error)
}
