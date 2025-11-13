package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"database/sql"
	_ "github.com/lib/pq" // драйвер PostgreSQL

	"fitboard/backend/internal/tgbot"
	"fitboard/backend/internal/db"
)

func main() {
	// Обработка Ctrl+C
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	connStr := "postgres://postgres:965478@localhost:5432/fitboard_db?sslmode=disable"

	database, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Ошибка подключения к БД: %v", err)
	}
	defer database.Close()

	db.InitRepo(database)

	bot, err := tgbot.NewBot()
	if err != nil {
		log.Fatalf("❌ Ошибка инициализации бота: %v", err)
	}

	tgbot.RegisterHandlers(bot)
	log.Printf("Бот запущен...")
	bot.Start(ctx)
}
