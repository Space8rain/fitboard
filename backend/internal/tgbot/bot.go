package tgbot

import (
  "github.com/go-telegram/bot"
	"github.com/joho/godotenv"
	"log"
	"os"

)


func NewBot() (*bot.Bot, error) {
  	// Загружаем .env
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatal("Ошибка загрузки .env файла")
	}

	token := os.Getenv("TOKEN_BOT")
	if token == "" {
		log.Fatal("TOKEN_BOT не найден")
	}

  return bot.New(token)
}