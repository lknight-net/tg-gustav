package callback

import (
	"context"
	"log/slog"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

type Service interface {
	CallbackMessage(chatID int64, text string, data string) (string, error)
}

type Callback struct {
	CallbackService Service
	log             *slog.Logger
}

func (c *Callback) CallbackHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	txt, err := c.CallbackService.CallbackMessage(update.CallbackQuery.Message.Message.Chat.ID, update.CallbackQuery.Message.Message.Text, update.CallbackQuery.Data)
	if err != nil {
		c.log.Error("error callback message in service level", "error", err)
	}

	_, err = b.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{
		CallbackQueryID: update.CallbackQuery.ID,
		ShowAlert:       false,
	})
	if err != nil {
		c.log.Error("error answer callback query", "error", err)
	}

	_, err = b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.CallbackQuery.Message.Message.Chat.ID,
		Text:   txt,
	})
	if err != nil {
		c.log.Error("error send message", "error", err)
	}
}

func NewCallbackHandler(callbackService Service, log *slog.Logger) *Callback {
	return &Callback{CallbackService: callbackService, log: log}
}
