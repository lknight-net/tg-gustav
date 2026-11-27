package command

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"gustav-tgbot/internal/model"
	userErr "gustav-tgbot/internal/service/user"
	"log/slog"
	"slices"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

type UserService interface {
	CreateUser(ctx context.Context, username *string, chatID int64) (int64, error)
	UpdateStateUserByUserID(ctx context.Context, userID int64, state string) error
	GetUserByChatID(ctx context.Context, chatID int64) (*model.User, error)
}
type FileService interface {
	GetFilesByUserID(ctx context.Context, userID int64) (map[int64][]byte, error)
}

type Command struct {
	UserService UserService
	FileService FileService

	log *slog.Logger
}

func NewCommandHandler(userService UserService, fileService FileService, log *slog.Logger) *Command {
	return &Command{
		UserService: userService,
		FileService: fileService,
		log:         log,
	}
}

func (c *Command) StartMessage(ctx context.Context, b *bot.Bot, update *models.Update) {
	chatID := update.Message.Chat.ID
	username := update.Message.Chat.Username
	var id int64
	if username == "" {
		c.log.Warn("username is empty", "command", "/start", "chat_id", chatID)
		idUser, err := c.UserService.CreateUser(ctx, nil, chatID)
		if err != nil {
			c.log.Error("failed to create user", "command", "/start", "chat_id", chatID, "username", username, "error", err)
			return
		}
		id = idUser
	}

	idUser, err := c.UserService.CreateUser(ctx, &username, chatID)
	if err != nil {
		if errors.Is(err, userErr.ErrUserAlreadyExist) {
			_, err := b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: chatID,
				Text:   "Вы уже есть в системе.",
			})
			if err != nil {
				c.log.Error("failed to send message", "command", "/start", "chat_id", chatID, "error", err)
				return
			}
		}

		c.log.Error("failed to create user", "command", "/start", "chat_id", chatID, "username", username, "error", err)
		return
	}
	id = idUser

	c.log.Info("New user created", "command", "/start", "chat_id", chatID, "username", username)

	_, err = b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: chatID,
		Text:   fmt.Sprintf("Добро пожаловать. Ваш программный идентификатор № %d", id),
	})
	if err != nil {
		c.log.Error("failed to send message", "command", "/start", "chat_id", chatID, "error", err)
		return
	}
}

func (c *Command) WatermarkMessage(ctx context.Context, b *bot.Bot, update *models.Update) {
	chatID := update.Message.Chat.ID

	user, err := c.UserService.GetUserByChatID(ctx, chatID)
	if err != nil {
		if errors.Is(err, userErr.ErrUserNotFound) {
			_, err := b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: chatID,
				Text:   "Регистрация не прошла. Напишите /start для начала.",
			})
			if err != nil {
				c.log.Error("failed to send message", "command", "/watermark", "chat_id", chatID, "error", err)
				return
			}
		}
		c.log.Error("failed to get user by chat id", "command", "/watermark", "chat_id", chatID, "error", err)
		return
	}

	if user.StateUser != model.StateInactivity {
		c.log.Info("user is not in inactivity state", "command", "/watermark", "chat_id", chatID)
		_, err := b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: chatID,
			Text:   "Вы не нажали /cancel."})
		if err != nil {
			c.log.Error("failed to send message", "command", "/watermark", "chat_id", chatID, "error", err)
			return
		}
		return
	}

	err = c.UserService.UpdateStateUserByUserID(ctx, user.ID, model.StateSendFile)
	if err != nil {
		c.log.Error("failed to update state user", "command", "/watermark", "chat_id", chatID, "error", err)
		return
	}

	_, err = b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: chatID,
		Text: "Пожалуйста, пришлите одну фотографию или документ для удаления водяного знака.\n" +
			"Важно: если Вы присылаете документы, то расширение файлов должно быть PNG, JPEG или JPG. В случае отправки WEBP программа отработает нестандартно(для теста надо попробовать).\n" +
			"Важно: обработка каждой фотографии не мгновенная, и Вам надо немного подождать. Вы можете следить за прогрессом изменения фотографии в меню /menu (не разработано).\n" +
			"Важно: если Вы закончили отправлять фотографии, напиши команду: /cancel (не разработано).",
	})
	if err != nil {
		c.log.Error("failed to send message", "command", "/watermark", "chat_id", chatID, "error", err)
		return
	}
}

func (c *Command) CancelMessage(ctx context.Context, b *bot.Bot, update *models.Update) {
	chatID := update.Message.Chat.ID

	user, err := c.UserService.GetUserByChatID(ctx, chatID)
	if err != nil {
		if errors.Is(err, userErr.ErrUserNotFound) {
			_, err := b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: chatID,
				Text:   "Регистрация не прошла. Напишите /start для начала.",
			})
			if err != nil {
				c.log.Error("failed to send message", "command", "/cancel", "chat_id", chatID, "error", err)
				return
			}
		}
		c.log.Error("failed to get user by chat id", "command", "/cancel", "chat_id", chatID, "error", err)
		return
	}

	if user.StateUser != model.StateSendFile {
		c.log.Info("user is in inactivity state", "command", "/cancel", "chat_id", chatID)
		_, err := b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: chatID,
			Text:   "Вы не нажали /watermark."})
		if err != nil {
			c.log.Error("failed to send message", "command", "/cancel", "chat_id", chatID, "error", err)
			return
		}
		return
	}

	err = c.UserService.UpdateStateUserByUserID(ctx, user.ID, model.StateInactivity)
	if err != nil {
		c.log.Error("failed to update state user", "command", "/cancel", "chat_id", chatID, "error", err)
		return
	}

	_, err = b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: chatID,
		Text:   "Вы можете ознакомиться с результатом /menu.(не разработано)",
	})
	if err != nil {
		c.log.Error("failed to send message", "command", "/cancel", "chat_id", chatID, "error", err)
		return
	}

	filesMap, err := c.FileService.GetFilesByUserID(ctx, user.ID)
	if err != nil {
		c.log.Error("failed to get files", "command", "/cancel", "chat_id", chatID, "error", err)
		return
	}

	// Собираем даты.
	times := make([]int64, 0, len(filesMap))
	for t := range filesMap {
		times = append(times, t)
	}

	// Сортируем от старых к новым.
	slices.Sort(times)

	// Обрабатываем файлы в FIFO-порядке.
	for _, t := range times {
		fileBytes := filesMap[t]
		_, err = b.SendPhoto(ctx, &bot.SendPhotoParams{
			ChatID: chatID,
			Photo: &models.InputFileUpload{
				//Filename: entry.Name(),
				Data: bytes.NewReader(fileBytes),
			},
		})

		if err != nil {
			c.log.Error("failed to send image", "command", "/cancel", "chat_id", chatID, "error", err)
			continue
		}

		_, err = b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: chatID,
			Text:   fmt.Sprintf("%d", t),
		})

		if err != nil {
			c.log.Error("failed to send image", "command", "/cancel", "chat_id", chatID, "error", err)
			continue
		}
	}
}
