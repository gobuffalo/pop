package translators

import (
	"errors"
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/markbates/pop/fizz"
	_ "github.com/mattes/migrate/driver/sqlite3"
)

type tableInfo struct {
	CID     int         `db:"cid"`
	Name    string      `db:"name"`
	Type    string      `db:"type"`
	NotNull bool        `db:"notnull"`
	Default interface{} `db:"dflt_value"`
	PK      bool        `db:"pk"`
}

func (t tableInfo) ToColumn() fizz.Column {
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
}

func (p SQLite) CreateTable(t fizz.Table) (string, error) {
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
	return fmt.Sprintf("CREATE TABLE IF NOT EXISTS \"%s\" (\n%s\n);", t.Name, strings.Join(cols, ",\n")), nil
}

func (p SQLite) DropTable(t fizz.Table) (string, error) {
	return fmt.Sprintf("DROP TABLE IF EXISTS \"%s\";", t.Name), nil
}

func (p SQLite) RenameTable(t []fizz.Table) (string, error) {
	if len(t) < 2 {
		return "", errors.New("Not enough table names supplied!")
	}
	return fmt.Sprintf("ALTER TABLE \"%s\" RENAME TO \"%s\";", t[0].Name, t[1].Name), nil
}

func (p SQLite) AddColumn(t fizz.Table) (string, error) {
	if len(t.Columns) == 0 {
		return "", errors.New("Not enough columns supplied!")
	}
	c := t.Columns[0]
	s := fmt.Sprintf("ALTER TABLE \"%s\" ADD COLUMN %s;", t.Name, p.buildColumn(c))
	return s, nil
}

func (p SQLite) DropColumn(t fizz.Table) (string, error) {
	return "", errors.New("NEEDS TO BE IMPLEMENTED!")
	// if len(t.Columns) == 0 {
	// 	return "", errors.New("Not enough columns supplied!")
	// }
	// c := t.Columns[0]
	// return fmt.Sprintf("ALTER TABLE \"%s\" DROP COLUMN \"%s\";", t.Name, c.Name), nil
}

func (p SQLite) RenameColumn(t fizz.Table) (string, error) {
	if len(t.Columns) < 2 {
		return "", errors.New("Not enough columns supplied!")
	}

	oldColumn := t.Columns[0]
	newColumn := t.Columns[1]

	rows, err := p.rows(t.Name)
	if err != nil {
		return "", err
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

	sql := []string{
		renameTableSQL,
		createTableSQL,
		fmt.Sprintf("INSERT INTO \"%s\" (%s) SELECT %s FROM \"%s\";", t.Name, strings.Join(newColumns, ", "), strings.Join(oldColumns, ", "), tempTable.Name),
	}
	sql = append(sql, fmt.Sprintf("DROP TABLE \"%s\";", tempTable.Name))
	return strings.Join(sql, "\n"), nil
}

func (p SQLite) AddIndex(t fizz.Table) (string, error) {
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

func (p SQLite) DropIndex(i fizz.Index) (string, error) {
	return fmt.Sprintf("DROP INDEX IF EXISTS \"%s\";", i.Name), nil
}

func (p SQLite) RenameIndex(ix []fizz.Index) (string, error) {
	return "", errors.New("NEEDS TO BE IMPLEMENTED!")
	// if len(ix) < 2 {
	// 	return "", errors.New("Not enough indexes supplied!")
	// }
	// oi := ix[0]
	// ni := ix[1]
	// return fmt.Sprintf("ALTER INDEX \"%s\" RENAME TO \"%s\";", oi.Name, ni.Name), nil
}

func (p SQLite) buildColumn(c fizz.Column) string {
	s := fmt.Sprintf("\"%s\" %s", c.Name, p.colType(c))
	if c.Options["null"] == nil {
		s = fmt.Sprintf("%s NOT NULL", s)
	}
	if c.Options["default"] != nil {
		s = fmt.Sprintf("%s DEFAULT '%v'", s, c.Options["default"])
	}
	return s
}

func (p SQLite) colType(c fizz.Column) string {
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

func (p SQLite) rows(table string) ([]tableInfo, error) {
	rows := []tableInfo{}

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
		ti := tableInfo{}
		err = res.StructScan(&ti)
		if err != nil {
			return rows, err
		}
		rows = append(rows, ti)
	}
	return rows, err
}
