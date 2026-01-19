package tgbot

import (
	"fitboard/backend/config"
	"log"

	"github.com/go-telegram/bot"
)

func NewBot() (*bot.Bot, error) {
	cfg := config.Load()

	if cfg.BotToken == "" {
		log.Fatal("TOKEN_BOT не найден")
	}

	return bot.New(cfg.BotToken)
}
