package cache

import (
	"strings"
	"time"
)

// IndexFunc knows how to compute the set of indexed values for an object.
type IndexFunc func(obj interface{}) ([]string, error)

//
type Loader func() (value interface{}, expiration time.Duration, err error)

type Filter func(key string, entry CacheEntry) bool

type Cache interface {
	Set(key string, value interface{}, expiration time.Duration) (present interface{}, err error)
	// Returns the value associated with key in this cache, obtaining that value from loader if necessary.
	Get(key string, loader Loader) (value interface{}, err error)
	// Returns the value associated with key in this cache, or nil if there is no cached value for key.
	GetIfPresent(key string) (value interface{})
	GetIfPresentE(key string) (value interface{}, err error)
	List(filter Filter) error
	DeleteE(key string) (value interface{}, err error)
	Delete(key string) interface{}
}

type CacheEntry interface {
	// Key() string
	Value() interface{}
	IsExpired() bool
}

func FilterWithKey(key string) Filter {
	return func(k string, _ CacheEntry) bool {
		return strings.EqualFold(key, k)
	}
}

type Empty struct{}

type String map[string]Empty

// Index maps the indexed value to a set of keys in the store that match on that value
type Index map[string]String

// Indexers maps a name to a IndexFunc
type Indexers map[string]IndexFunc

// Indices maps a name to an Index
type Indices map[string]Index
