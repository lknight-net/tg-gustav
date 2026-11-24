package service

import (
	"gustav-tgbot/internal/repository"
	"gustav-tgbot/internal/service/user"
	"log/slog"
)

type Service struct {
	UserService *user.User
}

func NewService(repository *repository.Repository, log *slog.Logger) *Service {
	return &Service{
		UserService: user.NewUserService(repository.UserRepository, log),
	}
}
