package cache

import (
	"github.com/mhchlib/mconfig/pkg/config"
	"github.com/mhchlib/mconfig/pkg/mconfig"
	"sync"
)

// MconfigCache ...
type MconfigCache struct {
	cache map[mconfig.Appkey]*config.AppConfigsMap
	sync.RWMutex
}

func (cache *MconfigCache) getConfigCache(key mconfig.Appkey) (*config.AppConfigsMap, error) {
	cache.RLock()
	value, ok := cache.cache[key]
	cache.RUnlock()
	if ok {
		return value, nil
	}
	return nil, mconfig.Error_AppConfigNotFound
}

func (cache *MconfigCache) putConfigCache(key mconfig.Appkey, configs *config.AppConfigs) error {
	configsMap := &config.AppConfigsMap{
		AppConfigs: configs,
	}
	cache.Lock()
	defer cache.Unlock()
	cache.cache[key] = configsMap
	return nil
}
