package translators

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"

	"github.com/markbates/pop/fizz"
)

type Postgres struct {
}

func NewPostgres() *Postgres {
	return &Postgres{}
}

func (p *Postgres) CreateTable(t fizz.Table) (string, error) {
	sql := []string{}
	cols := []string{}
	var s string
	for _, c := range t.Columns {
		if c.Primary {
			switch c.ColType {
			case "string", "uuid":
				s = fmt.Sprintf("\"%s\" %s PRIMARY KEY", c.Name, p.colType(c))
			case "integer":
				s = fmt.Sprintf("\"%s\" SERIAL PRIMARY KEY", c.Name)
			default:
				return "", errors.Errorf("can not use %s as a primary key", c.ColType)
			}
		} else {
			s = p.buildColumn(c)
		}
		cols = append(cols, s)
	}
	s = fmt.Sprintf("CREATE TABLE \"%s\" (\n%s\n);", t.Name, strings.Join(cols, ",\n"))
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

func (p *Postgres) DropTable(t fizz.Table) (string, error) {
	return fmt.Sprintf("DROP TABLE \"%s\";", t.Name), nil
}

func (p *Postgres) RenameTable(t []fizz.Table) (string, error) {
	if len(t) < 2 {
		return "", errors.New("Not enough table names supplied!")
	}
	return fmt.Sprintf("ALTER TABLE \"%s\" RENAME TO \"%s\";", t[0].Name, t[1].Name), nil
}

func (p *Postgres) AddColumn(t fizz.Table) (string, error) {
	if len(t.Columns) == 0 {
		return "", errors.New("Not enough columns supplied!")
	}
	c := t.Columns[0]
	s := fmt.Sprintf("ALTER TABLE \"%s\" ADD COLUMN %s;", t.Name, p.buildColumn(c))
	return s, nil
}

func (p *Postgres) DropColumn(t fizz.Table) (string, error) {
	if len(t.Columns) == 0 {
		return "", errors.New("Not enough columns supplied!")
	}
	c := t.Columns[0]
	return fmt.Sprintf("ALTER TABLE \"%s\" DROP COLUMN \"%s\";", t.Name, c.Name), nil
}

func (p *Postgres) RenameColumn(t fizz.Table) (string, error) {
	if len(t.Columns) < 2 {
		return "", errors.New("Not enough columns supplied!")
	}
	oc := t.Columns[0]
	nc := t.Columns[1]
	s := fmt.Sprintf("ALTER TABLE \"%s\" RENAME COLUMN \"%s\" TO \"%s\";", t.Name, oc.Name, nc.Name)
	return s, nil
}

func (p *Postgres) AddIndex(t fizz.Table) (string, error) {
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

func (p *Postgres) DropIndex(t fizz.Table) (string, error) {
	if len(t.Indexes) == 0 {
		return "", errors.New("Not enough indexes supplied!")
	}
	i := t.Indexes[0]
	return fmt.Sprintf("DROP INDEX \"%s\";", i.Name), nil
}

func (p *Postgres) RenameIndex(t fizz.Table) (string, error) {
	ix := t.Indexes
	if len(ix) < 2 {
		return "", errors.New("Not enough indexes supplied!")
	}
	oi := ix[0]
	ni := ix[1]
	return fmt.Sprintf("ALTER INDEX \"%s\" RENAME TO \"%s\";", oi.Name, ni.Name), nil
}

func (p *Postgres) buildColumn(c fizz.Column) string {
	s := fmt.Sprintf("\"%s\" %s", c.Name, p.colType(c))
	if c.Options["null"] == nil {
		s = fmt.Sprintf("%s NOT NULL", s)
	}
	if c.Options["default"] != nil {
		s = fmt.Sprintf("%s DEFAULT '%v'", s, c.Options["default"])
	}
	return s
}

func (p *Postgres) colType(c fizz.Column) string {
	switch c.ColType {
	case "string":
		s := "255"
		if c.Options["size"] != nil {
			s = fmt.Sprintf("%d", c.Options["size"])
		}
		return fmt.Sprintf("VARCHAR (%s)", s)
	case "uuid":
		return "UUID"
	default:
		return c.ColType
	}
}
