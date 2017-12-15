package translators

import (
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/markbates/pop/fizz"
)

type cockroachIndexListInfo struct {
	Seq     int    `db:"seq"`
	Name    string `db:"name"`
	Unique  bool   `db:"unique"`
	Origin  string `db:"origin"`
	Partial string `db:"partial"`
}

type cockroachIndexInfo struct {
	Seq  int    `db:"seqno"`
	CID  int    `db:"cid"`
	Name string `db:"name"`
}

type cockroachTableInfo struct {
	CID     int         `db:"cid"`
	Name    string      `db:"name"`
	Type    string      `db:"type"`
	NotNull bool        `db:"notnull"`
	Default interface{} `db:"dflt_value"`
	PK      bool        `db:"pk"`
}

func (t cockroachTableInfo) ToColumn() fizz.Column {
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

type cockroachSchema struct {
	Schema
}

func (p *cockroachSchema) Build() error {
	var err error
	p.db, err = sqlx.Open("postgres", p.URL)
	if err != nil {
		return err
	}
	defer p.db.Close()

	res, err := p.db.Queryx("SELECT table_name FROM information_schema.table;")
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
		if table.Name != "cockroach_sequence" {
			err = p.buildTableData(table)
			if err != nil {
				return err
			}
		}

	}
	return nil
}

func (p *cockroachSchema) buildTableData(table *fizz.Table) error {
	prag := fmt.Sprintf("PRAGMA table_info(%s)", table.Name)

	res, err := p.db.Queryx(prag)
	if err != nil {
		return nil
	}

	for res.Next() {
		ti := cockroachTableInfo{}
		err = res.StructScan(&ti)
		if err != nil {
			return err
		}
		table.Columns = append(table.Columns, ti.ToColumn())
	}
	err = p.buildTableIndexes(table)
	if err != nil {
		return err
	}
	p.schema[table.Name] = table
	return nil
}

func (p *cockroachSchema) buildTableIndexes(t *fizz.Table) error {
	prag := fmt.Sprintf("PRAGMA index_list(%s)", t.Name)
	res, err := p.db.Queryx(prag)
	if err != nil {
		return err
	}

	for res.Next() {
		li := cockroachIndexListInfo{}
		err = res.StructScan(&li)
		if err != nil {
			return err
		}

		i := fizz.Index{
			Name:    li.Name,
			Unique:  li.Unique,
			Columns: []string{},
		}

		prag = fmt.Sprintf("PRAGMA index_info(%s)", i.Name)
		iires, err := p.db.Queryx(prag)
		if err != nil {
			return err
		}

		for iires.Next() {
			ii := cockroachIndexInfo{}
			err = iires.StructScan(&ii)
			if err != nil {
				return err
			}
			i.Columns = append(i.Columns, ii.Name)
		}

		t.Indexes = append(t.Indexes, i)

	}
	return nil
}
