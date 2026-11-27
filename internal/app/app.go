package app

import (
	"context"
	"database/sql"
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

	//initMigrations(db, log)

	log.Info("Configuration loaded: ", cfg)

	repository := repository.NewRepository(db, log)
	service := service.NewService(repository, log)

	initTelegramBot(ctx, cfg, service, log)

}

func initMigrations(db *sql.DB, log *slog.Logger) {
	migrate, err := database.NewMigrator(db)
	if err != nil {
		log.Error("failed to create migrations: ", err)
	}

	defer func(migrate *database.Migrator) {
		err := migrate.Close()
		if err != nil {
			log.Error("failed to close migrations: ", err)
		}
	}(migrate)

	err = migrate.Up()
	if err != nil {
		log.Error("failed to up migrations: ", err)
	}
}

func initTelegramBot(ctx context.Context, cfg *config.Config, srvc *service.Service, log *slog.Logger) {
	opts := []bot.Option{
		bot.WithDefaultHandler(handler.MiddlewareMessage),
	}

	b, err := bot.New(cfg.Telegram.Token, opts...)
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
