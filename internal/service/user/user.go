package user

import (
	"context"
	"database/sql"
	"errors"
	"gustav-tgbot/internal/model"
	"gustav-tgbot/internal/repository/user"
	"log/slog"
)

type UserRepository interface {
	CreateUser(ctx context.Context, user *model.User) (int64, error)
	UpdateStateUserByUserID(ctx context.Context, userID int64, state string) error
	GetUserByChatID(ctx context.Context, chatID int64) (*model.User, error)
}

type User struct {
	UserRepository UserRepository

	log *slog.Logger
}

func (u *User) GetUserByChatID(ctx context.Context, chatID int64) (*model.User, error) {
	user, err := u.UserRepository.GetUserByChatID(ctx, chatID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		u.log.Error("failed get user by chat id", "chat_id", chatID, "err", err)
		return nil, err
	}
	return user, err
}

func (u *User) UpdateStateUserByUserID(ctx context.Context, userID int64, state string) error {
	err := u.UserRepository.UpdateStateUserByUserID(ctx, userID, state)
	if err != nil {
		u.log.Error("failed to update state user by user_id", "user_id", userID, "state", state, "err", err)
		return err
	}
	u.log.Info("state user updated", "user_id", userID, "state", state)
	return nil
}

func (u *User) CreateUser(ctx context.Context, username *string, chatID int64) (int64, error) {
	id, err := u.UserRepository.CreateUser(ctx, model.NewUser(username, chatID))
	if err != nil {
		if errors.Is(err, user.ErrUserAlreadyExist) {
			return 0, ErrUserAlreadyExist
		}

		u.log.Error("failed to create user", "chat_id", chatID, "err", err)
		return 0, err
	}
	u.log.Info("user is created", "chat_id", chatID, "id", id)
	return id, nil
}

func NewUserService(userRepository UserRepository, log *slog.Logger) *User {
	return &User{
		UserRepository: userRepository,
		log:            log,
	}
}
