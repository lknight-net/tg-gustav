package handler

import (
	"gustav-tgbot/internal/service"
	"log/slog"

	"github.com/go-telegram/bot"
)

type Router struct {
	log *slog.Logger
	B   *bot.Bot
}

func NewRouter(service *service.Service, b *bot.Bot, log *slog.Logger) *Router {
	handler := NewHandler(service, log)
	b.RegisterHandler(bot.HandlerTypeMessageText, "/start", bot.MatchTypeExact, handler.CommandHandler.StartMessage)
	return &Router{log: log, B: b}
}
