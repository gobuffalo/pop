package a

type User struct{}

func (u User) TableName() string {
	return "usera"
}
