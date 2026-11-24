package app

import (
	"context"
	"gustav-tgbot/config"
	"gustav-tgbot/internal/database"
	"gustav-tgbot/internal/handler"
	"gustav-tgbot/internal/repository"
	"gustav-tgbot/internal/service"
	"gustav-tgbot/pkg/logger"
	"log/slog"
	"os/signal"
	"syscall"

	"github.com/go-telegram/bot"
)

func Run(cfg *config.Config) {
	initCheckConfig(cfg)

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	log := logger.NewLogger(cfg.App.State)

	db, err := database.NewSQLite(cfg.Database.Path)
	if err != nil {
		panic(err)
	}
	defer db.Close()
	log.Info("Configuration loaded: ", cfg)

	repository := repository.NewRepository(db, log)
	service := service.NewService(repository, log)

	initTelegramBot(ctx, cfg, service, log)

}

func initTelegramBot(ctx context.Context, cfg *config.Config, srvc *service.Service, log *slog.Logger) {
	b, err := bot.New(cfg.Telegram.Token)
	if err != nil {
		log.Error("failed to create new bot: %w", err)
	}

	router := handler.NewRouter(srvc, b, log)

	router.B.Start(ctx)
}

func initCheckConfig(cfg *config.Config) {
	if cfg.App.Name == "" {
		panic("app name is required")
	}

	if cfg.Telegram.Token == "" {
		panic("telegram token is required")
	}
}
