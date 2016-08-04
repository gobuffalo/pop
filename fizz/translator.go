package fizz

type Translator interface {
	CreateTable(Table) (string, error)
	DropTable(Table) (string, error)
	RenameTable([]Table) (string, error)
	AddColumn(Table) (string, error)
	DropColumn(Table) (string, error)
	RenameColumn(Table) (string, error)
	AddIndex(Table) (string, error)
	DropIndex(Index) (string, error)
	RenameIndex([]Index) (string, error)
}
