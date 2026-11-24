package user

import (
	"context"
	"gustav-tgbot/internal/model"
	"log/slog"
)

type UserRepository interface {
	CreateUser(ctx context.Context, user *model.User) (int64, error)
	GetUserByChatID(ctx context.Context, chatID int64) (*model.User, error)
	GetUserByUsername(ctx context.Context, username string) (*model.User, error)
	GetUserByID(ctx context.Context, id int64) (*model.User, error)
}

type User struct {
	UserRepository UserRepository

	log *slog.Logger
}

func (u *User) CreateUser(ctx context.Context, username *string, chatID int64) error {
	id, err := u.UserRepository.CreateUser(ctx, model.NewUser(username, chatID))
	if err != nil {
		u.log.Error("CreateUser", "err", err)
	}
	u.log.Info("CreateUser", "id", id)
	return nil
}

func NewUserService(userRepository UserRepository, log *slog.Logger) *User {
	return &User{
		UserRepository: userRepository,
		log:            log,
	}
}
