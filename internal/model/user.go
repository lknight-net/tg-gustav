package model

type User struct {
	ID       int64
	UserName *string
	ChatID   int64
}

func NewUser(username *string, chatID int64) *User {
	return &User{
		UserName: username,
		ChatID:   chatID,
	}
}
