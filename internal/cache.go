package internal

import (
	"sync"
	"time"
)

type Cache struct {
	mu      sync.Mutex
	entries map[string]CacheEntry
}

type CacheEntry struct {
	createdAt time.Time
	val       []byte
}

func NewCache(interval time.Duration) *Cache {
	c := &Cache{
		entries: make(map[string]CacheEntry),
	}
	go c.reapLoop(interval)

	return c
}

func (c *Cache) Get(key string) ([]byte, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	entry, ok := c.entries[key]
	return entry.val, ok
}

func (c *Cache) Add(key string, val []byte) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.entries[key] = CacheEntry{createdAt: time.Now(), val: val}
}

func (c *Cache) reapLoop(interval time.Duration) {
	ticker := time.NewTicker(interval)
	for range ticker.C {
		c.mu.Lock()
		for k, v := range c.entries {
			if time.Since(v.createdAt) > interval {
				delete(c.entries, k)
			}
		}
		c.mu.Unlock()
	}
}
