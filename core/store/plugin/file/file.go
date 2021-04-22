package file

import (
	"context"
	"encoding/json"
	log "github.com/mhchlib/logger"
	"github.com/mhchlib/mconfig/core/config"
	"github.com/mhchlib/mconfig/core/event"
	"github.com/mhchlib/mconfig/core/filter"
	"github.com/mhchlib/mconfig/core/mconfig"
	"github.com/mhchlib/mconfig/core/store"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/util"
)

// KeyNamespce ...
type KeyNamespce string

// KeyMode ...
type KeyMode string

// KeyClass ...
type KeyClass string

const (
	// PLUGIN_NAME ...
	PLUGIN_NAME = "file"
	// SEPARATOR ...
	SEPARATOR = "/"
	// CLASS_CONFIG ...
	CLASS_CONFIG KeyClass = "config"
	// CLASS_FILTER ...
	CLASS_FILTER KeyClass = "filter"
	// MAX_SUFFIX ...
	MAX_SUFFIX = "zzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzz"
)

var namespce KeyNamespce = "com.github.hchlib.mconfig"

var prefix_common = SEPARATOR + string(namespce)

var prefix_config = prefix_common + SEPARATOR + string(CLASS_CONFIG)
var prefix_filter = SEPARATOR + string(namespce) + SEPARATOR + string(CLASS_FILTER)

// KeyEntity ...
type KeyEntity struct {
	namespace KeyNamespce
	class     KeyClass
	appKey    mconfig.AppKey
	configKey mconfig.ConfigKey
	env       mconfig.ConfigEnv
}

// Event_EventType ...
type Event_EventType int32

const (
	// PUT ...
	PUT Event_EventType = 0
	// DELETE ...
	DELETE Event_EventType = 1
)

// ConfigEvent ...
type ConfigEvent struct {
	Key   []byte
	Value []byte
	Type  Event_EventType
}

// SIZE_CHANGEEVENTBUS ...
const SIZE_CHANGEEVENTBUS = 20

// FileStore ...
type FileStore struct {
	cancelFunc context.CancelFunc
	watchChan  chan *ConfigEvent
	config     *leveldb.DB
	filter     *leveldb.DB
}

// PutConfigVal ...
func (f *FileStore) PutConfigVal(appKey mconfig.AppKey, env mconfig.ConfigEnv, configKey mconfig.ConfigKey, val mconfig.StoreVal) error {
	entity := &KeyEntity{
		namespace: namespce,
		class:     CLASS_CONFIG,
		appKey:    appKey,
		configKey: configKey,
		env:       env,
	}
	key, err := getStoreKey(entity)
	if err != nil {
		return err
	}
	data, _ := json.Marshal(val)
	err = f.config.Put([]byte(key), data, nil)
	return err
}

// PutFilterVal ...
func (f *FileStore) PutFilterVal(appKey mconfig.AppKey, env mconfig.ConfigEnv, val mconfig.StoreVal) error {
	entity := &KeyEntity{
		namespace: namespce,
		class:     CLASS_FILTER,
		appKey:    appKey,
		env:       env,
	}
	key, err := getStoreKey(entity)
	if err != nil {
		return err
	}
	data, _ := json.Marshal(val)
	err = f.filter.Put([]byte(key), data, nil)
	f.watchChan <- &ConfigEvent{
		Key:   []byte(key),
		Value: data,
		Type:  PUT,
	}
	return err
}

// DeleteConfig ...
func (f *FileStore) DeleteConfig(appKey mconfig.AppKey, configKey mconfig.ConfigKey, env mconfig.ConfigEnv) error {
	k := &KeyEntity{
		namespace: namespce,
		class:     CLASS_CONFIG,
		appKey:    appKey,
		configKey: configKey,
		env:       env,
	}
	storeKey, err := getStoreKey(k)
	if err != nil {
		return err
	}
	err = f.config.Delete([]byte(storeKey), nil)
	if err != nil {
		return err
	}
	return nil
}

// DeleteFilter ...
func (f *FileStore) DeleteFilter(appKey mconfig.AppKey, env mconfig.ConfigEnv) error {
	k := &KeyEntity{
		namespace: namespce,
		class:     CLASS_FILTER,
		appKey:    appKey,
		env:       env,
	}
	storeKey, err := getStoreKey(k)
	if err != nil {
		return err
	}
	err = f.filter.Delete([]byte(storeKey), nil)
	if err != nil {
		return err
	}
	return nil
}

