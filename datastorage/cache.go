package datastorage

import (
	cache "github.com/patrickmn/go-cache"
	"time"
)

var GlobalCache *cache.Cache

func InitCache() *cache.Cache {
	// Create a cache with a default expiration time of 5 minutes, and which
	// purges expired items every 10 minutes
	GlobalCache = cache.New(5*time.Minute, 10*time.Minute)

	return GlobalCache
}
