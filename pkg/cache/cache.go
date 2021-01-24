package cache

import (
	"errors"
	"sync"
)

var cache *Cache

type CacheKey interface{}
type CacheValue interface{}

type Cache struct {
	cache map[CacheKey]CacheValue
	sync.RWMutex
}

func InitCacheManagement() {
	cache = &Cache{
		cache: make(map[CacheKey]CacheValue),
	}
}

func PutConfigToCache(key CacheKey, val CacheValue) error {
	return cache.putConfigCache(key, val)
}

func GetConfigFromCache(key CacheKey) (CacheKey, error) {
	return cache.getConfigCache(key)
}

func (cache *Cache) getConfigCache(key CacheKey) (CacheValue, error) {
	cache.RLock()
	value, ok := cache.cache[key]
	cache.RUnlock()
	if ok {
		return value, nil
	}
	return nil, errors.New("not found")
}

func (cache *Cache) putConfigCache(key CacheKey, val CacheValue) error {
	cache.Lock()
	defer cache.Unlock()
	cache.cache[key] = val
	return nil
}
