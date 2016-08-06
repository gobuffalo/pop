package translators

import (
	"errors"
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/markbates/pop/fizz"
	_ "github.com/mattes/migrate/driver/sqlite3"
)

type sqliteIndexListInfo struct {
	Seq     int    `db:"seq"`
	Name    string `db:"name"`
	Unique  bool   `db:"unique"`
	Origin  string `db:"origin"`
	Partial string `db:"partial"`
}

type sqliteIndexInfo struct {
	Seq  int    `db:"seqno"`
	CID  int    `db:"cid"`
	Name string `db:"name"`
}

type sqliteTableInfo struct {
	CID     int         `db:"cid"`
	Name    string      `db:"name"`
	Type    string      `db:"type"`
	NotNull bool        `db:"notnull"`
	Default interface{} `db:"dflt_value"`
	PK      bool        `db:"pk"`
}

func (t sqliteTableInfo) ToColumn() fizz.Column {
	c := fizz.Column{
		Name:    t.Name,
		ColType: t.Type,
		Primary: t.PK,
		Options: fizz.Options{},
	}
	if !t.NotNull {
		c.Options["null"] = true
	}
	if t.Default != nil {
		c.Options["default"] = strings.TrimSuffix(strings.TrimPrefix(fmt.Sprintf("%s", t.Default), "'"), "'")
	}
	return c
}

type SQLite struct {
	URL string
	sql []string
}

func NewSQLite(url string) *SQLite {
	return &SQLite{
		URL: url,
		sql: []string{},
	}
}

func (p *SQLite) add(s string) {
	p.sql = append(p.sql, s)
}

func (p *SQLite) CreateTable(t fizz.Table) (string, error) {
	cols := []string{}
	var s string
	for _, c := range t.Columns {
		if c.Primary {
			s = fmt.Sprintf("\"%s\" INTEGER PRIMARY KEY AUTOINCREMENT", c.Name)
		} else {
			s = p.buildColumn(c)
		}
		cols = append(cols, s)
	}
	return fmt.Sprintf("CREATE TABLE \"%s\" (\n%s\n);", t.Name, strings.Join(cols, ",\n")), nil
}

func (p *SQLite) DropTable(t fizz.Table) (string, error) {
	return fmt.Sprintf("DROP TABLE \"%s\";", t.Name), nil
}

func (p *SQLite) RenameTable(t []fizz.Table) (string, error) {
	if len(t) < 2 {
		return "", errors.New("Not enough table names supplied!")
	}
	return fmt.Sprintf("ALTER TABLE \"%s\" RENAME TO \"%s\";", t[0].Name, t[1].Name), nil
}

func (p *SQLite) AddColumn(t fizz.Table) (string, error) {
	if len(t.Columns) == 0 {
		return "", errors.New("Not enough columns supplied!")
	}
	c := t.Columns[0]
	s := fmt.Sprintf("ALTER TABLE \"%s\" ADD COLUMN %s;", t.Name, p.buildColumn(c))
	return s, nil
}

func (p *SQLite) DropColumn(t fizz.Table) (string, error) {
	if len(t.Columns) < 1 {
		return "", errors.New("Not enough columns supplied!")
	}

	droppedColumn := t.Columns[0]

	rows, err := p.rows(t.Name)
	if err != nil {
		return "", err
	}
	if len(rows) == 0 {
		return "", errors.New("No table data was returned from SQLite!")
	}

	newTable := fizz.Table{
		Name:    t.Name,
		Columns: []fizz.Column{},
	}

	oldColumns := []string{}
	newColumns := []string{}
	for _, row := range rows {
		if row.Name != droppedColumn.Name {
			c := row.ToColumn()
			oldColumns = append(oldColumns, row.Name)
			newColumns = append(newColumns, c.Name)
			newTable.Columns = append(newTable.Columns, c)
		}
	}

	tempTable := fizz.Table{Name: fmt.Sprintf("%s_old", t.Name)}
	renameTableSQL, _ := p.RenameTable([]fizz.Table{t, tempTable})
	createTableSQL, _ := p.CreateTable(newTable)

	p.add(renameTableSQL)
	p.add(createTableSQL)
	p.add(fmt.Sprintf("INSERT INTO \"%s\" (%s) SELECT %s FROM \"%s\";", t.Name, strings.Join(newColumns, ", "), strings.Join(oldColumns, ", "), tempTable.Name))

	p.withRebuildIndexes(t, func() error {
		p.add(fmt.Sprintf("DROP TABLE \"%s\";", tempTable.Name))
		return nil
	})

	return strings.Join(p.sql, "\n"), nil
}

func (p *SQLite) withRebuildIndexes(t fizz.Table, fn func() error) error {
	indexes, err := p.indexes(t.Name)
	if err != nil {
		return err
	}

	for _, ix := range indexes {
		s, err := p.DropIndex(ix)
		if err != nil {
			return err
		}
		p.add(s)
	}

	err = fn()

	if err != nil {
		return err
	}

	for _, ix := range indexes {
		it := fizz.Table{
			Name:    t.Name,
			Indexes: []fizz.Index{ix},
		}
		s, err := p.AddIndex(it)
		if err != nil {
			return err
		}
		p.add(s)
	}

	return err
}

