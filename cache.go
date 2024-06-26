package cache

import (
	"sync"
	"time"
)

type Cache[K comparable, T any] interface {
	Set(key K, value T, ttl time.Duration) error
	Get(key K) (T, bool)
	Delete(key K)
}

type cache[K comparable, T any] struct {
	mu    sync.RWMutex
	cache map[K]cacheItem[T]
}

type cacheItem[T any] struct {
	value    T
	lifeTime time.Time
}

func New[K comparable, T any]() Cache[K, T] {
	cache := &cache[K, T]{
		cache: make(map[K]cacheItem[T]),
		mu:    sync.RWMutex{},
	}

	go func() {
		for {
			cache.clear()
			time.Sleep(time.Minute)
		}
	}()

	return cache
}

func (c *cache[K, T]) Set(key K, value T, ttl time.Duration) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.cache[key] = cacheItem[T]{
		value:    value,
		lifeTime: time.Now().Add(ttl),
	}
	return nil
}

func (c *cache[K, T]) Get(key K) (T, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	cache, exist := c.cache[key]
	return cache.value, exist
}

func (c *cache[K, T]) Delete(key K) {
	c.mu.Lock()
	delete(c.cache, key)
	c.mu.Unlock()
}

func (c *cache[K, T]) clear() {
	for key, value := range c.cache {
		if time.Now().After(value.lifeTime) {
			delete(c.cache, key)
		}
	}
}
