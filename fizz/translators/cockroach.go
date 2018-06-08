package translators

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"

	"github.com/gobuffalo/pop/fizz"
)

// Cockroach is the fizz translator implementation for cockroach.
type Cockroach struct {
	Schema SchemaQuery
}

// NewCockroach constructs a new Cockroach translator.
func NewCockroach(url string, name string) *Cockroach {
	schema := &cockroachSchema{
		Schema{
			URL:    url,
			Name:   name,
			schema: map[string]*fizz.Table{},
		},
	}
	schema.Builder = schema
	return &Cockroach{Schema: schema}
}

// CreateTable backs fizz create table command.
func (p *Cockroach) CreateTable(t fizz.Table) (string, error) {
	p.Schema.SetTable(&t)
	sql := []string{}
	cols := []string{}
	var s string
	for _, c := range t.Columns {
		if c.Primary {
			switch c.ColType {
			case "string", "uuid":
				s = fmt.Sprintf("\"%s\" %s PRIMARY KEY", c.Name, p.colType(c))
			case "integer", "int", "INT":
				s = fmt.Sprintf("\"%s\" SERIAL PRIMARY KEY", c.Name)
			default:
				return "", errors.Errorf("can not use %s as a primary key", c.ColType)
			}
		} else {
			s = p.buildAddColumn(c)
		}
		cols = append(cols, s)
	}

	for _, fk := range t.ForeignKeys {
		cols = append(cols, p.buildForeignKey(t, fk, true))
	}

	s = fmt.Sprintf("CREATE TABLE \"%s\" (\n%s\n);COMMIT TRANSACTION;BEGIN TRANSACTION;", t.Name, strings.Join(cols, ",\n"))
	sql = append(sql, s)

	for _, i := range t.Indexes {
		s, err := p.AddIndex(fizz.Table{
			Name:    t.Name,
			Indexes: []fizz.Index{i},
		})
		if err != nil {
			return "", err
		}
		sql = append(sql, s)
	}

	return strings.Join(sql, "\n"), nil
}

// DropTable backs fizz drop table command.
func (p *Cockroach) DropTable(t fizz.Table) (string, error) {
	p.Schema.Delete(t.Name)
	return fmt.Sprintf("DROP TABLE \"%s\";COMMIT TRANSACTION;BEGIN TRANSACTION;", t.Name), nil
}

// RenameTable backs fizz rename table command.
func (p *Cockroach) RenameTable(t []fizz.Table) (string, error) {
	if len(t) < 2 {
		return "", errors.New("not enough table names supplied")
	}
	oldName := t[0].Name
	newName := t[1].Name
	tableInfo, err := p.Schema.TableInfo(oldName)
	if err != nil {
		return "", err
	}
	tableInfo.Name = newName

	return fmt.Sprintf("ALTER TABLE \"%s\" RENAME TO \"%s\";COMMIT TRANSACTION;BEGIN TRANSACTION;", oldName, newName), nil
}

// ChangeColumn backs fizz change column command.
func (p *Cockroach) ChangeColumn(t fizz.Table) (string, error) {
	if len(t.Columns) == 0 {
		return "", errors.New("not enough columns supplied")
	}
	c := t.Columns[0]

	tableInfo, err := p.Schema.TableInfo(t.Name)

	if err != nil {
		return "", err
	}

	for i := range tableInfo.Columns {
		if tableInfo.Columns[i].Name == t.Columns[0].Name {
			tableInfo.Columns[i] = c
			break
		}
	}

	sql := []string{}
	s, err := p.withTempColumn(t.Name, c.Name, func(table fizz.Table, origCol fizz.Column, tempCol string) (string, error) {
		newCol := p.buildChangeColumn(origCol, c)
		err1 := p.Schema.ReplaceColumn(table.Name, origCol.Name, newCol)
		if err1 != nil {
			return "", err1
		}

		createColumnSQL := fmt.Sprintf("ALTER TABLE \"%s\" ADD COLUMN %s;COMMIT TRANSACTION;BEGIN TRANSACTION;", table.Name, p.buildAddColumn(newCol))
		ins := fmt.Sprintf("UPDATE \"%s\" SET \"%s\" = \"%s\";COMMIT TRANSACTION;BEGIN TRANSACTION;", t.Name, c.Name, tempCol)
		return strings.Join([]string{createColumnSQL, ins}, "\n"), nil
	})

	if err != nil {
		return "", err
	}

	sql = append(sql, s)

	return strings.Join(sql, "\n"), nil
}

