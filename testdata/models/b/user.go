package b

type User struct{}

func (u User) TableName() string {
	return "userb"
}
