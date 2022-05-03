package handlers

import (
	"NFTracker/cmd/message"
	"NFTracker/pkg/db"
	"NFTracker/pkg/opensea"
	"fmt"
	"github.com/go-pg/pg/v10"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
)

const PopularCollectionNumOwners = 400

func Introduction(bot *tgbotapi.BotAPI, chatID int64) {
	msg := tgbotapi.NewMessage(chatID, message.WelcomeMessage)
	msg.ParseMode = "Markdown"
	msg.DisableWebPagePreview = true
	if _, e := bot.Send(msg); e != nil {
		log.Printf("Error sending message to telegram.\nMessage: %v\nError: %v", msg, e)
	}
	return
}

func PriceCheck(pgdb *pg.DB, bot *tgbotapi.BotAPI, chatID int64, userName, slugQuery string) {
	// empty slugQuery
	if slugQuery == "" {
		msg := "No slugQuery detected."
		message.SendMessage(bot, chatID, msg)
		return
	}

	//// Check Cache, if found return from cache
	//if x, found := datastorage.GlobalCache.Get(slugQuery); found {
	//	// assert type
	//	osResponse := x.(*opensea.OSResponse)
	//
	//	// Send price check message
	//	message := message.PriceCheckMessage(osResponse.Collection.Name, opensea.CreateUrlFromSlug(slugQuery), osResponse)
	//	message.SendMessage(bot, chatID, message)
	//	return
	//}

	// Else query web
	osResponse, err := opensea.QueryAPI(slugQuery)
	if err != nil {
		// Find the closest match from list of popular NFTs from database
		match, err := closestMatchHelper(pgdb, slugQuery)
		if err != nil {
			log.Printf(fmt.Sprintf("[closestMatchHelper] Err: %s", err))
			return
		}

		// Send notification of 404 but suggest the closest match
		msg := fmt.Sprintf("⚠️ \"%s\" does not exist, please double-check `<slug>`", slugQuery)
		message.SendInlineSlugMissMessage(bot, chatID, msg, slugQuery, match)

		return
	}

	defer func() {
		// Collecting customer usernames
		_, err = db.GetCustomer(pgdb, userName)
		if err == pg.ErrNoRows {
			log.Printf("[db.GetCustomer] New user... adding user to DB...")
			_, _ = db.CreateCustomer(pgdb, userName)
		}
	}()

	//// Update cache - Set the value of the key "slugQuery" to fp with the default expiration time
	//datastorage.GlobalCache.Set(slugQuery, osResponse, cache.DefaultExpiration)

	// If collection is popular, trust user's query and return osresponse
	if osResponse.Collection.Stats.NumOwners > PopularCollectionNumOwners {
		// Send price check message
		msg := message.PriceCheckMessage(osResponse.Collection.Name, opensea.CreateOpenseaUrlFromSlug(slugQuery), osResponse)
		message.SendMessage(bot, chatID, msg)

		// Update DB with popular collection
		popularCollectionHelper(pgdb, osResponse, slugQuery)

		return

		// If collection is not popular, return query but suggest the closest match
	} else {
		// Return query
		msg := message.PriceCheckMessage(osResponse.Collection.Name, opensea.CreateOpenseaUrlFromSlug(slugQuery), osResponse)

		// But suggest the closest match from list of popular NFTs from database
		match, err := closestMatchHelper(pgdb, slugQuery)
		if err != nil {
			log.Printf(fmt.Sprintf("[closestMatchHelper] Err: %s", err))
			return
		}

		message.SendInlineSlugMissMessage(bot, chatID, msg, slugQuery, match)
		return
	}
}

func EditMessage(pgdb *pg.DB, bot *tgbotapi.BotAPI, messageID int, chatID int64, slugQuery string) {
	// Query web
	osResponse, err := opensea.QueryAPI(slugQuery)
	if err != nil {
		msg := "Collection does not exist!"
		message.SendMessage(bot, chatID, msg)
		log.Printf("[opensea.QueryAPI] %v", err)
		return
	}

	// Send price check message
	msg := message.PriceCheckMessage(osResponse.Collection.Name, opensea.CreateOpenseaUrlFromSlug(slugQuery), osResponse)
	message.SendEditMessage(bot, chatID, messageID, msg)
}
