package fizz

type BubbleType int

const (
	E_CREATE_TABLE BubbleType = iota
	E_DROP_TABLE
	E_RAW_SQL
	E_ADD_COLUMN
)

type Bubble struct {
	Type BubbleType
	Data interface{}
}
