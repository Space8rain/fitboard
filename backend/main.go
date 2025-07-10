package main

import (
	"context"
	"log"
	"os"
	"os/signal"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/joho/godotenv"
)

func main() {
	// Загружаем .env
	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatal("Ошибка загрузки .env файла")
	}

	token := os.Getenv("TOKEN_BOT")
	if token == "" {
		log.Fatal("TOKEN_BOT не найден")
	}

	// Обработка Ctrl+C
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	// Инициализация бота
	b, err := bot.New(token, bot.WithDefaultHandler(handler))
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Бот запущен...")
	b.Start(ctx)
}

// Обработчик сообщений
func handler(ctx context.Context, b *bot.Bot, update *models.Update) {
	if update.Message != nil {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   "Привет! 🤖",
		})
	}
}
