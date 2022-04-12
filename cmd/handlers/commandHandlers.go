package handlers

import (
	"NFTracker/cmd"
	"NFTracker/datastorage"
	"NFTracker/pkg/db"
	"NFTracker/pkg/opensea"
	"fmt"
	"github.com/go-pg/pg/v10"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/patrickmn/go-cache"
	"log"
	"strconv"
)

func Introduction(bot *tgbotapi.BotAPI, chatID int64) {
	msg := tgbotapi.NewMessage(chatID, cmd.WelcomeMessage)
	msg.ParseMode = "Markdown"
	msg.DisableWebPagePreview = true
	if _, e := bot.Send(msg); e != nil {
		log.Printf("Error sending message to telegram.\nMessage: %v\nError: %v", msg, e)
	}
	return
}

func PriceCheck(pgdb *pg.DB, bot *tgbotapi.BotAPI, chatID int64, userName, slug string) {
	if slug == "" {
		msg := tgbotapi.NewMessage(chatID, "No slug detected.")
		if _, e := bot.Send(msg); e != nil {
			log.Printf("Error sending message to telegram.\nMessage: %v\nError: %v", msg, e)
		}
		return
	}

	defer func() {
		_, err := db.GetCustomer(pgdb, userName)
		if err == pg.ErrNoRows {
			log.Printf("[db.GetCustomer] New user... adding user to DB...")
			_, _ = db.CreateCustomer(pgdb, userName)
		}
	}()

	// Check Cache, if found return from cache
	if x, found := datastorage.GlobalCache.Get(slug); found {
		osResponse := x.(*opensea.OSResponse)
		msg := tgbotapi.NewMessage(
			chatID,
			cmd.PriceCheckMessage(slug, opensea.CreateUrlFromSlug(slug), osResponse),
		)
		msg.ParseMode = "Markdown"
		msg.DisableWebPagePreview = true
		if _, e := bot.Send(msg); e != nil {
			log.Printf("Error sending message to telegram.\nMessage: %v\nError: %v", msg, e)
		}
		return
	}

	// Else query web and update cache
	osResponse, err := opensea.Scrape(slug)
	if err != nil {
		log.Printf("[opensea.Scrape] %v", err)
		msg := tgbotapi.NewMessage(chatID, "TEST")
		if _, e := bot.Send(msg); e != nil {
			log.Printf("Error sending message to telegram.\nMessage: %v\nError: %v", msg, e)
		}
		return
	}

	// Update cache - Set the value of the key "slug" to fp,
	// with the default expiration time
	datastorage.GlobalCache.Set(slug, osResponse, cache.DefaultExpiration)

	msg := tgbotapi.NewMessage(
		chatID,
		cmd.PriceCheckMessage(slug, opensea.CreateUrlFromSlug(slug), osResponse),
	)
	msg.ParseMode = "Markdown"
	msg.DisableWebPagePreview = true
	if _, e := bot.Send(msg); e != nil {
		log.Printf("Error sending message to telegram.\nMessage: %v\nError: %v", msg, e)
	}
	return
}

func Alert(db *pg.DB, bot *tgbotapi.BotAPI, chatID int64, arguments string) {
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
