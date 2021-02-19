package cache

import (
	"sync"
)

type CacheKey interface{}
type CacheValue interface{}

type Cache struct {
	cache map[CacheKey]CacheValue
	sync.RWMutex
}

func NewCache() *Cache {
	return &Cache{
		cache: make(map[CacheKey]CacheValue),
	}
}

func (cache *Cache) GetCache(key CacheKey) (CacheValue, error) {
	cache.RLock()
	value, ok := cache.cache[key]
	cache.RUnlock()
	if ok {
		return value, nil
	}
	return nil, ERROR_CACHE_NOT_FOUND
}

func (cache *Cache) PutCache(key CacheKey, val CacheValue) error {
	cache.Lock()
	defer cache.Unlock()
	cache.cache[key] = val
	return nil
}

func (cache *Cache) GetCacheMap() map[CacheKey]CacheValue {
	cloneVal := make(map[CacheKey]CacheValue)
	cache.Lock()
	for key, value := range cache.cache {
		cloneVal[key] = value
	}
	cache.Unlock()
	return cloneVal
}
