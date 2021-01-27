package ac

import "context"

type User struct{}

func (u User) TableName(ctx context.Context) string {
	return ctx.Value("name").(string) + "_useras"
}
