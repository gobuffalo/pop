package fizz

func init() {
	fizzers["raw"] = RawSQL
}

func RawSQL(ch chan *Bubble) interface{} {
	return func(sql string) {
		ch <- &Bubble{BubbleType: E_RAW_SQL, Data: sql}
	}
}