// AddColumn backs fizz add column command.
func (p *Cockroach) AddColumn(t fizz.Table) (string, error) {
	if len(t.Columns) == 0 {
		return "", errors.New("not enough columns supplied")
	}
	c := t.Columns[0]
	s := fmt.Sprintf("ALTER TABLE \"%s\" ADD COLUMN %s;COMMIT TRANSACTION;BEGIN TRANSACTION;", t.Name, p.buildAddColumn(c))

	//Update schema cache if we can
	tableInfo, err := p.Schema.TableInfo(t.Name)
	if err == nil {
		found := false
		for i := range tableInfo.Columns {
			if tableInfo.Columns[i].Name == c.Name {
				tableInfo.Columns[i] = c
				break
			}
		}
		if !found {
			tableInfo.Columns = append(tableInfo.Columns, c)
		}
	}

	return s, nil
}

// DropColumn backs fizz drop column command.
func (p *Cockroach) DropColumn(t fizz.Table) (string, error) {
	if len(t.Columns) == 0 {
		return "", errors.New("not enough columns supplied")
	}
	c := t.Columns[0]
	p.Schema.DeleteColumn(t.Name, c.Name)
	return fmt.Sprintf("ALTER TABLE \"%s\" DROP COLUMN \"%s\";COMMIT TRANSACTION;BEGIN TRANSACTION;", t.Name, c.Name), nil
}

// RenameColumn backs fizz rename column command.
func (p *Cockroach) RenameColumn(t fizz.Table) (string, error) {
	if len(t.Columns) < 2 {
		return "", errors.New("not enough columns supplied")
	}

	oc := t.Columns[0]
	nc := t.Columns[1]

	tableInfo, err := p.Schema.TableInfo(t.Name)
	if err != nil {
		return "", err
	}

	for ic, c := range tableInfo.Columns {
		if c.Name == oc.Name {
			tableInfo.Columns[ic].Name = nc.Name
		}
	}

	s := fmt.Sprintf("ALTER TABLE \"%s\" RENAME COLUMN \"%s\" TO \"%s\";COMMIT TRANSACTION;BEGIN TRANSACTION;", t.Name, oc.Name, nc.Name)
	return s, nil
}

// AddIndex backs fizz add index command.
func (p *Cockroach) AddIndex(t fizz.Table) (string, error) {
	if len(t.Indexes) == 0 {
		return "", errors.New("not enough indexes supplied")
	}
	i := t.Indexes[0]
	s := fmt.Sprintf("CREATE INDEX \"%s\" ON \"%s\" (%s);COMMIT TRANSACTION;BEGIN TRANSACTION;", i.Name, t.Name, strings.Join(i.Columns, ", "))
	if i.Unique {
		s = strings.Replace(s, "CREATE", "CREATE UNIQUE", 1)
	}

	tableInfo, err := p.Schema.TableInfo(t.Name)
	if err != nil {
		return "", err
	}
	tableInfo.Indexes = append(tableInfo.Indexes, i)

	return s, nil
}

// DropIndex backs fizz drop index command.
func (p *Cockroach) DropIndex(t fizz.Table) (string, error) {
	if len(t.Indexes) == 0 {
		return "", errors.New("not enough indexes supplied")
	}
	i := t.Indexes[0]

	tableInfo, err := p.Schema.TableInfo(t.Name)
	if err != nil {
		return "", err
	}

	newIndexes := []fizz.Index{}
	for _, c := range tableInfo.Indexes {
		if c.Name != i.Name {
			newIndexes = append(newIndexes, c)
		}
	}
	tableInfo.Indexes = newIndexes

	return fmt.Sprintf("DROP INDEX IF EXISTS \"%s\";COMMIT TRANSACTION;BEGIN TRANSACTION;", i.Name), nil
}

// RenameIndex backs fizz rename index command.
func (p *Cockroach) RenameIndex(t fizz.Table) (string, error) {
	ix := t.Indexes
	if len(ix) < 2 {
		return "", errors.New("not enough indexes supplied")
	}
	oi := ix[0]
	ni := ix[1]

	tableInfo, err := p.Schema.TableInfo(t.Name)
	if err != nil {
		return "", err
	}

	for _, c := range tableInfo.Indexes {
		if c.Name == oi.Name {
			c.Name = ni.Name
		}
	}

	return fmt.Sprintf("ALTER INDEX \"%s\"@\"%s\" RENAME TO \"%s\";COMMIT TRANSACTION;BEGIN TRANSACTION;", t.Name, oi.Name, ni.Name), nil
}

// AddForeignKey backs fizz add foreign key command.
func (p *Cockroach) AddForeignKey(t fizz.Table) (string, error) {
	if len(t.ForeignKeys) == 0 {
		return "", errors.New("not enough foreign keys supplied")
	}

	tableInfo, err := p.Schema.TableInfo(t.Name)
	if err != nil {
		return "", err
	}
	tableInfo.ForeignKeys = append(tableInfo.ForeignKeys, t.ForeignKeys[0])

	return p.buildForeignKey(t, t.ForeignKeys[0], false), nil
}

