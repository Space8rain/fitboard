package main

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	_ "github.com/lib/pq" // драйвер PostgreSQL
	"golang.org/x/sync/errgroup"

	"fitboard/backend/internal/config"
	"fitboard/backend/internal/db"
	"fitboard/backend/internal/handlers"
	"fitboard/backend/internal/httpserver"
	"fitboard/backend/internal/tgbot"
)

func main() {
	// Обработка Ctrl+C
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	cfg := config.Load()

	// Подключение к БД
	database, err := sql.Open("postgres", cfg.DBConnStr)
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

	// Инициализация HTTP сервера
	h := handlers.New()
	r := httpserver.NewRouter(h)

	srv := &http.Server{
		Addr: ":3000", Handler: r,
	}

	// Группа параллельных задач
	g, gctx := errgroup.WithContext(ctx)

	// Запуск бота
	g.Go(func() error {
		log.Println("Бот запущен...")
		bot.Start(gctx)
		return nil
	})

	// Запуск HTTP сервера
	g.Go(func() error {
		log.Println("HTTP server started on :3000")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			return err
		}
		return nil
	})

	// Грейсфул-шатдаун
	g.Go(func() error {
		<-gctx.Done()
		log.Println("Завершение работы...")

		// Остановка HTTP сервера
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		srv.Shutdown(shutdownCtx)

		return nil
	})

	// Ожидание завершения
	if err := g.Wait(); err != nil {
		log.Println("Ошибка:", err)
	}
}
