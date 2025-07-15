package tgbot

import (
	"context"
	"log"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

func LoggerMiddleware(next bot.HandlerFunc) bot.HandlerFunc {
	return func(ctx context.Context, b *bot.Bot, update *models.Update) {
		if update.Message != nil {
			user := update.Message.From
			text := update.Message.Text
			chat := update.Message.Chat

			log.Printf("👤 Пользователь: %s (@%s)\n", user.FirstName, user.Username)
			log.Printf("🆔 Telegram ID: %d\n", user.ID)
			log.Printf("📣 Чат ID: %d\n", chat.ID)
			log.Printf("✉️ Сообщение: %s\n", text)

			next(ctx, b, update)
		}
	}
}