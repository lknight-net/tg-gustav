package handler

import (
	"gustav-tgbot/internal/handler/command"
	"gustav-tgbot/internal/service"
	"log/slog"
)

type Handler struct {
	CommandHandler *command.Command
}

func NewHandler(service *service.Service, log *slog.Logger) *Handler {
	return &Handler{
		CommandHandler: command.NewCommandHandler(service.UserService, log),
	}
}
