package user

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"gustav-tgbot/internal/model"
	"log/slog"
)

type User struct {
	db *sql.DB

	log *slog.Logger
}

func (u *User) CreateUser(ctx context.Context, user *model.User) (int64, error) {
	query := `
        INSERT INTO users (chat_id, username)
        VALUES (?, ?)
    `

	result, err := u.db.ExecContext(ctx, query, user.ChatID, user.UserName)
	if err != nil {
		return 0, fmt.Errorf("create user: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("get inserted user ID: %w", err)
	}

	return id, nil
}

func (u *User) GetUserByChatID(ctx context.Context, chatID int64) (*model.User, error) {
	query := `
        SELECT id, chat_id, username
        FROM users
        WHERE chat_id = ?
    `

	var user model.User

	err := u.db.QueryRowContext(ctx, query, chatID).Scan(
		&user.ID,
		&user.ChatID,
		&user.UserName,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		return nil, fmt.Errorf("get user by chat ID: %w", err)
	}

	return &user, nil
}

func (u *User) GetUserByUsername(ctx context.Context, username string) (*model.User, error) {
	query := `
        SELECT id, chat_id, username
        FROM users
        WHERE username = ?
    `

	var user model.User

	err := u.db.QueryRowContext(ctx, query, username).Scan(
		&user.ID,
		&user.ChatID,
		&user.UserName,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		return nil, fmt.Errorf("get user by username: %w", err)
	}

	return &user, nil
}

func (u *User) GetUserByID(ctx context.Context, id int64) (*model.User, error) {
	query := `
        SELECT id, chat_id, username
        FROM users
        WHERE id = ?
    `

	var user model.User

	err := u.db.QueryRowContext(ctx, query, id).Scan(
		&user.ID,
		&user.ChatID,
		&user.UserName,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		return nil, fmt.Errorf("get user by ID: %w", err)
	}

	return &user, nil
}

func NewUserRepository(db *sql.DB, log *slog.Logger) *User {
	return &User{db: db, log: log}
}