func (p *SQLite) RenameColumn(t fizz.Table) (string, error) {
	if len(t.Columns) < 2 {
		return "", errors.New("Not enough columns supplied!")
	}

	oldColumn := t.Columns[0]
	newColumn := t.Columns[1]

	rows, err := p.rows(t.Name)
	if err != nil {
		return "", err
	}
	if len(rows) == 0 {
		return "", errors.New("No table data was returned from SQLite!")
	}

	newTable := fizz.Table{
		Name:    t.Name,
		Columns: []fizz.Column{},
	}

	oldColumns := []string{}
	newColumns := []string{}
	for _, row := range rows {
		c := row.ToColumn()
		if row.Name == oldColumn.Name {
			c.Name = newColumn.Name
		}
		oldColumns = append(oldColumns, row.Name)
		newColumns = append(newColumns, c.Name)
		newTable.Columns = append(newTable.Columns, c)
	}

	tempTable := fizz.Table{Name: fmt.Sprintf("%s_old", t.Name)}
	renameTableSQL, _ := p.RenameTable([]fizz.Table{t, tempTable})
	createTableSQL, _ := p.CreateTable(newTable)

	p.add(renameTableSQL)
	p.add(createTableSQL)
	p.add(fmt.Sprintf("INSERT INTO \"%s\" (%s) SELECT %s FROM \"%s\";", t.Name, strings.Join(newColumns, ", "), strings.Join(oldColumns, ", "), tempTable.Name))

	p.withRebuildIndexes(t, func() error {
		p.add(fmt.Sprintf("DROP TABLE \"%s\";", tempTable.Name))
		return nil
	})

	return strings.Join(p.sql, "\n"), nil
}

func (p *SQLite) AddIndex(t fizz.Table) (string, error) {
	if len(t.Indexes) == 0 {
		return "", errors.New("Not enough indexes supplied!")
	}
	i := t.Indexes[0]
	s := fmt.Sprintf("CREATE INDEX \"%s\" ON \"%s\" (%s);", i.Name, t.Name, strings.Join(i.Columns, ", "))
	if i.Unique {
		s = strings.Replace(s, "CREATE", "CREATE UNIQUE", 1)
	}
	return s, nil
}

func (p *SQLite) DropIndex(i fizz.Index) (string, error) {
	return fmt.Sprintf("DROP INDEX \"%s\";", i.Name), nil
}

func (p *SQLite) RenameIndex(t fizz.Table) (string, error) {
	ix := t.Indexes
	if len(ix) < 2 {
		return "", errors.New("Not enough indexes supplied!")
	}
	oi := ix[0]
	ni := ix[1]
	return fmt.Sprintf("ALTER INDEX \"%s\" RENAME TO \"%s\";", oi.Name, ni.Name), nil
}

func (p *SQLite) buildColumn(c fizz.Column) string {
	s := fmt.Sprintf("\"%s\" %s", c.Name, p.colType(c))
	if c.Options["null"] == nil {
		s = fmt.Sprintf("%s NOT NULL", s)
	}
	if c.Options["default"] != nil {
		s = fmt.Sprintf("%s DEFAULT '%v'", s, c.Options["default"])
	}
	return s
}

func (p *SQLite) colType(c fizz.Column) string {
	switch c.ColType {
	case "timestamp":
		return "DATETIME"
	case "boolean", "DATE":
		return "NUMERIC"
	case "string":
		return "TEXT"
	default:
		return c.ColType
	}
}

func (p *SQLite) indexes(table string) ([]fizz.Index, error) {
	indexes := []fizz.Index{}
	db, err := sqlx.Open("sqlite3", p.URL)
	if err != nil {
		return indexes, err
	}
	defer db.Close()

	prag := fmt.Sprintf("PRAGMA index_list(%s)", table)
	res, err := db.Queryx(prag)
	if err != nil {
		return indexes, err
	}
	for res.Next() {
		li := sqliteIndexListInfo{}
		err = res.StructScan(&li)
		if err != nil {
			return indexes, err
		}

		i := fizz.Index{
			Name:    li.Name,
			Unique:  li.Unique,
			Columns: []string{},
		}

		prag = fmt.Sprintf("PRAGMA index_info(%s)", i.Name)
		iires, err := db.Queryx(prag)
		if err != nil {
			return indexes, err
		}

		for iires.Next() {
			ii := sqliteIndexInfo{}
			err = iires.StructScan(&ii)
			if err != nil {
				return indexes, err
			}
			i.Columns = append(i.Columns, ii.Name)
		}

		indexes = append(indexes, i)

	}
	return indexes, nil
}

func (p *SQLite) rows(table string) ([]sqliteTableInfo, error) {
	rows := []sqliteTableInfo{}

	db, err := sqlx.Open("sqlite3", p.URL)
	if err != nil {
		return rows, err
	}
	defer db.Close()

	prag := fmt.Sprintf("PRAGMA table_info(%s)", table)

	res, err := db.Queryx(prag)
	if err != nil {
		return rows, err
	}

	for res.Next() {
		ti := sqliteTableInfo{}
		err = res.StructScan(&ti)
		if err != nil {
			return rows, err
		}
		rows = append(rows, ti)
	}
	return rows, err
}
