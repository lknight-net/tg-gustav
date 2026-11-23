package text

import (
	"context"
	"log/slog"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

type Service interface {
	StartMessage(chatID int64) (string, error)
}
type Text struct {
	TextService Service
	log         *slog.Logger
}

func (t *Text) StartMessage(ctx context.Context, b *bot.Bot, update *models.Update) {
	chatID := update.Message.Chat.ID
	response, err := t.TextService.StartMessage(chatID)
	if err != nil {
		t.log.Error("error starting message in service level", "error", err)
	}

	kb := &models.InlineKeyboardMarkup{
		InlineKeyboard: [][]models.InlineKeyboardButton{
			{
				{Text: "Button 1", CallbackData: "button_1"},
				{Text: "Button 2", CallbackData: "button_2"},
			}, {
				{Text: "Button 3", CallbackData: "button_3"},
			},
		},
	}

	_, err = b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:      chatID,
		Text:        response,
		ReplyMarkup: kb,
	})
	if err != nil {
		t.log.Error("error send message", "error", err)
	}
}
func NewTextHandler(textService Service, log *slog.Logger) *Text {
	return &Text{TextService: textService, log: log}
}
