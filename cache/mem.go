package cache

import (
	"errors"
	"sync"
	"time"
)

type entry struct {
	// key        string
	expiration time.Duration
	value      interface{}
}

// func (e *entry) Key() string {
// 	if e == nil {
// 		return ""
// 	}
// 	return e.key
// }
func (e *entry) Value() interface{} {
	if e == nil {
		return nil
	}
	return e.value
}
func (e *entry) IsExpired() bool {
	if e == nil {
		return true
	}
	if e.expiration <= 0 {
		return false
	}
	t := time.Now()
	return t.After(t.Add(e.expiration))
}

type threadSafeMap struct {
	lock  sync.RWMutex
	items map[string]*entry
}

func NewMemCache() Cache {
	return &threadSafeMap{
		items: map[string]*entry{},
	}
}

func (c *threadSafeMap) update(key string, value interface{}, expiration time.Duration) (interface{}, error) {
	if value == nil {
		return c.delete(key)
	}
	item := c.items[key]
	if item == nil {
		item = &entry{} //key: key
		c.items[key] = item
	}

	item.expiration = expiration
	item.value = value
	return item.value, nil
}

func (c *threadSafeMap) delete(key string) (interface{}, error) {
	item := c.items[key]
	delete(c.items, key)
	if item == nil {
		return nil, nil
	}
	return item.value, nil
}

func (c *threadSafeMap) Set(key string, value interface{}, expiration time.Duration) (interface{}, error) {
	c.lock.Lock()
	defer c.lock.Unlock()

	return c.update(key, value, expiration)
}
func (c *threadSafeMap) get(key string) (*entry, error) {
	c.lock.RLock()
	defer c.lock.RUnlock()

	return c.items[key], nil
}

func (c *threadSafeMap) Get(key string, loader Loader) (value interface{}, err error) {
	item, err := c.get(key)
	if err != nil {
		return nil, err
	}

	if item.IsExpired() {
		if loader == nil {
			return nil, errors.New("loader cannt be nil")
		}
		value, expiration, err := loader()
		if err != nil || value == nil {
			return nil, err
		}
		return c.Set(key, value, expiration)
	}
	return item.value, nil
}

func (c *threadSafeMap) GetIfPresentE(key string) (value interface{}, err error) {
	item, err := c.get(key)
	if err != nil || item == nil {
		return nil, err
	}
	if item.IsExpired() {
		return c.DeleteE(key)
	}
	return item.value, nil
}

func (c *threadSafeMap) GetIfPresent(key string) (value interface{}) {
	value, _ = c.GetIfPresentE(key)
	return
}

func (c *threadSafeMap) DeleteE(key string) (value interface{}, err error) {
	c.lock.Lock()
	defer c.lock.Unlock()

	return c.delete(key)
}

func (c *threadSafeMap) Delete(key string) (value interface{}) {
	value, _ = c.DeleteE(key)
	return
}

func (c *threadSafeMap) List(filter Filter) error {
	// TODO: 效率
	c.lock.RLock()
	copied := make(map[string]*entry, len(c.items))
	for k, v := range c.items {
		copied[k] = v
	}
	c.lock.RUnlock()
	for k, v := range copied {
		if stop := filter(k, v); stop {
			return nil
		}
	}
	return nil
}
