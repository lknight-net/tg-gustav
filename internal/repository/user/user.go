package user

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"gustav-tgbot/internal/model"
	"log/slog"

	"modernc.org/sqlite"
)

const sqliteConstraintUnique = 2067

type User struct {
	db *sql.DB

	log *slog.Logger
}

func (u *User) UpdateStateUserByUserID(ctx context.Context, userID int64, state string) error {
	query := `UPDATE users SET state = ? WHERE id = ?`

	result, err := u.db.ExecContext(ctx, query, state, userID)
	if err != nil {
		return fmt.Errorf("update user: %w", err)
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("get affected rows: %w", err)
	}
	if affected == 0 {
		return ErrUserNotFound
	}
	return nil
}

func (u *User) GetUserByChatID(ctx context.Context, chatID int64) (*model.User, error) {
	query := `SELECT id, username, chat_id, state FROM users WHERE chat_id = ?`

	var user model.User

	err := u.db.QueryRowContext(ctx, query, chatID).Scan(
		&user.ID,
		&user.UserName,
		&user.ChatID,
		&user.StateUser,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrUserNotFound
		}

		return nil, fmt.Errorf("get user by chat id: %w", err)
	}

	return &user, nil
}

func (u *User) CreateUser(ctx context.Context, user *model.User) (int64, error) {
	query := `
        INSERT INTO users (chat_id, username, state)
        VALUES (?, ?, ?)
    `

	result, err := u.db.ExecContext(ctx, query, user.ChatID, user.UserName, user.StateUser)
	if err != nil {
		var sqliteErr *sqlite.Error

		if errors.As(err, &sqliteErr) && sqliteErr.Code() == sqliteConstraintUnique {
			return 0, ErrUserAlreadyExist
		}

		return 0, fmt.Errorf("create user: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("get last inserted user ID: %w", err)
	}

	return id, nil
}

func NewUserRepository(db *sql.DB, log *slog.Logger) *User {
	return &User{db: db, log: log}
}
