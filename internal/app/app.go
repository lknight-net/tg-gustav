package app

import (
	"context"
	"fmt"
	"gustav-tgbot/config"
	"gustav-tgbot/internal/handler"
	"gustav-tgbot/internal/service"
	"gustav-tgbot/pkg/logger"
	"log/slog"
	"os"
	"os/signal"

	"github.com/go-telegram/bot"
)

func Run() {
	cfg, err := config.LoadConfig()
	if err != nil {
		panic(fmt.Errorf("failed to load config: %w", err))
	}
	fmt.Println(cfg)

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	log := logger.NewLogger(cfg.App.State)

	srvc := service.NewService(log)

	initTelegramBot(ctx, cfg, srvc, log)
}

func initTelegramBot(ctx context.Context, cfg *config.Config, srvc *service.Service, log *slog.Logger) {
	b, err := bot.New(cfg.Telegram.Token)
	if err != nil {
		log.Error("failed to create new bot: %w", err)
	}

	router := handler.NewRouter(b, srvc, log)

	router.B.Start(ctx)
}
