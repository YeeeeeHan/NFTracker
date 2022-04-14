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
