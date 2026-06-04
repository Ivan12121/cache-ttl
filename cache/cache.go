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
	items map[string]Item
	mu    sync.RWMutex
	stop chan struct{}
}

func NewCache(cleanupInterval time.Duration) *Cache {
	cache := &Cache{
		items: make(map[string]Item),
		stop: make(chan struct{}),
	}

	go cache.startCleanup(cleanupInterval)
	return cache
}

func (c *Cache) Set(key string, value any, ttl time.Duration) {
	expiration := time.Now().Add(ttl).UnixNano()

	c.mu.Lock()
	defer c.mu.Unlock()

	c.items[key] = Item{
		Value: value,
		Expiration: expiration,
	}
}

func (c *Cache) Get(key string) (any, bool) {
	c.mu.RLock()
	item, found := c.items[key]
	c.mu.RUnlock()

	if  !found {
		return nil, false
	}

	if time.Now().UnixNano() > item.Expiration {
		c.Delete(key)
		return nil, false
	}

	return item.Value, true 
}

func (c *Cache) Delete(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	delete(c.items, key)
}

func (c *Cache) cleanup() {
	now := time.Now().UnixNano()

	c.mu.Lock()
	defer c.mu.Unlock()

	for key, item := range c.items {
		if now > item.Expiration {
			delete(c.items, key)
		}
	}
}

func (c *Cache) startCleanup(interval time.Duration) {
	ticker := time.NewTicker(interval)

	for {
		select {
		case <- ticker.C:
			c.cleanup()
		case <- c.stop:
			ticker.Stop()
			return
		}
	}
}

func (c *Cache) Stop() {
	close(c.stop)
}