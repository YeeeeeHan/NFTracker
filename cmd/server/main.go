package main

import (
	"NFTracker/cmd/handlers"
	"NFTracker/datastorage"
	"NFTracker/pkg/api"
	"NFTracker/pkg/db"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
)

func main() {
	// If env is not prod, use read env from local.env
	if os.Getenv("ENV") != "PROD" {
		log.Printf("[env] Reading from local.env")
		err := godotenv.Load("local.env")
		if err != nil {
			log.Fatalf("Some err occured. Err: %s", err)
		}
	}

	// Init API and DB
	db, err := db.NewDB()
	if err != nil {
		panic(err)
	}

	log.Printf("[main] We're up and running!")
	port := os.Getenv("PORT")

	go func() {
		router := api.NewAPI(db)
		err = http.ListenAndServe(fmt.Sprintf(":%s", port), router)
		if err != nil {
			log.Printf("err from  router: %v\n", err)
		}
	}()

	// Initialize cache
	_ = datastorage.InitCache()
	log.Printf("Cache init")

	var bot *tgbotapi.BotAPI
	var boterr error
	bot, boterr = tgbotapi.NewBotAPI(os.Getenv("BOTTOKEN"))
	if boterr != nil {
		log.Panic(boterr)
	}
	bot.Debug = false

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.CallbackQuery != nil {
			chatID := update.CallbackQuery.Message.Chat.ID
			username := update.CallbackQuery.Message.From.UserName
			callbackData := update.CallbackQuery.Data
			log.Printf("\n\nReceived message in (chatID: %d) from %s: (callbackData: %v) \n\n", chatID, username, callbackData)

			// Respond to the callback query, telling Telegram to show the user
			// a message with the data received.
			callback := tgbotapi.NewCallback(update.CallbackQuery.ID, callbackData)
			if _, err := bot.AnswerCallbackQuery(callback); err != nil {
				continue
			}

			// And finally, send a message containing the data received.
			handlers.PriceCheck(db, bot, chatID, username, callbackData)
		}

		if update.Message == nil { // ignore any non-Message Updates
			continue
		}

		chatID := update.Message.Chat.ID
		username := update.Message.From.UserName
		t := update.Message.Text
		log.Printf("\n\nReceived message in (chatID: %d) from %s: %s (command: %v) \n\n", chatID, username, t, update.Message.IsCommand())

		// Replying to a message
		//msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
		//msg.ReplyToMessageID = update.Message.MessageID
		// bot.Send(msg)

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
				handlers.PriceCheck(db, bot, chatID, username, update.Message.CommandArguments())
			case "alert":
				handlers.Alert(db, bot, chatID, update.Message.CommandArguments())
			case "start", "help":
				handlers.Introduction(bot, chatID)
			default:
				bot.Send(tgbotapi.NewMessage(chatID, "🤔 Command not recognised."))
			}
		}
	}
}
