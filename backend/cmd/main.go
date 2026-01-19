package main

import (
	"context"
	"database/sql"
	_ "github.com/lib/pq" // драйвер PostgreSQL
	"log"
	"net/http"
	"os"
	"os/signal"
	"github.com/go-chi/chi/v5"

	"fitboard/backend/internal/db"
	"fitboard/backend/internal/tgbot"
	"fitboard/backend/internal/handlers"
)

func main() {
	// Обработка Ctrl+C
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	// Подключение к БД
	connStr := "postgres://postgres:965478@localhost:5432/fitboard_db?sslmode=disable"

	database, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Ошибка подключения к БД: %v", err)
	}
	defer database.Close()

	db.InitRepo(database)

	// Инициализация бота
	bot, err := tgbot.NewBot()
	if err != nil {
		log.Fatalf("❌ Ошибка инициализации бота: %v", err)
	}

	tgbot.RegisterHandlers(bot)

	// HTTP‑маршруты для фронта
	r := chi.NewRouter()
	r.Get("/api/ping", handlers.Ping)
	r.Get("/api/users", handlers.Users)
	r.Get("/api/trainers", handlers.Trainers)

	// Запуск HTTP‑сервера
	go func() {
		log.Println("HTTP сервер запущен на :3000")
		if err := http.ListenAndServe(":3000", r); err != nil {
			log.Fatalf("Ошибка запуска HTTP сервера: %v", err)
		}
	}()

	// Запуск бота
	log.Printf("Бот запущен...")
	bot.Start(ctx)
}