// DropForeignKey backs fizz drop foreign key command.
func (p *Cockroach) DropForeignKey(t fizz.Table) (string, error) {
	if len(t.ForeignKeys) == 0 {
		return "", errors.New("not enough foreign keys supplied")
	}

	fk := t.ForeignKeys[0]

	tableInfo, err := p.Schema.TableInfo(t.Name)
	if err != nil {
		return "", err
	}
	newFKs := []fizz.ForeignKey{}
	for _, key := range tableInfo.ForeignKeys {
		if key.Name != fk.Name {
			newFKs = append(newFKs, key)
		}
	}
	tableInfo.ForeignKeys = newFKs

	var ifExists string
	if v, ok := fk.Options["if_exists"]; ok && v.(bool) {
		ifExists = "IF EXISTS"
	}

	s := fmt.Sprintf("ALTER TABLE %s DROP CONSTRAINT %s %s;COMMIT TRANSACTION;BEGIN TRANSACTION;", t.Name, ifExists, fk.Name)
	return s, nil
}

func (p *Cockroach) buildAddColumn(c fizz.Column) string {
	s := fmt.Sprintf("\"%s\" %s", c.Name, p.colType(c))

	if c.Options["null"] == nil {
		s = fmt.Sprintf("%s NOT NULL", s)
	}
	if c.Options["default"] != nil {
		s = fmt.Sprintf("%s DEFAULT '%v'", s, c.Options["default"])
	}
	if c.Options["default_raw"] != nil {
		s = fmt.Sprintf("%s DEFAULT %s", s, c.Options["default_raw"])
	}

	return s
}

func (p *Cockroach) buildChangeColumn(oldCol fizz.Column, c fizz.Column) fizz.Column {
	newCol := fizz.Column{
		Name:    c.Name,
		ColType: oldCol.ColType,
		Options: oldCol.Options,
		Primary: oldCol.Primary,
	}

	if c.Options["default"] != nil {
		newCol.Options["default"] = c.Options["default"]
	}
	if c.Options["null"] != nil {
		newCol.Options["null"] = c.Options["null"]
	}
	if c.Options["default_raw"] != nil {
		newCol.Options["default_raw"] = c.Options["default_raw"]
	}

	return newCol
}

func (p *Cockroach) withTempTable(table string, fn func(fizz.Table) (string, error)) (string, error) {
	tempTable := fizz.Table{Name: fmt.Sprintf("_%s_tmp", table)}

	sql := []string{fmt.Sprintf("ALTER TABLE \"%s\" RENAME TO \"%s\";COMMIT TRANSACTION;BEGIN TRANSACTION;", table, tempTable.Name)}
	s, err := fn(tempTable)
	if err != nil {
		return "", err
	}
	sql = append(sql, s, fmt.Sprintf("DROP TABLE \"%s\";COMMIT TRANSACTION;BEGIN TRANSACTION;", tempTable.Name))

	return strings.Join(sql, "\n"), nil
}

func (p *Cockroach) withTempColumn(tableName string, column string, fn func(fizz.Table, fizz.Column, string) (string, error)) (string, error) {
	table, err := p.Schema.TableInfo(tableName)
	if err != nil {
		return "", err
	}
	col, err1 := p.Schema.ColumnInfo(tableName, column)
	if err1 != nil {
		return "", err1
	}

	tempCol := fmt.Sprintf("_%s_tmp", column)

	sql := []string{fmt.Sprintf("ALTER TABLE \"%s\" RENAME COLUMN \"%s\" TO \"%s\";COMMIT TRANSACTION;BEGIN TRANSACTION;", tableName, column, tempCol)}

	s, err := fn(*table, *col, tempCol)
	if err != nil {
		return "", err
	}
	sql = append(sql, s, fmt.Sprintf("ALTER TABLE \"%s\" DROP COLUMN \"%s\";COMMIT TRANSACTION;BEGIN TRANSACTION;", tableName, tempCol))

	return strings.Join(sql, "\n"), nil
}

func (p *Cockroach) colType(c fizz.Column) string {
	switch c.ColType {
	case "string":
		s := "255"
		if c.Options["size"] != nil {
			s = fmt.Sprintf("%d", c.Options["size"])
		}
		return fmt.Sprintf("VARCHAR (%s)", s)
	case "uuid":
		return "UUID"
	case "time", "datetime":
		return "timestamp"
	case "blob":
		return "BYTES"
	default:
		return c.ColType
	}
}

func (p *Cockroach) buildForeignKey(t fizz.Table, fk fizz.ForeignKey, onCreate bool) string {
	refs := fmt.Sprintf("%s (%s)", fk.References.Table, strings.Join(fk.References.Columns, ", "))
	s := fmt.Sprintf("CONSTRAINT %s FOREIGN KEY (%s) REFERENCES %s", fk.Name, fk.Column, refs)

	if !onCreate {
		s = fmt.Sprintf("ALTER TABLE %s ADD %s;COMMIT TRANSACTION;BEGIN TRANSACTION;", t.Name, s)
	}

	return s
}
