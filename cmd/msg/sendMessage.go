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
	// Set message configs
	msg.ParseMode = "Markdown"
	msg.DisableWebPagePreview = true
	if _, e := bot.Send(msg); e != nil {
		log.Printf("Error sending message to telegram.\nMessage: %v\nError: %v", msg, e)
	}
}

func SendInlineSlugMissMessage(bot *tgbotapi.BotAPI, chatID int64, message, slugQuery string, matches ...string) {
	var buttons [][]tgbotapi.InlineKeyboardButton

	for _, match := range matches {
		buttons = append(buttons, tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(fmt.Sprintf("Did you mean:  /check %s ?", match), match),
		))
	}

	var inlineKeyboard = tgbotapi.NewInlineKeyboardMarkup(buttons...)

	// Set message configs
	msg := tgbotapi.NewMessage(chatID, message)
	msg.ReplyMarkup = inlineKeyboard
	msg.ParseMode = "Markdown"
	msg.DisableWebPagePreview = true
	if _, e := bot.Send(msg); e != nil {
		log.Printf("Error sending message to telegram.\nMessage: %v\nError: %v", msg, e)
	}
}

func SendEditMessage(bot *tgbotapi.BotAPI, chatID int64, messageID int, message string) {
	msg := tgbotapi.NewEditMessageText(
		chatID,
		messageID,
		message,
	)
	msg.ParseMode = "Markdown"
	msg.DisableWebPagePreview = true
	if _, e := bot.Send(msg); e != nil {
		log.Printf("Error sending message to telegram.\nMessage: %v\nError: %v", msg, e)
	}
}
