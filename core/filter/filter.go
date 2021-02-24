package filter

import (
	"errors"
	log "github.com/mhchlib/logger"
	"github.com/mhchlib/mconfig/core/mconfig"
)

func InitFilterEngine() {
	initCache()
	initEvent()
}

func GetEffectEnvKey(appkey mconfig.AppKey, metatdata map[string]string) (mconfig.ConfigEnv, error) {
	effectFilterCacheKey := &EffectFilterCacheKey{
		AppKey:      appkey,
		MetadataMd5: mconfig.GetInterfaceMd5(metatdata),
	}
	//effect cache
	effectFilterCacheVal, ok := getFromEffectFilterCache(effectFilterCacheKey)
	if ok {
		return effectFilterCacheVal, nil
	}
	filters, err := getFilterByAppKey(appkey)
	if err != nil {
		return "", err
	}
	//calculate effect env
	envKey, err := calculateEffectEnvKey(filters, metatdata)
	if err != nil {
		return "", err
	}
	putEffectFilterCache(effectFilterCacheKey, envKey)
	return envKey, nil
}

func calculateEffectEnvKey(filters []*mconfig.FilterEntity, metatdata map[string]string) (mconfig.ConfigEnv, error) {
	var effecfEnv mconfig.ConfigEnv
	maxWeight := -1
	for _, filter := range filters {
		if maxWeight >= filter.Weight {
			continue
		}
		calResult := false
		if filter.Code == "" {
			calResult = true
			goto result
		}
		switch filter.Mode {
		case mconfig.FilterMode_lua:
			calResult = CalLuaFilter(string(filter.Code), metatdata)
		case mconfig.FilterMode_simple:
			calResult = CalSimpleFilter(string(filter.Code), metatdata)
		case mconfig.FilterMode_mep:
			calResult = CalMepFilter(string(filter.Code), metatdata)
		default:
			log.Error("not support filter mode", filter.Mode, "in env", filter.Env)
		}
		log.Debug(string(filter.Env), filter.Weight, string(filter.Code), metatdata, calResult)

	result:
		if calResult {
			maxWeight = filter.Weight
			effecfEnv = filter.Env
		}
	}
	if effecfEnv == "" {
		return "", errors.New("not found effect env")
	}
	return effecfEnv, nil
}
