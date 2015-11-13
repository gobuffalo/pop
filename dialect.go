package pop

import (
	"fmt"

	. "github.com/markbates/pop/columns"
)

type Dialect interface {
	URL() string
	Details() *ConnectionDetails
	TranslateSQL(sql string) string
	Create(store Store, model *Model, cols Columns) error
	Update(store Store, model *Model, cols Columns) error
	Destroy(store Store, model *Model) error
	SelectOne(store Store, model *Model, query Query) error
	SelectMany(store Store, models *Model, query Query) error
	CreateDB() error
	DropDB() error
}

func genericCreate(store Store, model *Model, cols Columns) error {
	var id int64
	w := cols.Writeable()
	query := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)", model.TableName(), w.String(), w.SymbolizedString())
	if Debug {
		Log(query)
	}
	res, err := store.NamedExec(query, model.Value)
	if err != nil {
		return err
	}
	id, err = res.LastInsertId()
	if err == nil {
		model.SetID(int(id))
	}
	return err
}

func genericUpdate(store Store, model *Model, cols Columns) error {
	stmt := fmt.Sprintf("UPDATE %s SET %s where id = %d", model.TableName(), cols.Writeable().UpdateString(), model.ID())
	if Debug {
		Log(stmt)
	}
	_, err := store.NamedExec(stmt, model.Value)
	if err != nil {
		Log(fmt.Sprintf("Error: Update: %q", err))
	}
	return err
}

func genericDestroy(store Store, model *Model) error {
	stmt := fmt.Sprintf("DELETE FROM %s WHERE id = %d", model.TableName(), model.ID())
	if Debug {
		Log(stmt)
	}
	_, err := store.Exec(stmt)
	return err
}

func genericSelectOne(store Store, model *Model, query Query) error {
	sql, args := query.ToSQL(model)
	err := store.Get(model.Value, sql, args...)
	return err
}

func genericSelectMany(store Store, models *Model, query Query) error {
	sql, args := query.ToSQL(models)
	err := store.Select(models.Value, sql, args...)
	return err
}
