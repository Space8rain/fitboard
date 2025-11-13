package tgbot

import (
	"context"
	"fitboard/backend/internal/db"
	"fmt"
	"log"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

func RegisterHandlers(b *bot.Bot) {
	// b.RegisterHandler(bot.HandlerTypeMessageText, "/start", bot.MatchTypeCommandStartOnly, LoggerMiddleware(startHandler))
	b.RegisterHandler(bot.HandlerTypeCallbackQueryData, "confirm_role", bot.MatchTypeExact, confirmRoleHandler)
	b.RegisterHandler(bot.HandlerTypeCallbackQueryData, "disapprove_role", bot.MatchTypeExact, disapproveRoleHandler)

	b.RegisterHandler(bot.HandlerTypeMessageText, "req", bot.MatchTypePrefix, LoggerMiddleware(startHandler))
	b.RegisterHandler(bot.HandlerTypeMessageText, "del", bot.MatchTypePrefix, LoggerMiddleware(deleteUserHandler))
}

// var idDog int64 = 469895624
// var idMy int64 = 413870391

// Обработчик сообщений
func startHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	user := update.Message.From

	exists, err := db.Repo.Exists(user.ID)

	if err != nil {
		log.Fatalf("❌ Ошибка проверки пользователя: %v", err)
		return
	}

	if exists {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   fmt.Sprintf("👋 С возвращением, %s!", user.FirstName),
		})
		return
	}

	newUser := db.User{
		ID:        user.ID,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Username:  user.Username,
		Role:      "client",
	}

	err = db.Repo.CreateUser(newUser)
	if err != nil {
		log.Fatalf("❌ Ошибка регистрации пользователя: %v", err)
		return
	}

	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   fmt.Sprintf("👋 Привет, %s! Ты тренер?", user.FirstName),
		ReplyMarkup: &models.InlineKeyboardMarkup{
			InlineKeyboard: [][]models.InlineKeyboardButton{
				{
					{Text: "✅ Да", CallbackData: "confirm_role"},
				},
				{
					{Text: "❌ Нет", CallbackData: "disapprove_role"},
				},
			},
		},
	})
}

func confirmRoleHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	callback := update.CallbackQuery

	// 1. Обновляем роль в базе
	if err := db.Repo.UpdateUserRole(callback.From.ID, "trainer"); err != nil {
		log.Printf("❌ Ошибка обновления роли пользователя: %v", err)
		return
	}

	// 2. Отвечаем на сам callback (убираем "часики")
	b.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{
		CallbackQueryID: callback.ID,
		Text:            "✅ Теперь вам доступен функционал тренера",
		ShowAlert:       false,
	})

	// 3. Проверяем, что сообщение доступно
	if callback.Message.Message != nil {
		msg := callback.Message.Message

		// Меняем текст и убираем клавиатуру
		_, err := b.EditMessageText(ctx, &bot.EditMessageTextParams{
			ChatID:      msg.Chat.ID,
			MessageID:   msg.ID,
			Text:        "👋 Привет! Роль тренера подтверждена ✅",
			ReplyMarkup: nil,
		})
		if err != nil {
			log.Printf("Ошибка при редактировании сообщения: %v", err)
		}
	}
}

func disapproveRoleHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	callback := update.CallbackQuery

	// 1. Ответим на сам callback (убираем "часики")
	b.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{
		CallbackQueryID: callback.ID,
		Text:            "❌ Хорошо, ждем когда вам назначат тренировки",
		ShowAlert:       false,
	})

	// 2. Проверяем, что сообщение доступно
	if callback.Message.Message != nil {
		msg := callback.Message.Message

		// Меняем текст и убираем клавиатуру
		_, err := b.EditMessageText(ctx, &bot.EditMessageTextParams{
			ChatID:      msg.Chat.ID,
			MessageID:   msg.ID,
			Text:        "Роль тренера отклонена ❌",
			ReplyMarkup: nil,
		})
		if err != nil {
			log.Printf("Ошибка при редактировании сообщения: %v", err)
		}
	}
}

func deleteUserHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	user := update.Message.From

	deleted, err := db.Repo.DeleteUser(user.ID)
	if err != nil {
		log.Printf("❌ Ошибка удаления пользователя: %v", err)
		return
	}

	var text string
	if deleted {
		text = fmt.Sprintf("👋 Пользователь %s удален.", user.FirstName)
	} else {
		text = fmt.Sprintf("ℹ️ Пользователь %s не найден в базе.", user.FirstName)
	}

	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   text,
	})
}
