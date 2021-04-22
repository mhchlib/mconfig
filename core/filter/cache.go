package filter

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

// FilterCacheKey ...
type FilterCacheKey struct {
	AppKey mconfig.AppKey
	Env    mconfig.ConfigEnv
}

// FilterCacheValue ...
type FilterCacheValue struct {
	Weight int
	Code   mconfig.FilterVal
	Mode   mconfig.FilterMode
	mconfig.DataVersion
}

var filterCache cache.Cache
var effectFilterCache cache.Cache

// EffectFilterCacheKey ...
type EffectFilterCacheKey struct {
	AppKey      mconfig.AppKey `json:"app_key"`
	MetadataMd5 string         `json:"metadata"`
}

func initCache() {
	filterCache = cache.NewCache()
	effectFilterCache = cache.NewLRUCache(20)
}

func putEffectFilterCache(key *EffectFilterCacheKey, value mconfig.ConfigEnv) {
	_ = effectFilterCache.PutCache(*key, value)
}

func getFromEffectFilterCache(key *EffectFilterCacheKey) (mconfig.ConfigEnv, bool) {
	value, err := effectFilterCache.GetCache(*key)
	if err != nil {
		return "", false
	}
	return mconfig.ConfigEnv(fmt.Sprintf("%v", value)), true
}

func deleteEffectFilterCacheWithAppKey(appKey mconfig.AppKey) error {
	return effectFilterCache.ExecuteForEachItem(func(key cache.CacheKey, value cache.CacheValue, param ...interface{}) {
		cacheKey := key.(EffectFilterCacheKey)
		if cacheKey.AppKey == appKey {
			_ = effectFilterCache.DeleteCache(key)
		}
	})

}

// PutFilterToCache ...
func PutFilterToCache(appKey mconfig.AppKey, env mconfig.ConfigEnv, val *mconfig.StoreVal) error {
	key := &FilterCacheKey{
		AppKey: appKey,
		Env:    env,
	}
	exist := filterCache.CheckExist(*key)
	if exist {
		value, err := filterCache.GetCache(*key)
		if err != nil {
			log.Error(err)
			return err
		}
		cacheValue, ok := value.(*FilterCacheValue)
		if !ok {
			log.Error("filter cache value transform fail:", fmt.Sprintf("%+v", value))
			return nil
		}
		if val.Version < cacheValue.Version {
			log.Info("filter update version", val.Version, "is smaller than cache version", cacheValue.Version)
			return nil
		}
		if val.Version == cacheValue.Version {
			if val.Md5 == cacheValue.Md5 {
				log.Info("filter update md5", val.Md5, "is equal with cache md5", cacheValue.Md5)
				return nil
			}
		}
	}
	storeVal, err := mconfig.TransformMap2FilterStoreVal(val.Data)
	if err != nil {
		return err
	}

	err = filterCache.PutCache(*key, &FilterCacheValue{
		Weight: storeVal.Weight,
		Code:   storeVal.Code,
		Mode:   storeVal.Mode,
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
			AppKey: appKey,
			Type:   mconfig.Event_Type_Filter,
		},
	})
	if err != nil {
		return err
	}
	return nil
}

// DeleteFilterFromCacheByApp ...
func DeleteFilterFromCacheByApp(appKey mconfig.AppKey) error {
	err := filterCache.ExecuteForEachItem(func(key cache.CacheKey, value cache.CacheValue, param ...interface{}) {
		k := key.(FilterCacheKey)
		if appKey == k.AppKey {
			_ = filterCache.DeleteCache(k)
			log.Info("recycle filter cache with app key:", fmt.Sprintf("%+v", k))
		}
	})
	if err != nil {
		return err
	}
	return nil
}

// DeleteFilterFromCache ...
func DeleteFilterFromCache(appKey mconfig.AppKey, env mconfig.ConfigEnv) error {
	_ = filterCache.DeleteCache(&FilterCacheKey{
		AppKey: appKey,
		Env:    env,
	})
	return nil
}

// GetFilterFromCache ...
func GetFilterFromCache(appKey mconfig.AppKey) ([]*mconfig.FilterEntity, error) {
	filters := make([]*mconfig.FilterEntity, 0)
	mutex := sync.Mutex{}
	err := filterCache.ExecuteForEachItem(func(key cache.CacheKey, value cache.CacheValue, param ...interface{}) {
		k := key.(FilterCacheKey)
		v := value.(*FilterCacheValue)
		if appKey == k.AppKey {
			mutex.Lock()
			filters = append(filters, &mconfig.FilterEntity{
				Env:    k.Env,
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
			val, err := mconfig.TransformMap2FilterStoreVal(filter.Data)
			if err != nil {
				return nil, err
			}
			_ = PutFilterToCache(appKey, val.Env, filter)
		}
		filters, _ = GetFilterFromCache(appKey)
	}
	//for _, filter := range filters {
	//	log.Info(fmt.Sprintf("%v", filter))
	//}
	return filters, nil
}

// CheckCacheUpToDateWithStore ...
func CheckCacheUpToDateWithStore() error {
	return filterCache.ExecuteForEachItem(func(key cache.CacheKey, value cache.CacheValue, param ...interface{}) {
		cacheKey := key.(FilterCacheKey)
		cacheValue := value.(*FilterCacheValue)
		filterVal, err := store.GetFilterVal(cacheKey.AppKey, cacheKey.Env)
		if err != nil {
			if errors.Is(err, mconfig.ERROR_STORE_NOT_FOUND) {
				_ = DeleteFilterFromCache(cacheKey.AppKey, cacheKey.Env)
			} else {
				log.Error(fmt.Sprintf("cron sync filter -- store get config val %v fail:", cacheKey), err.Error())
				return
			}
		}
		//log.Info(cacheValue.Version, filterVal.Version, cacheValue.Md5, filterVal.Md5)
		if cacheValue.Version != filterVal.Version || cacheValue.Md5 != filterVal.Md5 {
			err = PutFilterToCache(cacheKey.AppKey, cacheKey.Env, filterVal)
			if err != nil {
				log.Error(fmt.Sprintf("cron sync filter -- store put filter val key: %v value: %v fail:", cacheKey, filterVal), err.Error())
				return
			}

		}
	})
}
