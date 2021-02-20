package config

import (
	"errors"
	"fmt"
	log "github.com/mhchlib/logger"
	"github.com/mhchlib/mconfig/pkg/cache"
	"github.com/mhchlib/mconfig/pkg/mconfig"
	"github.com/mhchlib/mconfig/pkg/store"
	"sync"
)

type ConfigCacheKey struct {
	appKey    mconfig.AppKey
	configKey mconfig.ConfigKey
	env       mconfig.ConfigEnv
}

var configCache *cache.Cache
var appRegisterCache *cache.Cache

var registerLock sync.Mutex

func initCache() {
	configCache = cache.NewCache()
	appRegisterCache = cache.NewCache()

	registerLock = sync.Mutex{}

}

func PutConfigToCache(appKey mconfig.AppKey, configKey mconfig.ConfigKey, env mconfig.ConfigEnv, val mconfig.ConfigVal) error {
	key := &ConfigCacheKey{
		appKey:    appKey,
		configKey: configKey,
		env:       env,
	}
	return configCache.PutCache(*key, val)
}

func GetConfigFromCache(appKey mconfig.AppKey, configKey mconfig.ConfigKey, env mconfig.ConfigEnv) (mconfig.ConfigVal, error) {
	key := &ConfigCacheKey{
		appKey:    appKey,
		configKey: configKey,
		env:       env,
	}
	c, err := configCache.GetCache(*key)
	if err != nil {
		return "", err
	}
	return mconfig.ConfigVal(fmt.Sprintf("%v", c)), nil
}

func DeleteConfigFromCacheByApp(appKey mconfig.AppKey) error {
	err := configCache.ExecuteForEachItem(func(key cache.CacheKey, value cache.CacheValue, param ...interface{}) {
		k := key.(ConfigCacheKey)
		if appKey == k.appKey {
			_ = configCache.DeleteCache(k)
			log.Info("recycle config cache with app key:", fmt.Sprintf("%+v", k))
		}
	})
	if err != nil {
		return err
	}
	return nil
}

func GetConfig(appKey mconfig.AppKey, configKeys []mconfig.ConfigKey, env mconfig.ConfigEnv) ([]*mconfig.ConfigEntity, error) {
	configs := make([]*mconfig.ConfigEntity, 0)
	for _, configKey := range configKeys {
		val, err := GetConfigFromCache(appKey, configKey, env)
		if err != nil {
			val, err = store.GetCurrentMConfigStore().GetConfigVal(appKey, configKey, env)
			if err != nil {
				return nil, err
			}
			//sync to store
			go func() {
				err := PutConfigToCache(appKey, configKey, env, val)
				if err != nil {
					log.Info(err)
				}
			}()
		}
		configs = append(configs, &mconfig.ConfigEntity{
			Key: configKey,
			Val: val,
		})
	}
	return configs, nil
}

func RegisterAppNotify(app mconfig.AppKey) error {
	registerLock.Lock()
	defer registerLock.Unlock()
	v, err := appRegisterCache.GetCache(app)
	count := 0
	if err != nil && !errors.Is(err, cache.ERROR_CACHE_NOT_FOUND) {
		return err
	}
	if v == nil {
		count = 0
	} else {
		count = v.(int)
	}
	count = count + 1
	return appRegisterCache.PutCache(app, count)
}

func UnRegisterAppNotify(app mconfig.AppKey) error {
	registerLock.Lock()
	defer registerLock.Unlock()
	v, err := appRegisterCache.GetCache(app)
	count := 0
	if err != nil && !errors.Is(err, cache.ERROR_CACHE_NOT_FOUND) {
		return err
	}
	if v == nil {
		return nil
	} else {
		count = v.(int)
	}
	count = count - 1
	if count == 0 {
		return appRegisterCache.DeleteCache(app)
	}
	return appRegisterCache.PutCache(app, count)
}

func CheckRegisterAppNotifyExist(app mconfig.AppKey) bool {
	return appRegisterCache.CheckExist(app)
}
