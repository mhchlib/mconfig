package service

import (
	"sync"
)

type MconfigCache struct {
	cache map[ConfigId]ConfigJSONStr
	sync.RWMutex
}

var mconfigCache *MconfigCache

func init() {
	mconfigCache = &MconfigCache{
		cache: make(map[ConfigId]ConfigJSONStr),
	}
}

func (cache *MconfigCache) getConfigCache(key ConfigId) (ConfigJSONStr, error) {
	cache.RLock()
	value, ok := cache.cache[key]
	cache.RUnlock()
	if ok {
		return value, nil
	}
	return "", nil
}

func (cache *MconfigCache) putConfigCache(key ConfigId, value ConfigJSONStr) error {
	cache.Lock()
	defer cache.Unlock()
	cache.cache[key] = value
	return nil
}
