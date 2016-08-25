package pop

import (
	"fmt"

	. "github.com/markbates/pop/columns"
	"github.com/markbates/pop/fizz"
)

type dialect interface {
	URL() string
	MigrationURL() string
	Details() *ConnectionDetails
	TranslateSQL(string) string
	Create(store, *Model, Columns) error
	Update(store, *Model, Columns) error
	Destroy(store, *Model) error
	SelectOne(store, *Model, Query) error
	SelectMany(store, *Model, Query) error
	CreateDB() error
	DropDB() error
	FizzTranslator() fizz.Translator
	Lock(func() error) error
}

func genericCreate(s store, model *Model, cols Columns) error {
	var id int64
	w := cols.Writeable()
	query := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)", model.TableName(), w.String(), w.SymbolizedString())
	Log(query)
	res, err := s.NamedExec(query, model.Value)
	if err != nil {
		return err
	}
	id, err = res.LastInsertId()
	if err == nil {
		model.setID(int(id))
	}
	return err
}

func genericUpdate(s store, model *Model, cols Columns) error {
	stmt := fmt.Sprintf("UPDATE %s SET %s where id = %d", model.TableName(), cols.Writeable().UpdateString(), model.ID())
	Log(stmt)
	_, err := s.NamedExec(stmt, model.Value)
	return err
}

func genericDestroy(s store, model *Model) error {
	stmt := fmt.Sprintf("DELETE FROM %s WHERE id = %d", model.TableName(), model.ID())
	return genericExec(s, stmt)
}

func genericExec(s store, stmt string) error {
	Log(stmt)
	_, err := s.Exec(stmt)
	return err
}

func genericSelectOne(s store, model *Model, query Query) error {
	sql, args := query.ToSQL(model)
	Log(sql, args...)
	err := s.Get(model.Value, sql, args...)
	return err
}

func genericSelectMany(s store, models *Model, query Query) error {
	sql, args := query.ToSQL(models)
	Log(sql, args...)
	err := s.Select(models.Value, sql, args...)
	return err
}
