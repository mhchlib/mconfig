package pkg

import (
	"sync"
)

type MconfigCache struct {
	cache map[AppId][]ConfigEntity
	sync.RWMutex
}

var mconfigCache *MconfigCache

func init() {
	mconfigCache = &MconfigCache{
		cache: make(map[AppId][]ConfigEntity),
	}
}

func (cache *MconfigCache) getConfigCache(key AppId) ([]ConfigEntity, error) {
	cache.RLock()
	value, ok := cache.cache[key]
	cache.RUnlock()
	if ok {
		return value, nil
	}
	return nil, nil
}

func (cache *MconfigCache) putConfigCache(key AppId, value ConfigJSONStr) ([]ConfigEntity, error) {
	configs, err := ParseConfigJSONStr(value)
	if err != nil {
		return nil, err
	}
	cache.Lock()
	defer cache.Unlock()
	cache.cache[key] = configs
	return configs, nil
}
