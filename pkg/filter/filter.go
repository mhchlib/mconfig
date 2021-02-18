package filter

import (
	log "github.com/mhchlib/logger"
	"github.com/mhchlib/mconfig/pkg/mconfig"
)

type FilterMode string

const (
	FilterMode_lua    FilterMode = "lua"
	FilterMode_simple FilterMode = "simple"
	FilterMode_mep    FilterMode = "mep"
)

func GetEffectEnvKey(appkey mconfig.AppKey, metatdata map[string]string) (mconfig.ConfigEnv, error) {
	filters := getFilterByAppKey(appkey)
	//calculate effect env
	calculateEffectEnvKey(filters, metatdata)

	return "", nil
}

func calculateEffectEnvKey(filters []*FilterEntity, metatdata map[string]string) (mconfig.ConfigEnv, error) {
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
		case FilterMode_lua:
			calResult = CalLuaFilter(string(filter.Code), metatdata)
		case FilterMode_simple:
			calResult = CalSimpleFilter(string(filter.Code), metatdata)
		case FilterMode_mep:
			calResult = CalMepFilter(string(filter.Code), metatdata)
		default:
			log.Error("not support filter mode", filter.Mode, "in env", filter.Env)
		}
	result:
		if calResult {
			maxWeight = filter.Weight
			effecfEnv = filter.Env
		}
	}
	return effecfEnv, nil
}
