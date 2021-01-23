package config

import (
	"github.com/mhchlib/mconfig/pkg/cache"
	"github.com/mhchlib/mconfig/pkg/mconfig"
)

type ConfigCacheKey struct {
	appKey    mconfig.Appkey
	configKey mconfig.ConfigKey
	env       mconfig.ConfigEnv
}

func PutConfigToCache(appKey mconfig.Appkey, configKey mconfig.ConfigKey, env mconfig.ConfigEnv, val mconfig.ConfigVal) error {
	key := &ConfigCacheKey{
		appKey:    appKey,
		configKey: configKey,
		env:       env,
	}
	return cache.PutConfigToCache(*key, val)
}

func GetConfigFromCache(appKey mconfig.Appkey, configKey mconfig.ConfigKey, env mconfig.ConfigEnv) (mconfig.ConfigVal, error) {
	key := &ConfigCacheKey{
		appKey:    appKey,
		configKey: configKey,
		env:       env,
	}
	return cache.GetConfigFromCache(*key)
}
