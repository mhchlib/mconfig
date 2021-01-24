package etcd

import (
	"context"
	"errors"
	"github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/mvcc/mvccpb"
	log "github.com/mhchlib/logger"
	"github.com/mhchlib/mconfig/pkg/config"
	"github.com/mhchlib/mconfig/pkg/event"
	"github.com/mhchlib/mconfig/pkg/mconfig"
	"github.com/mhchlib/mconfig/pkg/store"
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
	PLUGIN_NAME          = "etcd"
	SEPARATOR            = "/"
	MODE_WATCH   KeyMode = "watch"
	MODE_DEFAULT KeyMode = "default"

	CLASS_CONFIG  KeyClass = "config"
	CLASS_FILTER  KeyClass = "filter"
	CLASS_VERSION KeyClass = "version"
	CLASS_META    KeyClass = "metadata"
)

var namespce KeyNamespce = "com.github.hchlib.mconfig"

var prefix_mode_watch = SEPARATOR + string(namespce) + SEPARATOR + string(MODE_WATCH)
var prefix_mode_default = SEPARATOR + string(namespce) + SEPARATOR + string(MODE_DEFAULT)

type KeyEntity struct {
	namespace KeyNamespce
	mode      KeyMode
	class     KeyClass
	appKey    mconfig.Appkey
	configKey mconfig.ConfigKey
	env       mconfig.ConfigEnv
}

type EtcdStore struct {
	cancelFunc context.CancelFunc
}

func (e *EtcdStore) WatchDynamicVal(consumers *store.Consumer) error {
	var watchChan clientv3.WatchChan
	ctx, cancelFunc := context.WithCancel(context.Background())
	e.cancelFunc = cancelFunc
	watchChan = watcher.Watch(ctx, Prefix(prefix_mode_watch, ""), clientv3.WithPrefix())
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
					key, err := parseEventKey(string(e.Kv.Key))
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
							Val:       mconfig.ConfigVal(e.Kv.Value),
						}
						eventKey = config.EVENT_KEY
					case CLASS_FILTER:
					//metadate := config.ConfigEventMetadata{
					//	AppKey:    key.appKey,
					//	ConfigKey: key.configKey,
					//	Env:       key.env,
					//	Val:       mconfig.ConfigVal(e.Kv.Value),
					//}
					//eventKey := config.EVENT_KEY
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
					log.Debug("etcd update key:", string(e.Kv.Key), "value:", string(e.Kv.Value))
				case mvccpb.DELETE:
					key, err := parseEventKey(string(e.Kv.Key))
					if err != nil {
						log.Error(err)
						break
					}
					err = consumers.AddEvent(&event.Event{
						EventDesc: event.EventDesc{
							EventType: event.Event_Delete,
							EventKey:  config.EVENT_KEY,
						},
						Metadata: config.ConfigEventMetadata{
							AppKey:    key.appKey,
							ConfigKey: key.configKey,
							Env:       key.env,
						},
					})
					if err != nil {
						log.Error(err)
					}
				}
			}
		}
	}
}

func (e *EtcdStore) GetConfigVal(appKey mconfig.Appkey, configKey mconfig.ConfigKey, env mconfig.ConfigEnv) (mconfig.ConfigVal, error) {
	entity := &KeyEntity{
		namespace: namespce,
		mode:      MODE_WATCH,
		class:     CLASS_CONFIG,
		appKey:    appKey,
		configKey: configKey,
		env:       env,
	}
	key, err := getEventKey(entity)
	if err != nil {
		return "", err
	}
	response, err := cli.Get(context.Background(), key)
	if err != nil {
		return "", err
	}
	if response.Count != 1 {
		return "", errors.New("not found")
	}
	return mconfig.ConfigVal(response.Kvs[0].Value), nil
}

func (e *EtcdStore) PutConfigVal(appKey mconfig.Appkey, configKey mconfig.ConfigKey, env mconfig.ConfigEnv, content mconfig.ConfigVal) error {

	return nil
}

func (e *EtcdStore) NewAppMetaData(meta *mconfig.AppMetaData) error {
	panic("implement me")
}

func (e *EtcdStore) NewConfigMetaData(meta *mconfig.ConfigMetaData) error {
	panic("implement me")
}

func (e *EtcdStore) ListAppConfigsMeta(limit int, offset int, filter string, appKey mconfig.Appkey) ([]*mconfig.ConfigMetaData, error) {
	panic("implement me")
}

func (e *EtcdStore) UpdateAppMetaData(meta *mconfig.AppMetaData) error {
	panic("implement me")
}

func (e *EtcdStore) UpdateConfigMetaData(meta *mconfig.ConfigMetaData) error {
	panic("implement me")
}

func (e *EtcdStore) DeleteApp(appKey mconfig.Appkey) error {
	panic("implement me")
}

func (e *EtcdStore) DeleteConfig(appKey mconfig.Appkey, configKey mconfig.ConfigKey) error {
	panic("implement me")
}

func (e *EtcdStore) ListAppMetaData(limit int, offset int, filter string) error {
	panic("implement me")
}

func (e *EtcdStore) Close() error {
	e.cancelFunc()
	return nil
}
