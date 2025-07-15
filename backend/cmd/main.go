package main

import (
	"context"
	"log"
	"os"
	"os/signal"

	"fitboard/backend/internal/tgbot"
)

func main() {
	// Обработка Ctrl+C
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	// Инициализация бота
	bot, err := tgbot.NewBot()
	if err != nil {
		log.Fatalf("❌ Ошибка инициализации бота: %v", err)
	}

	tgbot.RegisterHandlers(bot)

	log.Printf("Бот запущен...")
	bot.Start(ctx)
}
