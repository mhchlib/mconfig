package etcd

import (
	"context"
	"encoding/json"
	"github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/mvcc/mvccpb"
	log "github.com/mhchlib/logger"
	"github.com/mhchlib/mconfig/core/config"
	"github.com/mhchlib/mconfig/core/event"
	"github.com/mhchlib/mconfig/core/filter"
	"github.com/mhchlib/mconfig/core/mconfig"
	"github.com/mhchlib/mconfig/core/store"
)

var (
	cli     *clientv3.Client
	kv      clientv3.KV
	watcher clientv3.Watcher
)

type KeyNamespce string
type KeyMode string
type KeyClass string

const (
	PLUGIN_NAME           = "etcd"
	SEPARATOR             = "/"
	CLASS_CONFIG KeyClass = "config"
	CLASS_FILTER KeyClass = "filter"
)

var namespce KeyNamespce = "com.github.hchlib.mconfig"

var prefix_common = SEPARATOR + string(namespce)

var prefix_config = prefix_common + SEPARATOR + string(CLASS_CONFIG)
var prefix_filter = SEPARATOR + string(namespce) + SEPARATOR + string(CLASS_FILTER)

type KeyEntity struct {
	namespace KeyNamespce
	class     KeyClass
	appKey    mconfig.AppKey
	configKey mconfig.ConfigKey
	env       mconfig.ConfigEnv
}

type EtcdStore struct {
	cancelFunc context.CancelFunc
}

func (e *EtcdStore) PutConfigVal(appKey mconfig.AppKey, env mconfig.ConfigEnv, configKey mconfig.ConfigKey, val mconfig.StoreVal) error {
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
	_, err = kv.Put(context.Background(), key, string(data))
	return err
}

func (e *EtcdStore) PutFilterVal(appKey mconfig.AppKey, env mconfig.ConfigEnv, val mconfig.StoreVal) error {
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
	_, err = kv.Put(context.Background(), key, string(data))
	return err
}

func (e *EtcdStore) DeleteConfig(appKey mconfig.AppKey, configKey mconfig.ConfigKey, env mconfig.ConfigEnv) error {
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
	_, err = kv.Delete(context.Background(), storeKey)
	if err != nil {
		return err
	}
	return nil
}

func (e *EtcdStore) DeleteFilter(appKey mconfig.AppKey, env mconfig.ConfigEnv) error {
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
	_, err = kv.Delete(context.Background(), storeKey)
	if err != nil {
		return err
	}
	return nil
}

func (e *EtcdStore) GetAppFilters(appKey mconfig.AppKey) ([]*mconfig.StoreVal, error) {
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
	response, err := kv.Get(context.Background(), storeKey, clientv3.WithPrefix())
	if err != nil {
		return nil, err
	}
	for _, kv := range response.Kvs {
		k := string(kv.Key)
		v := kv.Value
		//key, err := parseStoreKey(k)
		//if err != nil {
		//	log.Error(err)
		//	return nil, err
		//}
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

func (e *EtcdStore) GetAppConfigs(appKey mconfig.AppKey, env mconfig.ConfigEnv) ([]*mconfig.StoreVal, error) {
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
	response, err := kv.Get(context.Background(), storeKey, clientv3.WithPrefix())
	if err != nil {
		return nil, err
	}
	for _, kv := range response.Kvs {
		k := string(kv.Key)
		v := kv.Value
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

func (e *EtcdStore) GetSyncData() (mconfig.AppData, error) {
	syncData := make(map[mconfig.AppKey]map[mconfig.ConfigEnv]*mconfig.EnvData)
	Response, err := kv.Get(context.Background(), prefix_common, clientv3.WithPrefix())
	if err != nil {
		return nil, err
	}
	for _, v := range Response.Kvs {
		key := v.Key
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
			err = json.Unmarshal(v.Value, val)
			if err != nil {
				log.Error(err, v.Value)
				continue
			}
			configs[storeKey.configKey] = *val
		}
		if storeKey.class == CLASS_FILTER {
			val := &mconfig.StoreVal{}
			err = json.Unmarshal(v.Value, val)
			if err != nil {
				log.Error(err, v.Value)
				continue
			}
			envData.Filter = *val
		}
	}
	return syncData, nil
}

func (e *EtcdStore) PutSyncData(data *mconfig.AppData) error {
	d, _ := json.Marshal(data)
	log.Info(string(d))
	return nil
}

func (e *EtcdStore) WatchDynamicVal(consumers *store.Consumer) error {
	var watchChan clientv3.WatchChan
	ctx, cancelFunc := context.WithCancel(context.Background())
	e.cancelFunc = cancelFunc
	watchChan = watcher.Watch(ctx, Prefix(prefix_common, ""), clientv3.WithPrefix())
	for {
		select {
		case v, ok := <-watchChan:
			if ok == false {
				return store.Error_FAIL_WATCH
			}
			if v.Canceled {
				log.Error("watcher err ..." + v.Err().Error())
				return store.Error_FAIL_WATCH
			}
			events := v.Events
			for _, e := range events {
				switch e.Type {
				case mvccpb.PUT:
					key, err := parseStoreKey(string(e.Kv.Key))
					if err != nil {
						log.Error(err)
						continue
					}
					var metadate interface{}
					var eventKey event.EventKey

					val := &mconfig.StoreVal{}
					err = json.Unmarshal(e.Kv.Value, val)
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
					log.Info("etcd update key:", string(e.Kv.Key), "value:", string(e.Kv.Value))
				case mvccpb.DELETE:
					key, err := parseStoreKey(string(e.Kv.Key))
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
					log.Info("etcd delete key:", string(e.Kv.Key))
				}
			}
		}
	}
}

func (e *EtcdStore) GetConfigVal(appKey mconfig.AppKey, configKey mconfig.ConfigKey, env mconfig.ConfigEnv) (*mconfig.StoreVal, error) {
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
	Response, err := cli.Get(context.Background(), key)
	if err != nil {
		return nil, err
	}
	if Response.Count != 1 {
		return nil, mconfig.ERROR_STORE_NOT_FOUND
	}
	val := &mconfig.StoreVal{}
	err = json.Unmarshal(Response.Kvs[0].Value, val)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	return val, nil
}

func (e *EtcdStore) GetFilterVal(appKey mconfig.AppKey, env mconfig.ConfigEnv) (*mconfig.StoreVal, error) {
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
	Response, err := cli.Get(context.Background(), key)
	if err != nil {
		return nil, err
	}
	if Response.Count != 1 {
		return nil, mconfig.ERROR_STORE_NOT_FOUND
	}
	val := &mconfig.StoreVal{}
	err = json.Unmarshal(Response.Kvs[0].Value, val)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	return val, nil
}

func (e *EtcdStore) Close() error {
	e.cancelFunc()
	return nil
}
