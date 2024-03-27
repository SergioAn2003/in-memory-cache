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
	cache := &Cache[K, T]{
		cache: make(map[K]CacheItem[T]),
		mu:    sync.RWMutex{},
	}

	go func() {
		for {
			cache.Clear()
			time.Sleep(time.Minute)
		}
	}()

	return cache
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

func (c *Cache[K, T]) Get(key K) (T, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	cache, exist := c.cache[key]
	return cache.value, exist
}

func (c *Cache[K, T]) Delete(key K) {
	c.mu.Lock()
	delete(c.cache, key)
	c.mu.Unlock()
}

func (c *Cache[K, T]) Clear() {
	for key, value := range c.cache {
		if time.Now().After(value.lifeTime) {
			delete(c.cache, key)
		}
	}
}
