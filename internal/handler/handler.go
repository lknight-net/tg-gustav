package handler

import (
	"context"
	"gustav-tgbot/internal/handler/callback"
	"gustav-tgbot/internal/handler/text"
	"log/slog"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

type TextHandler interface {
	StartMessage(ctx context.Context, b *bot.Bot, update *models.Update)
}

type CallbackHandler interface {
	CallbackHandler(ctx context.Context, b *bot.Bot, update *models.Update)
}

type Handler struct {
	TextHandler     TextHandler
	CallbackHandler CallbackHandler
}

func NewHandler(textService text.Service, callbackService callback.Service, log *slog.Logger) *Handler {
	return &Handler{
		TextHandler:     text.NewTextHandler(textService, log),
		CallbackHandler: callback.NewCallbackHandler(callbackService, log),
	}
}
