package cache

import (
	"sync"
	"time"
)

type Item struct {
	Value      any
	Expiration int64
}

type Cache struct {
	sync.RWMutex
	items map[string]Item
}

// New creates a new Cache with a background eviction routine running at the interval of the default TTL.
func New(ttl time.Duration) *Cache {
	c := &Cache{
		items: make(map[string]Item),
	}
	if ttl > 0 {
		go c.startJanitor(ttl)
	}
	return c
}

func (c *Cache) startJanitor(interval time.Duration) {
	ticker := time.NewTicker(interval)
	for range ticker.C {
		c.Lock()
		now := time.Now().UnixNano()
		for k, item := range c.items {
			if item.Expiration > 0 && now > item.Expiration {
				delete(c.items, k)
			}
		}
		c.Unlock()
	}
}

func (c *Cache) Set(key string, value any, duration time.Duration) {
	var expiration int64
	if duration > 0 {
		expiration = time.Now().Add(duration).UnixNano()
	}
	c.Lock()
	defer c.Unlock()
	c.items[key] = Item{
		Value:      value,
		Expiration: expiration,
	}
}

func (c *Cache) Get(key string) (any, bool) {
	c.RLock()
	defer c.RUnlock()
	item, ok := c.items[key]
	if !ok {
		return nil, false
	}
	if item.Expiration > 0 && time.Now().UnixNano() > item.Expiration {
		return nil, false
	}
	return item.Value, true
}
