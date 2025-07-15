package tgbot

import (
    "context"
    "fmt"
    "github.com/go-telegram/bot"
    "github.com/go-telegram/bot/models"
)

func RegisterHandlers(b *bot.Bot) {
	b.RegisterHandler(bot.HandlerTypeMessageText, "/start", bot.MatchTypeCommandStartOnly, LoggerMiddleware(startHandler))

	b.RegisterHandler(bot.HandlerTypeMessageText, "", bot.MatchTypePrefix, echoHandler)

}

// var idDog int64 = 469895624
// var idMy int64 = 413870391

// Обработчик сообщений
func startHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	user := update.Message.From
	text := update.Message.Text
	chat := update.Message.Chat

	fmt.Printf("👤 Пользователь: %s (@%s)\n", user.FirstName, user.Username)
	fmt.Printf("🆔 Telegram ID: %d\n", user.ID)
	fmt.Printf("📣 Чат ID: %d\n", chat.ID)
	fmt.Printf("✉️ Сообщение: %s\n", text)

	answer := fmt.Sprintf("👋 Привет, %s! Выбери ниже", user.FirstName)
	fmt.Println("⚡ Обработчик /start запущен")

	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   answer,
	})
}

func echoHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   update.Message.Text,
	})
}