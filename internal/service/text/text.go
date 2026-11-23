package text

import (
	"fmt"
	"log/slog"
)

type Text struct {
	log *slog.Logger
}

func (t *Text) StartMessage(chatID int64) (string, error) {
	str := fmt.Sprintf(
		"Добро пожаловать в тестового бота № 0, пользователь № %d"+
			"\nПожалуйста, воспользуйтесь кнопками ниже для понимания, что умеет этот бот", chatID)
	return str, nil
}

func NewTextService(log *slog.Logger) *Text {
	return &Text{
		log: log,
	}
}
