package pop

import "reflect"

func (c *Connection) Find(model interface{}, id int) error {
	q := Q(c).Where("id = ?", id)
	return q.First(model)
}

func (c *Connection) First(model interface{}) error {
	return Q(c).First(model)
}

func (q *Query) First(model interface{}) error {
	return q.Connection.timeFunc("First", func() error {
		q.Limit(1)
		m := &Model{Value: model}
		return q.Connection.Dialect.SelectOne(q.Connection.Store, m, *q)
	})
}

func (c *Connection) All(models interface{}) error {
	return Q(c).All(models)
}

func (q *Query) All(models interface{}) error {
	return q.Connection.timeFunc("All", func() error {
		m := &Model{Value: models}
		err := q.Connection.Dialect.SelectMany(q.Connection.Store, m, *q)
		if err == nil && q.Paginator != nil {
			ct, err := q.Count(models)
			if err == nil {
				q.Paginator.TotalEntriesSize = ct
				st := reflect.ValueOf(models).Elem()
				q.Paginator.CurrentEntriesSize = st.Len()
				q.Paginator.TotalPages = (q.Paginator.TotalEntriesSize / q.Paginator.PerPage)
				if q.Paginator.TotalEntriesSize%q.Paginator.PerPage > 0 {
					q.Paginator.TotalPages = q.Paginator.TotalPages + 1
				}
			}
		}
		return err
	})
}

func (c *Connection) Exists(model interface{}) (bool, error) {
	return Q(c).Exists(model)
}

func (q *Query) Exists(model interface{}) (bool, error) {
	i, err := q.Count(model)
	return i != 0, err
}

func (c *Connection) Count(model interface{}) (int, error) {
	return Q(c).Count(model)
}

func (q Query) Count(model interface{}) (int, error) {
	res := &rowCount{}
	err := q.Connection.timeFunc("Count", func() error {
		q.Paginator = nil
		col := "count(*) as row_count"
		query, args := q.ToSQL(&Model{Value: model}, col)
		return q.Connection.Store.Get(res, query, args...)
	})
	return res.Count, err
}

type rowCount struct {
	Count int `db:"row_count"`
}
