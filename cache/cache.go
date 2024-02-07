package cache

import (
	"sync"
	"time"
)

type Cache[K comparable, T any] struct {
	mu    sync.RWMutex
	cache map[K]CacheItem[T]
}

type CacheItem[T any] struct {
	value    T
	lifeTime time.Time
}

func New[K comparable, T any]() *Cache[K, T] {
	return &Cache[K, T]{
		cache: make(map[K]CacheItem[T]),
		mu:    sync.RWMutex{},
	}
}

func (c *Cache[K, T]) Set(key K, value T, ttl time.Duration) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.cache[key] = CacheItem[T]{
		value:    value,
		lifeTime: time.Now().Add(ttl),
	}
	return nil
}

func (c *Cache[K, T]) Get(key K) (t T, exists bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	cache, exist := c.cache[key]
	if !exist || time.Now().After(cache.lifeTime) {
		c.Delete(key)
		return t, false
	}
	return cache.value, true
}

func (c *Cache[K, T]) Delete(key K) {
	c.mu.Lock()
	delete(c.cache, key)
	c.mu.Unlock()
}

func (c *Cache[K, T]) CleanUpExpired() {
	c.mu.Lock()
	defer c.mu.Unlock()

	now := time.Now()

	for key, value := range c.cache {
		if now.After(value.lifeTime) {
			delete(c.cache, key)
		}
	}
}
