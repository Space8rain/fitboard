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
	b.RegisterHandler(bot.HandlerTypeMessageText, "/start", bot.MatchTypeCommandStartOnly, LoggerMiddleware(startHandler))
	b.RegisterHandler(bot.HandlerTypeCallbackQueryData, "confirm_role", bot.MatchTypeExact, confirmRoleHandler)
	b.RegisterHandler(bot.HandlerTypeCallbackQueryData, "disapprove_role", bot.MatchTypeExact, disapproveRoleHandler)

	b.RegisterHandler(bot.HandlerTypeMessageText, "req", bot.MatchTypePrefix, LoggerMiddleware(startHandler))
	b.RegisterHandler(bot.HandlerTypeMessageText, "del", bot.MatchTypePrefix, LoggerMiddleware(deleteUserHandler))
}

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

	if callback.Message.Message != nil {
		msg := callback.Message.Message

		// Меняем текст и убираем старую inline‑клавиатуру
		_, err := b.EditMessageText(ctx, &bot.EditMessageTextParams{
			ChatID:    msg.Chat.ID,
			MessageID: msg.ID,
			Text:      "Роль тренера подтверждена ✅",
			ReplyMarkup: &models.InlineKeyboardMarkup{
				InlineKeyboard: [][]models.InlineKeyboardButton{
					{
						{
							Text: "Открыть в приложении:",
							WebApp: &models.WebAppInfo{
								URL: "https://www.google.com/",
							},
						},
					},
				},
			},
		})
		if err != nil {
			log.Printf("Ошибка при редактировании сообщения: %v", err)
		}

		// 3. Отправляем новое сообщение с reply‑клавиатурой
		// _, err = b.SendMessage(ctx, &bot.SendMessageParams{
		// 	ChatID: msg.Chat.ID,
		// 	Text:   "Выберите действие:",
		// 	ReplyMarkup: &models.ReplyKeyboardMarkup{
		// 		Keyboard: [][]models.KeyboardButton{
		// 			{
		// 				{Text: "➕ Добавить тренировку"},
		// 				{Text: "👤 Добавить клиента"},
		// 			},
		// 			{
		// 				{Text: "✏️ Редактировать тренировку"},
		// 				{Text: "📋 Остальное"},
		// 			},
		// 			{
		// 				{Text: "📞 Отправить номер", RequestContact: true},
		// 				{Text: "📍 Отправить геопозицию", RequestLocation: true},
		// 				{Text: "📊 Создать опрос", RequestPoll: &models.KeyboardButtonPollType{Type: "regular"}},
		// 			},
		// 		},
		// 		ResizeKeyboard:  true,  // клавиатура подгоняется под экран
		// 		OneTimeKeyboard: false, // не исчезает сразу
		// 	},
		// })
		// if err != nil {
		// 	log.Printf("Ошибка при отправке клавиатуры: %v", err)
		// }
	}
}

func disapproveRoleHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	callback := update.CallbackQuery

	// 1. Ответим на сам callback (убираем "часики")
	// b.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{
	// 	CallbackQueryID: callback.ID,
	// 	Text:            "❌ Хорошо, ждем когда вам назначат тренировки",
	// 	ShowAlert:       false,
	// })

	// 2. Проверяем, что сообщение доступно
	if callback.Message.Message != nil {
		msg := callback.Message.Message

		// Меняем текст и убираем клавиатуру
		_, err := b.EditMessageText(ctx, &bot.EditMessageTextParams{
			ChatID:      msg.Chat.ID,
			MessageID:   msg.ID,
			Text:        "Хорошо, сообщим когда вам назначат тренировки 🏅",
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

	// 3. удаляем reply‑клавиатуру
	_, err = b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		// Text:   "удаление клавиатуры",
		ReplyMarkup: &models.ReplyKeyboardRemove{
			RemoveKeyboard: true,
			Selective:      false, // если true — убирается только у конкретного пользователя
		},
	})
	if err != nil {
		log.Printf("Ошибка при удалении клавиатуры: %v", err)
	}
}

func replyKeyboardMessageHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	if update.Message == nil {
		return
	}

	text := update.Message.Text
	// userID := update.Message.From.ID

	switch text {
	case "➕ Добавить тренировку":
		// Логика добавления тренировки
		_, _ = b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   "Введите данные для новой тренировки (дата, время, описание):",
		})

	case "👤 Добавить клиента":
		// Логика добавления клиента
		_, _ = b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   "Введите имя клиента:",
		})

	case "✏️ Редактировать тренировку":
		// Логика редактирования
		_, _ = b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   "Выберите тренировку для редактирования:",
		})

	case "📋 Остальное":
		// Дополнительные действия
		_, _ = b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   "Доступные функции: перенос, отмена, напоминания об оплате.",
		})

	default:
		// Ответ по умолчанию
		_, _ = b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   "Не понял команду. Используйте кнопки на клавиатуре 👇",
		})
	}
}
