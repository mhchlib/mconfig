package cache

import (
	"sync"
)

// GeneralCache ...
type GeneralCache struct {
	cache map[CacheKey]CacheValue
	sync.RWMutex
}

// NewGeneralCache ...
func NewGeneralCache() *GeneralCache {
	return &GeneralCache{
		cache: make(map[CacheKey]CacheValue),
	}
}

// GetCache ...
func (cache *GeneralCache) GetCache(key CacheKey) (CacheValue, error) {
	cache.RLock()
	value, ok := cache.cache[key]
	cache.RUnlock()
	if ok {
		return value, nil
	}
	return nil, ERROR_CACHE_NOT_FOUND
}

// PutCache ...
func (cache *GeneralCache) PutCache(key CacheKey, val CacheValue) error {
	cache.Lock()
	defer cache.Unlock()
	cache.cache[key] = val
	return nil
}

// DeleteCache ...
func (cache *GeneralCache) DeleteCache(key CacheKey) error {
	cache.Lock()
	defer cache.Unlock()
	delete(cache.cache, key)
	return nil
}

// CheckExist ...
func (cache *GeneralCache) CheckExist(key CacheKey) bool {
	cache.RLock()
	defer cache.RUnlock()
	_, ok := cache.cache[key]
	return ok
}

// ExecuteForEachItem ...
func (cache *GeneralCache) ExecuteForEachItem(f func(key CacheKey, value CacheValue, param ...interface{}), param ...interface{}) error {
	var wg sync.WaitGroup
	cache.RLock()
	for key, value := range cache.cache {
		wg.Add(1)
		go func(key interface{}, value interface{}, param ...interface{}) {
			defer wg.Done()
			f(key, value, param)
		}(key, value, param)
	}
	cache.RUnlock()
	wg.Wait()
	return nil
}
