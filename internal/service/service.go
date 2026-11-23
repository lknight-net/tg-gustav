package service

import (
	"gustav-tgbot/internal/service/callback"
	"gustav-tgbot/internal/service/text"
	"log/slog"
)

type Service struct {
	TextService     *text.Text
	CallbackService *callback.Callback
}

func NewService(log *slog.Logger) *Service {
	return &Service{
		TextService:     text.NewTextService(log),
		CallbackService: callback.NewCallbackService(log),
	}
}
