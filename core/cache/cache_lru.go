package cache

import (
	lru "github.com/hashicorp/golang-lru"
	"sync"
)

type LRUCache struct {
	Max int
	Lru *lru.Cache
}

func NewLRUCache(maxArr ...int) *LRUCache {
	max := LRU_MAX
	if len(maxArr) == 1 {
		max = maxArr[0]
	}
	cache, _ := lru.New(max)
	return &LRUCache{
		Max: max,
		Lru: cache,
	}
}

func (cache *LRUCache) GetCache(key CacheKey) (CacheValue, error) {
	value, ok := cache.Lru.Get(key)
	if !ok {
		return nil, ERROR_CACHE_NOT_FOUND
	}
	return value, nil
}

func (cache *LRUCache) PutCache(key CacheKey, val CacheValue) error {
	_ = cache.Lru.Add(key, val)
	return nil
}

func (cache *LRUCache) DeleteCache(key CacheKey) error {
	_ = cache.Lru.Remove(key)
	return nil
}

func (cache *LRUCache) CheckExist(key CacheKey) bool {
	return cache.Lru.Contains(key)
}

func (cache *LRUCache) ExecuteForEachItem(f func(key CacheKey, value CacheValue, param ...interface{}), param ...interface{}) error {
	keys := cache.Lru.Keys()
	var wg sync.WaitGroup
	for _, key := range keys {
		value, ok := cache.Lru.Get(key)
		if ok {
			wg.Add(1)
			go func(key interface{}, value interface{}, param ...interface{}) {
				defer wg.Done()
				f(key, value, param)
			}(key, value, param)
		}
	}
	wg.Wait()
	return nil
}
