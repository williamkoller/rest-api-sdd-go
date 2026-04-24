package memory

import (
	"context"
	"errors"
	"strings"
	"sync"
	"time"
)

var ErrCacheMiss = errors.New("cache: key not found")

type entry struct {
	value     []byte
	expiresAt time.Time
}

type MemoryCache struct {
	mu    sync.RWMutex
	items map[string]entry
}

func New() *MemoryCache {
	c := &MemoryCache{items: make(map[string]entry)}
	go c.evict()
	return c
}

func (c *MemoryCache) Get(_ context.Context, key string) ([]byte, error) {
	c.mu.RLock()
	e, ok := c.items[key]
	c.mu.RUnlock()
	if !ok || time.Now().After(e.expiresAt) {
		return nil, ErrCacheMiss
	}
	return e.value, nil
}

func (c *MemoryCache) Set(_ context.Context, key string, value []byte, ttl time.Duration) error {
	c.mu.Lock()
	c.items[key] = entry{value: value, expiresAt: time.Now().Add(ttl)}
	c.mu.Unlock()
	return nil
}

func (c *MemoryCache) Delete(_ context.Context, key string) error {
	c.mu.Lock()
	delete(c.items, key)
	c.mu.Unlock()
	return nil
}

func (c *MemoryCache) DeletePattern(_ context.Context, pattern string) error {
	prefix := strings.TrimSuffix(pattern, "*")
	c.mu.Lock()
	for k := range c.items {
		if strings.HasPrefix(k, prefix) {
			delete(c.items, k)
		}
	}
	c.mu.Unlock()
	return nil
}

func (c *MemoryCache) evict() {
	ticker := time.NewTicker(time.Minute)
	for range ticker.C {
		now := time.Now()
		c.mu.Lock()
		for k, e := range c.items {
			if now.After(e.expiresAt) {
				delete(c.items, k)
			}
		}
		c.mu.Unlock()
	}
}
