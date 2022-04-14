package handlers

import (
	"NFTracker/cmd/msg"
	"NFTracker/pkg/db"
	"NFTracker/pkg/opensea"
	"fmt"
	"github.com/go-pg/pg/v10"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
)

const PopularCollectionNumOwners = 400

func Introduction(bot *tgbotapi.BotAPI, chatID int64) {
	msg := tgbotapi.NewMessage(chatID, msg.WelcomeMessage)
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
		message := "No slugQuery detected."
		msg.SendMessage(bot, chatID, message)
		return
	}

	// Collecting customer usernames
	defer func() {
		_, err := db.GetCustomer(pgdb, userName)
		if err == pg.ErrNoRows {
			log.Printf("[db.GetCustomer] New user... adding user to DB...")
			_, _ = db.CreateCustomer(pgdb, userName)
		}
	}()

	//// Check Cache, if found return from cache
	//if x, found := datastorage.GlobalCache.Get(slugQuery); found {
	//	// assert type
	//	osResponse := x.(*opensea.OSResponse)
	//
	//	// Send price check message
	//	message := msg.PriceCheckMessage(osResponse.Collection.Name, opensea.CreateUrlFromSlug(slugQuery), osResponse)
	//	msg.SendMessage(bot, chatID, message)
	//	return
	//}

	// Else query web
	osResponse, err := opensea.QueryAPI(slugQuery)
	if err != nil {
		// Find the closest match from list of popular NFTs from database
		matches, err := closestMatchHelper(pgdb, slugQuery)
		if err != nil {
			log.Printf(fmt.Sprintf("[closestMatchHelper] Err: %s", err))
			return
		}

		// Send notification of 404 but suggest the closest match
		message := fmt.Sprintf("⚠️ \"%s\" does not exist, please double-check `<slug>`", slugQuery)
		msg.SendInlineSlugMissMessage(bot, chatID, message, slugQuery, matches...)

		return
	}

	//// Update cache - Set the value of the key "slugQuery" to fp with the default expiration time
	//datastorage.GlobalCache.Set(slugQuery, osResponse, cache.DefaultExpiration)
	//
	//// Update DB if collection is not already in database and is popular
	//retSlug, err := db.GetSlug(pgdb, slugQuery)
	//if err == pg.ErrNoRows && osResponse.Collection.Stats.NumOwners > PopularCollectionNumOwners {
	//	log.Printf(fmt.Sprintf("[db.GetSlug] New popular slug... adding slug to DB... %s"), slugQuery)
	//	_, _ = db.CreateSlug(pgdb, &db.Slugs{
	//		SlugName:   slugQuery,
	//		FloorPrice: osResponse.Collection.Stats.FloorPrice,
	//	})
	//}
	//
	//// If collection metadata is outdated
	//if retSlug.FloorPrice != osResponse.Collection.Stats.FloorPrice {
	//	updatedSlug, err := db.UpdateSlug(pgdb, &db.Slugs{
	//		SlugName:   osResponse.Collection.Slug,
	//		FloorPrice: osResponse.Collection.Stats.FloorPrice,
	//	})
	//	if err != nil {
	//		log.Printf(fmt.Sprintf("[db.UpdateSlug] Err: %s", err))
	//	}
	//	log.Printf(fmt.Sprintf("[db.UpdateSlug] Successfully updated slug: %s", updatedSlug))
	//}

	// If collection is popular, trust user's query and return osresponse
	if osResponse.Collection.Stats.NumOwners > PopularCollectionNumOwners {
		// Send price check message
		message := msg.PriceCheckMessage(osResponse.Collection.Name, opensea.CreateUrlFromSlug(slugQuery), osResponse)
		msg.SendMessage(bot, chatID, message)

		return
	} else {
		// If collection is not popular, return query but suggest closest match

		// Find the closest match from list of popular NFTs from database
		matches, err := closestMatchHelper(pgdb, slugQuery)
		if err != nil {
			log.Printf(fmt.Sprintf("[closestMatchHelper] Err: %s", err))
			return
		}

		// Send results but suggest closest match
		message := msg.PriceCheckMessage(osResponse.Collection.Name, opensea.CreateUrlFromSlug(slugQuery), osResponse)
		msg.SendInlineSlugMissMessage(bot, chatID, message, slugQuery, matches...)
		return
	}

	// Send price check message
	message := msg.PriceCheckMessage(osResponse.Collection.Name, opensea.CreateUrlFromSlug(slugQuery), osResponse)
	msg.SendMessage(bot, chatID, message)

	return
}

func EditMessage(pgdb *pg.DB, bot *tgbotapi.BotAPI, messageID int, chatID int64, slugQuery string) {

	// Query web
	osResponse, err := opensea.QueryAPI(slugQuery)
	if err != nil {
		message := "Collection does not exist LOL!"
		msg.SendMessage(bot, chatID, message)
		log.Printf("[opensea.QueryAPI] %v", err)
		return
	}

	// Send price check message
	message := msg.PriceCheckMessage(osResponse.Collection.Name, opensea.CreateUrlFromSlug(slugQuery), osResponse)
	msg.SendEditMessage(bot, chatID, messageID, message)
}
