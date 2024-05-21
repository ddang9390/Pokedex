package pokecache

import (
	"sync"
	"time"
)

type Cache struct {
	cache map[string]cacheEntry
	mu    *sync.RWMutex
}

type cacheEntry struct {
	createdAt time.Time
	val       []byte
}

func NewCache(interval time.Duration) Cache {
	cache := Cache{
		cache: make(map[string]cacheEntry),
		mu:    &sync.RWMutex{},
	}

	go cache.reapLoop(interval)
	return cache
}

func (cache Cache) Add(key string, val []byte) {
	cache.mu.Lock()
	defer cache.mu.Unlock()

	c := cacheEntry{}
	c.createdAt = time.Now()
	c.val = val
	cache.cache[key] = c

}

func (cache Cache) Get(key string) ([]byte, bool) {
	cache.mu.Lock()
	defer cache.mu.Unlock()
	val, ok := cache.cache[key]

	if ok {
		return val.val, true
	} else {
		return val.val, false
	}

}

func (cache Cache) reapLoop(interval time.Duration) {
	ticker := time.NewTicker(interval)
	for range ticker.C {
		cache.reap(interval)
	}
}

func (cache Cache) reap(interval time.Duration) {
	cache.mu.Lock()
	defer cache.mu.Unlock()

	t := time.Now().Add(-interval)
	for key, val := range cache.cache {
		if val.createdAt.Before(t) {
			delete(cache.cache, key)
		}
	}
}
