package main

import (
	"gustav-tgbot/config"
	"gustav-tgbot/internal/app"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		panic(err)
	}
	app.Run(cfg)
}
