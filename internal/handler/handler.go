package handler

import (
	"gustav-tgbot/internal/handler/command"
	"gustav-tgbot/internal/handler/file"
	"gustav-tgbot/internal/service"
	"log/slog"
)

type Handler struct {
	CommandHandler *command.Command
	FileHandler    *file.File
}

func NewHandler(service *service.Service, log *slog.Logger) *Handler {
	return &Handler{
		CommandHandler: command.NewCommandHandler(service.UserService, service.FileService, log),
		FileHandler:    file.NewFileHandler(service.FileService, service.UserService, log),
	}
}
