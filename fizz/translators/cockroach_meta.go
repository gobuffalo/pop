package translators

import (
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/markbates/pop/fizz"
)

type cockroachIndexListInfo struct {
	Name      string `db:"name"`
	NonUnique bool   `db:"non_unique"`
}

type cockroachIndexInfo struct {
	Name      string `db:"name"`
	Direction string `db:"direction"`
}

type cockroachTableInfo struct {
	Name       string      `db:"column_name"`
	Type       string      `db:"data_type"`
	IsNullable bool        `db:"is_nullable"`
	Default    interface{} `db:"column_default"`
	PK         bool        `db:"pk"`
}

func (t cockroachTableInfo) ToColumn() fizz.Column {
	c := fizz.Column{
		Name:    t.Name,
		ColType: t.Type,
		Primary: t.PK,
		Options: fizz.Options{},
	}
	if t.IsNullable {
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

	res, err := p.db.Queryx("SELECT table_name as name FROM information_schema.table;")
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
	prag := fmt.Sprintf("SELECT c.*, tc.table_schema IS NOT NULL AS \"pk\" FROM information_schema.columns AS c LEFT JOIN information_schema.key_column_usage as kcu ON ((c.table_schema = kcu.table_schema) AND (c.table_name = kcu.table_name) AND (c.column_name = kcu.column_name)) LEFT JOIN information_schema.table_constraints AS tc ON ((tc.table_schema = kcu.table_schema) AND (tc.table_name = kcu.table_name) AND (tc.constraint_name = kcu.constraint_name)) AND (tc.constraint_name = 'primary') WHERE c.table_name = '%s';", table.Name)

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
	prag := fmt.Sprintf("SELECT distinct index_name as name, non_unique FROM information_schema.statistics where table_name = '%s';", t.Name)
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
			Unique:  !li.NonUnique,
			Columns: []string{},
		}

		prag = fmt.Sprintf("SELECT column_name as name, direction FROM information_schema.statistics where index_name = '%s';", i.Name)
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
