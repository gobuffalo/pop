package ctable

import "github.com/gobuffalo/fizz"

type mockTranslator struct{}

func (mockTranslator) Name() string {
	return "test"
}

func (mockTranslator) CreateTable(fizz.Table) (string, error) {
	return "create table;", nil
}

func (mockTranslator) DropTable(fizz.Table) (string, error) {
	return "drop table;", nil
}

func (mockTranslator) RenameTable([]fizz.Table) (string, error) {
	return "rename table;", nil
}

func (mockTranslator) AddColumn(fizz.Table) (string, error) {
	return "add column;", nil
}

func (mockTranslator) ChangeColumn(fizz.Table) (string, error) {
	return "change column;", nil
}

func (mockTranslator) DropColumn(fizz.Table) (string, error) {
	return "drop column;", nil
}

func (mockTranslator) RenameColumn(fizz.Table) (string, error) {
	return "rename column;", nil
}

func (mockTranslator) AddIndex(fizz.Table) (string, error) {
	return "add index;", nil
}

func (mockTranslator) DropIndex(fizz.Table) (string, error) {
	return "drop index;", nil
}

func (mockTranslator) RenameIndex(fizz.Table) (string, error) {
	return "rename index;", nil
}

func (mockTranslator) AddForeignKey(fizz.Table) (string, error) {
	return "add foreign key;", nil
}

func (mockTranslator) DropForeignKey(fizz.Table) (string, error) {
	return "drop foreign key;", nil
}
