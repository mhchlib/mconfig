package filter

import (
	"fmt"
	log "github.com/mhchlib/logger"
	"github.com/mhchlib/mconfig/core/cache"
	"github.com/mhchlib/mconfig/core/mconfig"
	"github.com/mhchlib/mconfig/core/store"
	"sync"
)

type FilterCacheKey struct {
	appKey mconfig.AppKey
	env    mconfig.ConfigEnv
}

type FilterCacheValue struct {
	Weight int
	Code   mconfig.FilterVal
	Mode   mconfig.FilterMode
	mconfig.DataVersion
}

var filterCache *cache.Cache

func initCache() {
	filterCache = cache.NewCache()
}

func PutFilterToCache(appKey mconfig.AppKey, env mconfig.ConfigEnv, val *mconfig.StoreVal) error {
	key := &FilterCacheKey{
		appKey: appKey,
		env:    env,
	}
	exist := filterCache.CheckExist(key)
	if exist {
		value, err := filterCache.GetCache(key)
		if err != nil {
			log.Error(err)
			return err
		}
		cacheValue, ok := value.(FilterCacheValue)
		if !ok {
			log.Error("filter cache value transform fail:", fmt.Sprintf("%+v", value))
			return nil
		}
		if val.Version < cacheValue.Version {
			log.Info("filter update version", val.Version, "is smaller than cache version", cacheValue.Version)
			return nil
		}
	}
	storeVal, err := mconfig.TransformMap2FilterStoreVal(val.Data)
	if err != nil {
		return err
	}
	return filterCache.PutCache(*key, &FilterCacheValue{
		Weight: storeVal.Weight,
		Code:   storeVal.Code,
		Mode:   storeVal.Mode,
		DataVersion: mconfig.DataVersion{
			Md5:     val.Md5,
			Version: val.Version,
		},
	})
}

func DeleteFilterFromCacheByApp(appKey mconfig.AppKey) error {
	err := filterCache.ExecuteForEachItem(func(key cache.CacheKey, value cache.CacheValue, param ...interface{}) {
		k := key.(FilterCacheKey)
		if appKey == k.appKey {
			_ = filterCache.DeleteCache(k)
			log.Info("recycle filter cache with app key:", fmt.Sprintf("%+v", k))
		}
	})
	if err != nil {
		return err
	}
	return nil
}

func GetFilterFromCache(appKey mconfig.AppKey) ([]*mconfig.FilterEntity, error) {
	filters := make([]*mconfig.FilterEntity, 0)
	mutex := sync.Mutex{}
	err := filterCache.ExecuteForEachItem(func(key cache.CacheKey, value cache.CacheValue, param ...interface{}) {
		k := key.(FilterCacheKey)
		v := value.(*FilterCacheValue)
		if appKey == k.appKey {
			mutex.Lock()
			filters = append(filters, &mconfig.FilterEntity{
				Env:    k.env,
				Weight: v.Weight,
				Code:   v.Code,
				Mode:   v.Mode,
			})
			mutex.Unlock()
		}
	})
	if err != nil {
		return nil, err
	}
	if len(filters) == 0 {
		return nil, cache.ERROR_CACHE_NOT_FOUND
	}
	return filters, nil
}

func getFilterByAppKey(appKey mconfig.AppKey) ([]*mconfig.FilterEntity, error) {
	var filters []*mconfig.FilterEntity
	filters, _ = GetFilterFromCache(appKey)
	if filters == nil {
		appFilters, err := store.GetAppFilters(appKey)
		if err != nil {
			return nil, err
		}
		//sync to cache
		for _, filter := range appFilters {
			//val,ok := filter.Data.(mconfig.FilterStoreVal)
			//val := &mconfig.FilterStoreVal{}
			//err := mapstructure.Decode(filter.Data, &val)
			//if err!=nil {
			//	log.Error("filter store value transform fail:",fmt.Sprintf("%+v",filter.Data),"err:",err)
			//}
			val, err := mconfig.TransformMap2FilterStoreVal(filter.Data)
			if err != nil {
				return nil, err
			}
			_ = PutFilterToCache(appKey, val.Env, filter)
			//filters = append(filters, &mconfig.FilterEntity{
			//	Env:    val.Env,
			//	Weight: val.Weight,
			//	Code:   val.Code,
			//	Mode:   val.Mode,
			//})
		}
		filters, _ = GetFilterFromCache(appKey)
	}
	//for _, filter := range filters {
	//	log.Info(fmt.Sprintf("%v", filter))
	//}
	return filters, nil
}
