package repository

import (
	"database/sql"
	"gustav-tgbot/internal/repository/user"
	"log/slog"
)

type Repository struct {
	UserRepository *user.User
}

func NewRepository(db *sql.DB, log *slog.Logger) *Repository {
	return &Repository{
		UserRepository: user.NewUserRepository(db, log),
	}
}
