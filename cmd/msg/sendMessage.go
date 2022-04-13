package msg

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
)

func SendMessage(bot *tgbotapi.BotAPI, chatID int64, message string) {
	msg := tgbotapi.NewMessage(
		chatID,
		message,
	)
	msg.ParseMode = "Markdown"
	msg.DisableWebPagePreview = true
	if _, e := bot.Send(msg); e != nil {
		log.Printf("Error sending message to telegram.\nMessage: %v\nError: %v", msg, e)
	}
}

func SendInlineSlugMissMessage(bot *tgbotapi.BotAPI, chatID int64, message string, matches []string) {
	var inlineKeyboard = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(matches[0], fmt.Sprintf(matches[0])),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(matches[1], fmt.Sprintf(matches[1])),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(matches[2], fmt.Sprintf(matches[2])),
		),
	)

	msg := tgbotapi.NewMessage(chatID, message)
	msg.ReplyMarkup = inlineKeyboard
	if _, e := bot.Send(msg); e != nil {
		log.Printf("Error sending message to telegram.\nMessage: %v\nError: %v", msg, e)
	}
}
