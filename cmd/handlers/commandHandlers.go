package handlers

import (
	"NFTracker/cmd/msg"
	"NFTracker/datastorage"
	"NFTracker/pkg/custerror"
	"NFTracker/pkg/db"
	"NFTracker/pkg/opensea"
	"fmt"
	"github.com/go-pg/pg/v10"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/patrickmn/go-cache"
	"log"
	"strconv"
	"strings"
)

func Introduction(bot *tgbotapi.BotAPI, chatID int64) {
	msg := tgbotapi.NewMessage(chatID, msg.WelcomeMessage)
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

	matches := opensea.FindClosestmatch(slug, 3)
	// Slug not found in top collections
	if slug != matches[0] {
		// Send clarification message
		message := "Collection not found, did you mean any of these:\n\n"
		msg.SendInlineSlugMissMessage(bot, chatID, message, matches)
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
		// assert type
		osResponse := x.(*opensea.OSResponse)

		// Send price check message
		message := msg.PriceCheckMessage(slug, opensea.CreateUrlFromSlug(slug), osResponse)
		msg.SendMessage(bot, chatID, message)
		return
	}

	// Else query web
	osResponse, err := opensea.QueryAPI(slug)

	// Slug not found
	if err == custerror.InvalidSlugErr {
		matches := opensea.FindClosestmatch(slug, 3)
		// Send clarification message
		message := "Collection not found, did you mean any of these:\n\n" + strings.Join(matches, ", ")
		msg.SendMessage(bot, chatID, message)
		return
	}

	if err != nil {
		message := "Sorry there is an internal error, please try again!"
		msg.SendMessage(bot, chatID, message)
		log.Printf("[opensea.QueryAPI] %v", err)
		return
	}

	// Update cache - Set the value of the key "slug" to fp,
	// with the default expiration time
	datastorage.GlobalCache.Set(slug, osResponse, cache.DefaultExpiration)

	// Send price check message
	message := msg.PriceCheckMessage(slug, opensea.CreateUrlFromSlug(slug), osResponse)
	msg.SendMessage(bot, chatID, message)

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
