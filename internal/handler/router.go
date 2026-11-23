package handler

import (
	"gustav-tgbot/internal/service"
	"log/slog"

	"github.com/go-telegram/bot"
)

type Router struct {
	B *bot.Bot
}

func NewRouter(b *bot.Bot, service *service.Service, log *slog.Logger) *Router {
	handlers := NewHandler(service.TextService, service.CallbackService, log)

	b.RegisterHandler(bot.HandlerTypeMessageText, "/start", bot.MatchTypeExact, handlers.TextHandler.StartMessage)
	b.RegisterHandler(bot.HandlerTypeCallbackQueryData, "button", bot.MatchTypePrefix, handlers.CallbackHandler.CallbackHandler)

	return &Router{
		B: b,
	}
}
