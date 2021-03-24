package config

import (
	"errors"
	"fmt"
	log "github.com/mhchlib/logger"
	"github.com/mhchlib/mconfig/core/cache"
	"github.com/mhchlib/mconfig/core/event"
	"github.com/mhchlib/mconfig/core/mconfig"
	"github.com/mhchlib/mconfig/core/store"
	"sync"
)

// ConfigCacheKey ...
type ConfigCacheKey struct {
	AppKey    mconfig.AppKey
	ConfigKey mconfig.ConfigKey
	Env       mconfig.ConfigEnv
}

// ConfigCacheValue ...
type ConfigCacheValue struct {
	Key mconfig.ConfigKey
	Val mconfig.ConfigVal
	mconfig.DataVersion
}

var configCache cache.Cache
var appRegisterCache cache.Cache

var registerLock sync.Mutex

func initCache() {
	configCache = cache.NewCache()
	appRegisterCache = cache.NewCache()
	registerLock = sync.Mutex{}
}

// PutConfigToCache ...
func PutConfigToCache(appKey mconfig.AppKey, configKey mconfig.ConfigKey, env mconfig.ConfigEnv, val *mconfig.StoreVal) error {
	key := &ConfigCacheKey{
		AppKey:    appKey,
		ConfigKey: configKey,
		Env:       env,
	}
	exist := configCache.CheckExist(*key)
	if exist {
		value, err := configCache.GetCache(*key)
		if err != nil {
			log.Error(err)
			return err
		}
		cacheValue, ok := value.(*ConfigCacheValue)
		if !ok {
			log.Error("config cache value transform fail:", fmt.Sprintf("%v", value))
			return nil
		}
		if val.Version < cacheValue.Version {
			log.Info("config update version", val.Version, "is smaller than cache version", cacheValue.Version)
			return nil
		}
		if val.Version == cacheValue.Version {
			if val.Md5 == cacheValue.Md5 {
				log.Info("config update md5", val.Md5, "is equal with cache md5", cacheValue.Md5)
				return nil
			}
		}
	}
	storeVal, err := mconfig.TransformMap2ConfigStoreVal(val.Data)
	if err != nil {
		return err
	}
	err = configCache.PutCache(*key, &ConfigCacheValue{
		Key: storeVal.Key,
		Val: storeVal.Val,
		DataVersion: mconfig.DataVersion{
			Md5:     val.Md5,
			Version: val.Version,
		},
	})
	if err != nil {
		return err
	}
	err = event.AddEvent(&event.Event{
		EventDesc: event.EventDesc{
			EventType: event.Event_Change,
			EventKey:  mconfig.EVENT_KEY_CLIENT_NOTIFY,
		},
		Metadata: mconfig.ClientNotifyEventMetadata{
			AppKey:    appKey,
			ConfigKey: configKey,
			Env:       env,
			Type:      mconfig.Event_Type_Config,
		},
	})
	if err != nil {
		return err
	}
	return nil
}

// GetConfigFromCache ...
func GetConfigFromCache(appKey mconfig.AppKey, configKey mconfig.ConfigKey, env mconfig.ConfigEnv) (*mconfig.ConfigEntity, error) {
	key := &ConfigCacheKey{
		AppKey:    appKey,
		ConfigKey: configKey,
		Env:       env,
	}
	value, err := configCache.GetCache(*key)
	if err != nil {
		return nil, err
	}
	cacheVal := value.(*ConfigCacheValue)
	return &mconfig.ConfigEntity{
		Key: cacheVal.Key,
		Val: cacheVal.Val,
	}, nil
}

// DeleteConfigFromCacheByApp ...
func DeleteConfigFromCacheByApp(appKey mconfig.AppKey) error {
	err := configCache.ExecuteForEachItem(func(key cache.CacheKey, value cache.CacheValue, param ...interface{}) {
		k := key.(ConfigCacheKey)
		if appKey == k.AppKey {
			_ = configCache.DeleteCache(k)
			log.Info("recycle config cache with app key:", fmt.Sprintf("%+v", k))
		}
	})
	if err != nil {
		return err
	}
	return nil
}

// DeleteConfigFromCache ...
func DeleteConfigFromCache(appKey mconfig.AppKey, configKey mconfig.ConfigKey, env mconfig.ConfigEnv) error {
	_ = configCache.DeleteCache(&ConfigCacheKey{
		AppKey:    appKey,
		ConfigKey: configKey,
		Env:       env,
	})
	return nil
}

// GetConfig ...
func GetConfig(appKey mconfig.AppKey, configKeys []mconfig.ConfigKey, env mconfig.ConfigEnv) ([]*mconfig.ConfigEntity, error) {
	configs := make([]*mconfig.ConfigEntity, 0)
	//when configKeys len is 0, will get all config
	if len(configKeys) == 0 {
		keys, err := store.GetAppConfigKeys(appKey, env)
		if err != nil {
			return nil, err
		}
		configKeys = keys
	}

	for _, configKey := range configKeys {
		cacheVal, err := GetConfigFromCache(appKey, configKey, env)
		if err != nil {
			storeVal, err := store.GetConfigVal(appKey, configKey, env)
			if err != nil {
				log.Error(fmt.Sprintf("store get config val %v fail:", configKey), err.Error())
				return nil, err
			}
			//put to store
			err = PutConfigToCache(appKey, configKey, env, storeVal)
			if err != nil {
				log.Error(fmt.Sprintf("store put config val key: %v value: %v fail:", configKey, storeVal), err.Error())
			}
			cacheVal, _ = GetConfigFromCache(appKey, configKey, env)
		}
		configs = append(configs, cacheVal)
	}
	return configs, nil
}

// RegisterAppNotify ...
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

// UnRegisterAppNotify ...
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

// CheckRegisterAppNotifyExist ...
func CheckRegisterAppNotifyExist(app mconfig.AppKey) bool {
	return appRegisterCache.CheckExist(app)
}

// CheckCacheUpToDateWithStore ...
func CheckCacheUpToDateWithStore() error {
	return configCache.ExecuteForEachItem(func(key cache.CacheKey, value cache.CacheValue, param ...interface{}) {
		cacheKey := key.(ConfigCacheKey)
		cacheValue := value.(*ConfigCacheValue)
		storeVal, err := store.GetConfigVal(cacheKey.AppKey, cacheKey.ConfigKey, cacheKey.Env)
		if err != nil {
			if errors.Is(err, mconfig.ERROR_STORE_NOT_FOUND) {
				_ = DeleteConfigFromCache(cacheKey.AppKey, cacheKey.ConfigKey, cacheKey.Env)
			} else {
				log.Error(fmt.Sprintf("cron sync config -- store get config val %v fail:", cacheKey), err.Error())
				return
			}
		}
		if storeVal.Version != cacheValue.Version || storeVal.Md5 != cacheValue.Md5 {
			//put to store
			err = PutConfigToCache(cacheKey.AppKey, cacheKey.ConfigKey, cacheKey.Env, storeVal)
			if err != nil {
				log.Error(fmt.Sprintf("cron sync config -- store put config val key: %v value: %v fail:", cacheKey, storeVal), err.Error())
				return
			}
		}
	})
}
