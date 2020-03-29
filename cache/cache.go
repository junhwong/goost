package cache

import (
	"strings"
	"time"
)

type Loader func() (value interface{}, err error)

type Filter func(key string, entry CacheEntry) bool

type Cache interface {
	// Returns the value associated with key in this cache, obtaining that value from loader if necessary.
	Get(key string, loader Loader) (value interface{}, err error)
	// Returns the value associated with key in this cache, or null if there is no cached value for key.
	GetIfPresent(key string) (value interface{})
	// Discards any cached value for key key.
	Invalidate(filter Filter)
}

type CacheEntry interface {
	Expiration() time.Duration
}

func FilterWithKey(key string) Filter {
	return func(k string, _ CacheEntry) bool {
		return strings.EqualFold(key, k)
	}
}
