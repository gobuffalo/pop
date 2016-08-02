package fizz

type BubbleType int

const (
	E_CREATE_TABLE BubbleType = iota
	E_DROP_TABLE
	E_RENAME_TABLE
	E_RAW_SQL
	E_ADD_COLUMN
	E_DROP_COLUMN
	E_RENAME_COLUMN
	E_ADD_INDEX
	E_DROP_INDEX
	E_RENAME_INDEX
)

type Bubble struct {
	BubbleType BubbleType
	Data       interface{}
}

type Bubbler struct {
	Bubbles chan *Bubble
}
