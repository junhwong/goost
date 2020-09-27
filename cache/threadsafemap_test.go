package cache

import (
	"fmt"
	"testing"
	"time"
)

func TestExpired(t *testing.T) {
	e := entry{expiration: time.Second}
	if e.IsExpired() {
		t.Fatal()
	}
}

func TestList(t *testing.T) {

	c := NewMemCache()
	c.Set("1", "1", 0)
	c.Set("2", "2", 0)
	c.Set("3", "3", 0)

	c.List(func(k string, v CacheEntry) bool {
		fmt.Println(k)
		return false
	})

	c.List(func(k string, v CacheEntry) bool {
		c.Delete(k)
		return false
	})

	c.List(func(k string, v CacheEntry) bool {
		fmt.Println(k)
		return false
	})

}
