package cache

import (
	"github.com/mhchlib/mconfig/pkg/config"
	"github.com/mhchlib/mconfig/pkg/mconfig"
)

var mconfigCache *MconfigCache

func init() {
	mconfigCache = &MconfigCache{
		cache: make(map[mconfig.Appkey]*config.AppConfigsMap),
	}
}

func PutConfigToCache(appKey mconfig.Appkey, configKey mconfig.ConfigKey, env mconfig.ConfigEnv, val mconfig.ConfigVal) error {
	return nil
}

func GetConfigFromCache(appKey mconfig.Appkey, configKey mconfig.ConfigKey, env mconfig.ConfigEnv) error {
	return nil
}

func PutFilterToCache(appKey mconfig.Appkey, configKey mconfig.ConfigKey, val mconfig.FilterVal) error {
	return nil
}

func GetFilterFromCache(appKey mconfig.Appkey, configKey mconfig.ConfigKey) error {
	return nil
}