// GetAppFilters ...
func (f *FileStore) GetAppFilters(appKey mconfig.AppKey) ([]*mconfig.StoreVal, error) {
	entity := &KeyEntity{
		namespace: namespce,
		class:     CLASS_FILTER,
		appKey:    appKey,
	}
	storeKey, err := getStoreKey(entity)
	if err != nil {
		return nil, err
	}
	filters := []*mconfig.StoreVal{}

	iterator := f.filter.NewIterator(&util.Range{
		Start: []byte(storeKey),
		Limit: []byte(storeKey + MAX_SUFFIX),
	}, nil)
	for iterator.Next() {
		k := string(iterator.Key())
		v := iterator.Value()
		f := &mconfig.StoreVal{}
		err = json.Unmarshal(v, f)
		if err != nil {
			log.Error(err, "key:", k, "value:", string(v))
			return nil, err
		}
		filters = append(filters, f)
	}
	return filters, nil
}

// GetAppConfigs ...
func (f *FileStore) GetAppConfigs(appKey mconfig.AppKey, env mconfig.ConfigEnv) ([]*mconfig.StoreVal, error) {
	entity := &KeyEntity{
		namespace: namespce,
		class:     CLASS_CONFIG,
		appKey:    appKey,
		env:       env,
	}
	storeKey, err := getStoreKey(entity)
	if err != nil {
		return nil, err
	}
	configs := []*mconfig.StoreVal{}

	iterator := f.config.NewIterator(&util.Range{
		Start: []byte(storeKey),
		Limit: []byte(storeKey + MAX_SUFFIX),
	}, nil)
	for iterator.Next() {
		k := string(iterator.Key())
		v := iterator.Value()
		f := &mconfig.StoreVal{}
		err = json.Unmarshal(v, f)
		if err != nil {
			log.Error(err, "key:", k, "value:", string(v))
			return nil, err
		}
		configs = append(configs, f)
	}

	return configs, nil
}

// GetAppConfigKeys ...
func (f *FileStore) GetAppConfigKeys(appKey mconfig.AppKey, env mconfig.ConfigEnv) ([]mconfig.ConfigKey, error) {
	entity := &KeyEntity{
		namespace: namespce,
		class:     CLASS_CONFIG,
		appKey:    appKey,
		env:       env,
	}
	storeKey, err := getStoreKey(entity)
	if err != nil {
		return nil, err
	}
	configs := []mconfig.ConfigKey{}
	iterator := f.config.NewIterator(&util.Range{
		Start: []byte(storeKey),
		Limit: []byte(storeKey + MAX_SUFFIX),
	}, nil)
	for iterator.Next() {
		k := string(iterator.Key())
		keyEntity, err := parseStoreKey(k)
		if err != nil {
			log.Error(err)
		}
		configs = append(configs, keyEntity.configKey)
	}
	return configs, nil
}

// GetSyncData ...
func (f *FileStore) GetSyncData() (mconfig.AppData, error) {
	syncData := make(map[mconfig.AppKey]map[mconfig.ConfigEnv]*mconfig.EnvData)
	iterator := f.config.NewIterator(&util.Range{
		Start: []byte(prefix_common),
		Limit: []byte(prefix_common + MAX_SUFFIX),
	}, nil)
	for iterator.Next() {
		key := iterator.Key()
		storeKey, err := parseStoreKey(string(key))
		if err != nil {
			log.Error(err)
			continue
		}
		appData, ok := syncData[storeKey.appKey]
		if !ok {
			appData = make(map[mconfig.ConfigEnv]*mconfig.EnvData)
			syncData[storeKey.appKey] = appData
		}
		envData, ok := appData[storeKey.env]
		if !ok {
			envData = &mconfig.EnvData{}
			appData[storeKey.env] = envData
		}
		if storeKey.class == CLASS_CONFIG {
			configs := envData.Configs
			if configs == nil {
				configs = make(map[mconfig.ConfigKey]mconfig.StoreVal)
				envData.Configs = configs
			}
			val := &mconfig.StoreVal{}
			err = json.Unmarshal(iterator.Value(), val)
			if err != nil {
				log.Error(err, iterator.Value())
				continue
			}
			configs[storeKey.configKey] = *val
		}
		if storeKey.class == CLASS_FILTER {
			val := &mconfig.StoreVal{}
			err = json.Unmarshal(iterator.Value(), val)
			if err != nil {
				log.Error(err, iterator.Value())
				continue
			}
			envData.Filter = *val
		}
	}

	iterator = f.filter.NewIterator(&util.Range{
		Start: []byte(prefix_common),
		Limit: []byte(prefix_common + MAX_SUFFIX),
	}, nil)
	for iterator.Next() {
		key := iterator.Key()
		storeKey, err := parseStoreKey(string(key))
		if err != nil {
			log.Error(err)
			continue
		}
		appData, ok := syncData[storeKey.appKey]
		if !ok {
			appData = make(map[mconfig.ConfigEnv]*mconfig.EnvData)
			syncData[storeKey.appKey] = appData
		}
		envData, ok := appData[storeKey.env]
		if !ok {
			envData = &mconfig.EnvData{}
			appData[storeKey.env] = envData
		}
		if storeKey.class == CLASS_CONFIG {
			configs := envData.Configs
			if configs == nil {
				configs = make(map[mconfig.ConfigKey]mconfig.StoreVal)
				envData.Configs = configs
			}
			val := &mconfig.StoreVal{}
			err = json.Unmarshal(iterator.Value(), val)
			if err != nil {
				log.Error(err, iterator.Value())
				continue
			}
			configs[storeKey.configKey] = *val
		}
		if storeKey.class == CLASS_FILTER {
			val := &mconfig.StoreVal{}
			err = json.Unmarshal(iterator.Value(), val)
			if err != nil {
				log.Error(err, iterator.Value())
				continue
			}
			envData.Filter = *val
		}
	}

	return syncData, nil
}

