package handlers

import (
	"NFTracker/cmd/msg"
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

const POPULAR_COLLECTION_NUM_OWNERS = 400

func Introduction(bot *tgbotapi.BotAPI, chatID int64) {
	msg := tgbotapi.NewMessage(chatID, msg.WelcomeMessage)
	msg.ParseMode = "Markdown"
	msg.DisableWebPagePreview = true
	if _, e := bot.Send(msg); e != nil {
		log.Printf("Error sending message to telegram.\nMessage: %v\nError: %v", msg, e)
	}
	return
}

func PriceCheckWithSlugMatch(pgdb *pg.DB, bot *tgbotapi.BotAPI, chatID int64, userName, slugQuery string) {
	if slugQuery == "" {
		msg := tgbotapi.NewMessage(chatID, "No slugQuery detected.")
		if _, e := bot.Send(msg); e != nil {
			log.Printf("Error sending message to telegram.\nMessage: %v\nError: %v", msg, e)
		}
		return
	}

	var namelist []string
	sluglist, err := db.GetAllSlugs(pgdb)
	if err != nil {
		log.Printf("@@@@@@@@@@@ err %v", err)
		return
	}
	for _, slug := range sluglist {
		namelist = append(namelist, slug.SlugName)
	}
	log.Printf("@@@@@@@@@@@@ namelist: %v", namelist)

	matches := opensea.FindClosestmatch(slugQuery, 1, namelist)
	// Slugs not found in top collections
	if matches[0] != "" && slugQuery != matches[0] {
		log.Printf(slugQuery, matches)
		// Send clarification message
		message := "Slug not found! Did you mean:\n\n"
		msg.SendInlineSlugMissMessage(bot, chatID, message, slugQuery, matches...)
		return
	}

	PriceCheck(pgdb, bot, chatID, userName, slugQuery)
}

func PriceCheck(pgdb *pg.DB, bot *tgbotapi.BotAPI, chatID int64, userName, slugQuery string) {
	defer func() {
		_, err := db.GetCustomer(pgdb, userName)
		if err == pg.ErrNoRows {
			log.Printf("[db.GetCustomer] New user... adding user to DB...")
			_, _ = db.CreateCustomer(pgdb, userName)
		}
	}()

	// Check Cache, if found return from cache
	if x, found := datastorage.GlobalCache.Get(slugQuery); found {
		// assert type
		osResponse := x.(*opensea.OSResponse)

		// Send price check message
		message := msg.PriceCheckMessage(osResponse.Collection.Name, opensea.CreateUrlFromSlug(slugQuery), osResponse)
		msg.SendMessage(bot, chatID, message)
		return
	}

	// Else query web
	osResponse, err := opensea.QueryAPI(slugQuery)

	// Slugs not found
	//if err == custerror.InvalidSlugErr {
	//	matches := opensea.FindClosestmatch(slugQuery, 3)
	//	// Send clarification message
	//	message := "Collection not found, did you mean any of these:\n\n" + strings.Join(matches, ", ")
	//	msg.SendMessage(bot, chatID, message)
	//	return
	//}

	if err != nil {
		message := "Collection does not exist LOL!"
		msg.SendMessage(bot, chatID, message)
		log.Printf("[opensea.QueryAPI] %v", err)
		return
	}

	// Update cache - Set the value of the key "slugQuery" to fp with the default expiration time
	datastorage.GlobalCache.Set(slugQuery, osResponse, cache.DefaultExpiration)

	// Update DB if collection is not already in database and is popular
	retSlug, err := db.GetSlug(pgdb, slugQuery)
	if err == pg.ErrNoRows && osResponse.Collection.Stats.NumOwners > POPULAR_COLLECTION_NUM_OWNERS {
		log.Printf(fmt.Sprintf("[db.GetSlug] New popular slug... adding slug to DB... %s"), slugQuery)
		_, _ = db.CreateSlug(pgdb, &db.Slugs{
			SlugName:   slugQuery,
			FloorPrice: osResponse.Collection.Stats.FloorPrice,
		})
	}

	// If collection metadata is outdated
	if retSlug.FloorPrice != osResponse.Collection.Stats.FloorPrice {
		updatedSlug, err := db.UpdateSlug(pgdb, &db.Slugs{
			SlugName:   osResponse.Collection.Slug,
			FloorPrice: osResponse.Collection.Stats.FloorPrice,
		})
		if err != nil {
			log.Printf(fmt.Sprintf("[db.UpdateSlug] Err: %s", err))
		}
		log.Printf(fmt.Sprintf("[db.UpdateSlug] Successfully updated slug: %s", updatedSlug))
	}

	// Send price check message
	message := msg.PriceCheckMessage(osResponse.Collection.Name, opensea.CreateUrlFromSlug(slugQuery), osResponse)
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
