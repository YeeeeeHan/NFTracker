package datastorage

import (
	"github.com/patrickmn/go-cache"
	"time"
)

var GlobalCache *cache.Cache

func InitCache() *cache.Cache {
	// Create a cache with a default expiration time of 1 minutes, and which
	// purges expired items every 2 minutes
	GlobalCache = cache.New(1*time.Minute, 2*time.Minute)

	return GlobalCache
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

//// Update cache - Set the value of the key "slugQuery" to fp with the default expiration time
//datastorage.GlobalCache.Set(slugQuery, osResponse, cache.DefaultExpiration)
