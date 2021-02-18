package filter

import (
	"fmt"
	"github.com/mhchlib/mconfig/pkg/cache"
	"github.com/mhchlib/mconfig/pkg/mconfig"
)

type FilterCacheKey struct {
	appKey mconfig.AppKey
	env    mconfig.ConfigEnv
}

type FilterCacheValue struct {
	weight int
	code   mconfig.FilterVal
	mode   FilterMode
}

type FilterEntity struct {
	Env    mconfig.ConfigEnv
	Weight int
	Code   mconfig.FilterVal
	Mode   FilterMode
}

var filterCache *cache.Cache

func initCache() {
	filterCache = cache.NewCache()
}

func PutFilterToCache(appKey mconfig.AppKey, env mconfig.ConfigEnv, val mconfig.FilterVal) error {
	key := &FilterCacheKey{
		appKey: appKey,
		env:    env,
	}
	return filterCache.PutCache(*key, val)
}

func GetFilterFromCache(appKey mconfig.AppKey, env mconfig.ConfigEnv) (mconfig.FilterVal, error) {
	key := &FilterCacheKey{
		appKey: appKey,
		env:    env,
	}
	c, err := filterCache.GetCache(*key)
	if err != nil {
		return "", err
	}
	return mconfig.FilterVal(fmt.Sprintf("%v", c)), nil
}

func getFilterByAppKey(appKey mconfig.AppKey) []*FilterEntity {
	cacheMap := filterCache.GetCacheMap()
	filters := make([]*FilterEntity, 0)
	for key, value := range cacheMap {
		k := key.(FilterCacheKey)
		v := value.(FilterCacheValue)
		if appKey == k.appKey {
			filters = append(filters, &FilterEntity{
				Env:    k.env,
				Weight: v.weight,
				Code:   v.code,
				Mode:   v.mode,
			})
		}
	}
	return filters
}
