package callback

import (
	"fmt"
	"log/slog"
)

type Callback struct {
	log *slog.Logger
}

func (c *Callback) CallbackMessage(chatID int64, text string, data string) (string, error) {
	return fmt.Sprintf("button: %s", data), nil
}

func NewCallbackService(log *slog.Logger) *Callback {
	return &Callback{
		log: log,
	}
}
