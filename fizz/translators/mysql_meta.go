package translators

import (
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/markbates/pop/fizz"
	"github.com/pkg/errors"
)

type mysqlTableInfo struct {
	Field   string      `db:"Field"`
	Type    string      `db:"Type"`
	Null    string      `db:"Null"`
	Key     string      `db:"Key"`
	Default interface{} `db:"Default"`
	Extra   string      `db:"Extra"`
}

type mysqlIndexListInfo struct {
	NonUnique  bool   `db:"non_unique"`
	IndexName  string `db:"index_name"`
	ColumnName string `db:"column_name"`
}

func (ti mysqlTableInfo) ToColumn() fizz.Column {
	c := fizz.Column{
		Name:    ti.Field,
		ColType: ti.Type,
		Primary: ti.Key == "PRI",
		Options: map[string]interface{}{},
	}
	if strings.ToLower(ti.Null) == "yes" {
		c.Options["null"] = true
	}
	if ti.Default != nil {
		d := fmt.Sprintf("%s", ti.Default)
		c.Options["default"] = d
	}
	return c
}

type mysqlSchema struct {
	Schema
}

func (p *mysqlSchema) Version() (string, error) {
	var version string
	var err error

	p.db, err = sqlx.Open("mysql", p.URL)
	if err != nil {
		return version, err
	}
	defer p.db.Close()

	res, err := p.db.Queryx("select VERSION()")
	if err != nil {
		return version, err
	}

	for res.Next() {
		err = res.Scan(&version)
		return version, err
	}
	return "", errors.New("could not locate MySQL version")
}

func (p *mysqlSchema) Build() error {
	var err error
	p.db, err = sqlx.Open("mysql", p.URL)
	if err != nil {
		return err
	}
	defer p.db.Close()

	res, err := p.db.Queryx(fmt.Sprintf("select TABLE_NAME as name from information_schema.TABLES where TABLE_SCHEMA = '%s'", p.Name))
	if err != nil {
		return err
	}
	for res.Next() {
		table := &fizz.Table{
			Columns: []fizz.Column{},
			Indexes: []fizz.Index{},
		}
		err = res.StructScan(table)
		if err != nil {
			return err
		}
		err = p.buildTableData(table)
		if err != nil {
			return err
		}
		err = p.buildTableIndexes(table)
		if err != nil {
			return err
		}
	}
	return nil
}

func (p *mysqlSchema) buildTableData(table *fizz.Table) error {
	prag := fmt.Sprintf("describe %s", table.Name)

	res, err := p.db.Queryx(prag)
	if err != nil {
		return nil
	}

	for res.Next() {
		ti := mysqlTableInfo{}
		err = res.StructScan(&ti)
		if err != nil {
			return err
		}
		table.Columns = append(table.Columns, ti.ToColumn())
	}

	p.schema[table.Name] = table
	return nil
}

func (p *mysqlSchema) buildTableIndexes(t *fizz.Table) error {
	indexes := map[string]fizz.Index{}

	prag := fmt.Sprintf("SELECT non_unique, index_name, column_name FROM INFORMATION_SCHEMA.STATISTICS WHERE table_name = '%s' order by seq_in_index", t.Name)
	res, err := p.db.Queryx(prag)
	if err != nil {
		return err
	}

	for res.Next() {
		li := mysqlIndexListInfo{}
		err = res.StructScan(&li)
		if err != nil {
			return err
		}

		i, ok := indexes[li.IndexName]
		if !ok {
			i := fizz.Index{
				Name:    li.IndexName,
				Unique:  !li.NonUnique,
				Columns: []string{},
			}
			indexes[li.IndexName] = i
		}

		i.Columns = append(i.Columns, li.ColumnName)
		indexes[li.IndexName] = i
	}

	t.Indexes = []fizz.Index{}
	for _, i := range indexes {
		t.Indexes = append(t.Indexes, i)
	}

	return nil
}
