package file

import (
	"context"
	"errors"
	"gustav-tgbot/internal/model"
	userErr "gustav-tgbot/internal/service/user"
	"log/slog"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

type FileService interface {
	StoreFile(ctx context.Context, filePath, fileURL string, chatID int64) error
}

type UserService interface {
	GetUserByChatID(ctx context.Context, chatID int64) (*model.User, error)
}

type File struct {
	FileService FileService
	UserService UserService
	log         *slog.Logger
}

func NewFileHandler(fileService FileService, userService UserService, log *slog.Logger) *File {
	return &File{
		FileService: fileService,
		UserService: userService,
		log:         log,
	}
}

func (f *File) FileAcceptMessage(ctx context.Context, b *bot.Bot, update *models.Update) {
	message := update.Message
	chatID := message.Chat.ID

	user, err := f.UserService.GetUserByChatID(ctx, chatID)
	if err != nil {
		if errors.Is(err, userErr.ErrUserNotFound) {
			_, err := b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: chatID,
				Text:   "Регистрация не прошла. Напишите /start для начала.",
			})
			if err != nil {
				f.log.Error("failed to send message", "file", "file accept", "chat_id", chatID, "error", err)
				return
			}
		}
		f.log.Error("failed to get user by chat id", "file", "file accept", "chat_id", chatID, "error", err)
		return
	}

	if user.StateUser != model.StateSendFile {
		f.log.Info("user is in inactivity state", "file", "file accept", "chat_id", chatID)
		_, err := b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: chatID,
			Text:   "Нет смысла отправлять фото. Напишите /watermark для начала работы."})
		if err != nil {
			f.log.Error("failed to send message", "file", "file accept", "chat_id", chatID, "error", err)
			return
		}
		return
	}

	//if message.MediaGroupID != "" {
	//	_, err := b.SendMessage(ctx, &bot.SendMessageParams{
	//		ChatID: chatID,
	//		Text:   "Дружище, одну фотографию.",
	//	})
	//	if err != nil {
	//		f.log.Error("failed to send message", "file", "file accept", "chat_id", chatID, "error", err)
	//		return
	//	}
	//	f.log.Error("not more one photo!!!", "file", "file accept", "chat_id", chatID)
	//}

	if len(message.Photo) > 1 {
		photo := message.Photo[len(message.Photo)-1]

		file, err := b.GetFile(ctx, &bot.GetFileParams{
			FileID: photo.FileID,
		})
		if err != nil {
			f.log.Error("failed to get file", "file", "file accept", "chat_id", chatID, "error", err)
			return
		}

		fileURL := b.FileDownloadLink(file)

		err = f.FileService.StoreFile(ctx, file.FilePath, fileURL, chatID)
		if err != nil {
			f.log.Error("failed store file", "file", "file accept", "chat_id", chatID, "error", err)
			return
		}
	}

	if message.Document != nil {
		document := message.Document

		file, err := b.GetFile(ctx, &bot.GetFileParams{
			FileID: document.FileID,
		})
		if err != nil {
			f.log.Error("failed to get file", "file", "document accept", "chat_id", chatID, "error", err)
			return
		}

		fileURL := b.FileDownloadLink(file)

		err = f.FileService.StoreFile(ctx, file.FilePath, fileURL, chatID)
		if err != nil {
			f.log.Error("failed store file", "file", "document accept", "chat_id", chatID, "error", err)
			return
		}
	}

}
