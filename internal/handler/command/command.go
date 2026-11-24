package command

import (
	"context"
	"log/slog"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

type UserService interface {
	CreateUser(ctx context.Context, username *string, chatID int64) error
}

type Command struct {
	UserService UserService

	log *slog.Logger
}

func NewCommandHandler(userService UserService, log *slog.Logger) *Command {
	return &Command{
		UserService: userService,
		log:         log,
	}
}

func (c *Command) StartMessage(ctx context.Context, b *bot.Bot, update *models.Update) {
	chatID := update.Message.Chat.ID
	username := update.Message.Chat.Username
	if username == "" {
		c.log.Warn("Username is empty", "chat_id", chatID)
		err := c.UserService.CreateUser(ctx, nil, chatID)
		if err != nil {
			c.log.Error("Failed to create user", "chat_id", chatID)
		}
	}
	c.log.Info("New user created", "chat_id", chatID, "username", username)

	err := c.UserService.CreateUser(ctx, &username, chatID)
	if err != nil {
		c.log.Error("Failed to create user", "chat_id", chatID, "username", username)
	}
}