// PutSyncData ...
func (e *FileStore) PutSyncData(appData *mconfig.AppData) error {
	for appKey, envData := range *appData {
		for env, data := range envData {
			e.PutFilterVal(appKey, env, data.Filter)
			for configKey, val := range data.Configs {
				e.PutConfigVal(appKey, env, configKey, val)
			}
		}
	}
	return nil
}

// WatchDynamicVal ...
func (e *FileStore) WatchDynamicVal(consumers *store.Consumer) error {
	for {
		select {
		case v, ok := <-e.watchChan:
			if ok == false {
				return store.Error_FAIL_WATCH
			}
			switch v.Type {
			case PUT:
				key, err := parseStoreKey(string(v.Key))
				if err != nil {
					log.Error(err)
					continue
				}
				var metadate interface{}
				var eventKey event.EventKey

				val := &mconfig.StoreVal{}
				err = json.Unmarshal(v.Value, val)
				if err != nil {
					log.Error(err)
					continue
				}
				switch key.class {
				case CLASS_CONFIG:
					metadate = config.ConfigEventMetadata{
						AppKey:    key.appKey,
						ConfigKey: key.configKey,
						Env:       key.env,
						Val:       val,
					}
					eventKey = config.EVENT_KEY
				case CLASS_FILTER:
					metadate = filter.FilterEventMetadata{
						AppKey: key.appKey,
						Env:    key.env,
						Val:    val,
					}
					eventKey = filter.EVENT_KEY
				default:
					log.Error("key class <" + key.class + ">is not declare")
					continue
				}
				err = consumers.AddEvent(&event.Event{
					EventDesc: event.EventDesc{
						EventType: event.Event_Update,
						EventKey:  eventKey,
					},
					Metadata: metadate,
				})
				if err != nil {
					log.Error(err)
				}
				log.Info("etcd update key:", string(v.Key), "value:", string(v.Value))
			case DELETE:
				key, err := parseStoreKey(string(v.Key))
				if err != nil {
					log.Error(err)
					continue
				}
				var metadate interface{}
				var eventKey event.EventKey
				switch key.class {
				case CLASS_CONFIG:
					metadate = config.ConfigEventMetadata{
						AppKey:    key.appKey,
						ConfigKey: key.configKey,
						Env:       key.env,
					}
					eventKey = config.EVENT_KEY
				case CLASS_FILTER:
					metadate = filter.FilterEventMetadata{
						AppKey: key.appKey,
						Env:    key.env,
					}
					eventKey = filter.EVENT_KEY
				default:
					log.Error("key class <" + key.class + ">is not declare")
					continue
				}
				err = consumers.AddEvent(&event.Event{
					EventDesc: event.EventDesc{
						EventType: event.Event_Delete,
						EventKey:  eventKey,
					},
					Metadata: metadate,
				})
				if err != nil {
					log.Error(err)
				}
				log.Info("etcd delete key:", string(v.Key))
			}
		}
	}
}

// GetConfigVal ...
func (f *FileStore) GetConfigVal(appKey mconfig.AppKey, configKey mconfig.ConfigKey, env mconfig.ConfigEnv) (*mconfig.StoreVal, error) {
	entity := &KeyEntity{
		namespace: namespce,
		class:     CLASS_CONFIG,
		appKey:    appKey,
		configKey: configKey,
		env:       env,
	}
	key, err := getStoreKey(entity)
	if err != nil {
		return nil, err
	}
	v, err := f.config.Get([]byte(key), nil)
	if err != nil {
		return nil, err
	}
	val := &mconfig.StoreVal{}
	err = json.Unmarshal(v, val)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	return val, nil
}

// GetFilterVal ...
func (f *FileStore) GetFilterVal(appKey mconfig.AppKey, env mconfig.ConfigEnv) (*mconfig.StoreVal, error) {
	entity := &KeyEntity{
		namespace: namespce,
		class:     CLASS_FILTER,
		appKey:    appKey,
		env:       env,
	}
	key, err := getStoreKey(entity)
	if err != nil {
		return nil, err
	}
	v, err := f.filter.Get([]byte(key), nil)
	val := &mconfig.StoreVal{}
	err = json.Unmarshal(v, val)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	return val, nil
}

// Close ...
func (e *FileStore) Close() error {
	e.cancelFunc()
	return nil
}
