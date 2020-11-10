package pkg

import (
	"sync"
)

type MconfigCache struct {
	cache map[AppId]AppConfigs
	sync.RWMutex
}

var mconfigCache *MconfigCache

func init() {
	mconfigCache = &MconfigCache{
		cache: make(map[AppId]AppConfigs),
	}
}

func (cache *MconfigCache) getConfigCache(key AppId) (AppConfigs, error) {
	cache.RLock()
	value, ok := cache.cache[key]
	cache.RUnlock()
	if ok {
		return value, nil
	}
	return nil, Error_AppConfigNotFound
}

func (cache *MconfigCache) putConfigCache(key AppId, configs AppConfigs) error {
	cache.Lock()
	defer cache.Unlock()
	cache.cache[key] = configs
	return nil
}
