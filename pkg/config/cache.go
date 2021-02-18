package config

import (
	"fmt"
	"github.com/mhchlib/mconfig/pkg/cache"
	"github.com/mhchlib/mconfig/pkg/client"
	"github.com/mhchlib/mconfig/pkg/mconfig"
	"github.com/mhchlib/mconfig/pkg/store"
)

type ConfigCacheKey struct {
	appKey    mconfig.AppKey
	configKey mconfig.ConfigKey
	env       mconfig.ConfigEnv
}

var configCache *cache.Cache

func initCache() {
	configCache = cache.NewCache()
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

func GetConfig(appKey mconfig.AppKey, configKeys []mconfig.ConfigKey, env mconfig.ConfigEnv) ([]*mconfig.ConfigEntity, error) {
	configs := make([]*mconfig.ConfigEntity, 0)
	for _, configKey := range configKeys {
		val, err := GetConfigFromCache(appKey, configKey, env)
		if err != nil {
			val, err = store.GetCurrentMConfigStore().GetConfigVal(appKey, configKey, env)
			if err != nil {
				return nil, err
			}
		}
		configs = append(configs, &mconfig.ConfigEntity{
			Key: configKey,
			Val: val,
		})
	}
	return configs, nil
}

func WatchConfig(c *client.Client, appKey mconfig.AppKey, configKeys []mconfig.ConfigKey, env mconfig.ConfigEnv) error {
	err := c.BuildClientConfigRelation(appKey, configKeys, env)
	if err != nil {
		return err
	}
	return nil
}
