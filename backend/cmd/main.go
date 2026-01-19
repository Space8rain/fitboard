package main

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"os"
	"os/signal"

	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq" // драйвер PostgreSQL

	"fitboard/backend/internal/db"
	"fitboard/backend/internal/handlers"
	"fitboard/backend/internal/tgbot"
)

func main() {
	// Обработка Ctrl+C
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	// Подключение к БД
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatal("Ошибка загрузки .env файла")
	}

	DB_CONN_STR := os.Getenv("DB_CONN_STR")
	
	database, err := sql.Open("postgres", DB_CONN_STR)
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
