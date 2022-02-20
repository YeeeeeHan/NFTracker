package handlers

import (
	"NFTracker/datastorage"
	"NFTracker/scraping"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/patrickmn/go-cache"
	"log"
	"strconv"
)

func Introduction(bot *tgbotapi.BotAPI, chatID int64) {
	msg := tgbotapi.NewMessage(chatID, "Welcome to the 🚗 *Car Park Telegram Bot* 🚗\n_Powered by "+
		"Apache Kafka® and [ksqlDB](https://ksqldb.io)_ 😃\n\n👉 Use `/alert \\<x\\>` to receive an alert "+
		"when a car park has more than \\<x\\> places available\n👉 Send me the name of a car park to find "+
		"out how many spaces are currently available in it\n👉 Send me your location to find out the nearest "+
		"car park to you with more than 10 spaces\\.")
	msg.ParseMode = "MarkdownV2"
	if _, e := bot.Send(msg); e != nil {
		log.Printf("Error sending message to telegram.\nMessage: %v\nError: %v", msg, e)
	}
	return
}

func PriceCheck(bot *tgbotapi.BotAPI, chatID int64, argument string) {
	if argument == "" {
		msg := tgbotapi.NewMessage(chatID, "Never put anything???")
		if _, e := bot.Send(msg); e != nil {
			log.Printf("Error sending message to telegram.\nMessage: %v\nError: %v", msg, e)
		}
		return
	}

	// Check Cache
	price, found := datastorage.GlobalCache.Get(argument)

	// If found return from cache
	if found {
		msg := tgbotapi.NewMessage(chatID, price.(string))
		if _, e := bot.Send(msg); e != nil {
			log.Printf("Error sending message to telegram.\nMessage: %v\nError: %v", msg, e)
		}
		return
	}

	// Else query web and update cache
	fp, err := scraping.Scrape("https://opensea.io/collection/doodles-official")
	if err != nil {
		log.Printf("@@@@ %v", fp)
	}

	// Set the value of the key "argument" to fp, with the default expiration time
	datastorage.GlobalCache.Set(argument, fp, cache.DefaultExpiration)

	log.Printf("@@@@ %v", fp)
	msg := tgbotapi.NewMessage(chatID, fp)
	if _, e := bot.Send(msg); e != nil {
		log.Printf("Error sending message to telegram.\nMessage: %v\nError: %v", msg, e)
	}
	return
}

func Alert(bot *tgbotapi.BotAPI, chatID int64, arguments string) {
	if th, e := strconv.Atoi(arguments); e == nil {
		// Use a Go Routine to invoke the population
		// of the alert channel and handling the returned
		// alerts
		go func() {
			ac := make(chan string)
			//go alertSpaces(ac, th)
			msg := tgbotapi.NewMessage(chatID, fmt.Sprintf("👍 Successfully created alert to be sent whenever more than %v spaces are available", th))
			if _, e := bot.Send(msg); e != nil {
				log.Printf("Error sending message to telegram.\nMessage: %v\nError: %v", msg, e)
			}

			for a := range ac {
				msg := tgbotapi.NewMessage(chatID, a)
				if _, e := bot.Send(msg); e != nil {
					log.Printf("Error sending message to telegram.\nMessage: %v\nError: %v", msg, e)
				}
			}
		}()
	} else {
		msg := tgbotapi.NewMessage(chatID, "Non-integer value specified for `/alert`")
		if _, e := bot.Send(msg); e != nil {
			log.Printf("Error sending message to telegram.\nMessage: %v\nError: %v", msg, e)
		}

	}
}
