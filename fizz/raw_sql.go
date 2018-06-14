package fizz

// RawSql executes a raw SQL statement.
//
// Deprecated: use RawSQL instead.
func (f fizzer) RawSql() interface{} {
	return f.RawSQL()
}

// RawSQL executes a raw SQL statement.
func (f fizzer) RawSQL() interface{} {
	return func(sql string) {
		f.add(sql, nil)
	}
}
