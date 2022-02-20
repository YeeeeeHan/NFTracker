package main

import (
	"NFTracker/config"
	"NFTracker/datastorage"
	"NFTracker/handlers"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
)

func main() {
	// Initialize cache
	_ = datastorage.InitCache()
	log.Printf("Cache init")

	bot, err := tgbotapi.NewBotAPI(config.BotToken)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil { // ignore any non-Message Updates
			continue
		}

		chatID := update.Message.Chat.ID
		username := update.Message.From.UserName
		t := update.Message.Text
		log.Printf("\n\nReceived message in (chatID: %d) from %s: %s (command: %v) \n\n", chatID, username, t, update.Message.IsCommand())

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
		//msg.ReplyToMessageID = update.Message.MessageID
		bot.Send(msg)

		switch {
		case update.Message.IsCommand():
			// Handle commands
			//
			// TODO: Check that the bot is set up for `alert` command
			// and add it if not.
			// Currently hardcoded in setup process, but outline function
			// has been added. Need to change it to take existing commands,
			// and add the new one (rather than overwrite)

			switch update.Message.Command() {
			case "check":
				handlers.PriceCheck(bot, chatID, update.Message.CommandArguments())
			case "alert":
				handlers.Alert(bot, chatID, update.Message.CommandArguments())
			case "start", "help":
				handlers.Introduction(bot, chatID)
			default:
				bot.Send(tgbotapi.NewMessage(chatID, "🤔 Command not recognised."))
			}
		}
	}
}
