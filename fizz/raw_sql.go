package fizz

func (f fizzer) RawSql() interface{} {
	return func(sql string) {
		f.add(Bubble{BubbleType: E_RAW_SQL, Data: sql})
	}
}
