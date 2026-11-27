package handler

import (
	"context"
	"fmt"
	"gustav-tgbot/internal/service"
	"log"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

type Router struct {
	log *slog.Logger
	B   *bot.Bot
}

func NewRouter(service *service.Service, b *bot.Bot, log *slog.Logger) *Router {
	handler := NewHandler(service, log)

	b.RegisterHandler(bot.HandlerTypeMessageText, "/start", bot.MatchTypeExact, handler.CommandHandler.StartMessage)
	b.RegisterHandler(bot.HandlerTypeMessageText, "/watermark", bot.MatchTypeExact, handler.CommandHandler.WatermarkMessage)
	b.RegisterHandler(bot.HandlerTypeMessageText, "/cancel", bot.MatchTypeExact, handler.CommandHandler.CancelMessage)
	b.RegisterHandlerMatchFunc(matchFileAccept, handler.FileHandler.FileAcceptMessage)
	return &Router{log: log, B: b}
}

func matchFileAccept(update *models.Update) bool {
	if len(update.Message.Photo) > 0 {
		return true
	}
	if update.Message.Document != nil {
		return true
	}
	return false
}

func MiddlewareMessage(ctx context.Context, b *bot.Bot, update *models.Update) {
	if update.Message == nil {
		return
	}

	message := update.Message

	chatID := message.Chat.ID
	username := message.Chat.Username
	text := message.Text
	caption := message.Caption

	log.Println("========== NEW MESSAGE ==========")
	log.Printf("Chat ID: %d\n", chatID)
	log.Printf("Username: @%s\n", username)
	log.Printf("Message ID: %d\n", message.ID)
	log.Printf("Date: %d\n", message.Date)

	if text != "" {
		log.Printf("Text: %s\n", text)
	}

	if caption != "" {
		log.Printf("Caption: %s\n", caption)
	}

	log.Println("=================================")
}

func saveTelegramFile(ctx context.Context, b *bot.Bot, fileID string, directory string) {
	// Получаем информацию о файле.
	file, err := b.GetFile(ctx, &bot.GetFileParams{
		FileID: fileID,
	})
	if err != nil {
		log.Printf("Get file error: %v\n", err)
		return
	}

	// Получаем содержимое файла.
	fileURL := b.FileDownloadLink(file)

	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	resp, err := client.Get(fileURL)
	if err != nil {
		log.Printf("Download file error: %v\n", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		log.Printf("Download failed: HTTP %d\n", resp.StatusCode)
		return
	}

	// Формируем имя файла.
	ext := filepath.Ext(file.FilePath)
	if ext == "" {
		ext = ".bin"
	}

	filename := fmt.Sprintf("%d_%s%s",
		time.Now().UnixNano(),
		strings.TrimSuffix(filepath.Base(file.FilePath), ext),
		ext,
	)

	path := filepath.Join(directory, filename)

	fmt.Println(path)

	// Создаём файл.
	dst, err := os.Create(path)
	if err != nil {
		log.Printf("Create file error: %v\n", err)
		return
	}
	defer dst.Close()

	if _, err := dst.ReadFrom(resp.Body); err != nil {
		log.Printf("Save file error: %v\n", err)
		return
	}

	log.Printf("File successfully saved: %s\n", path)
}
