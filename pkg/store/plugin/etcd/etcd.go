package etcd

import (
	"context"
	"github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/mvcc/mvccpb"
	log "github.com/mhchlib/logger"
	"github.com/mhchlib/mconfig/pkg/config"
	"github.com/mhchlib/mconfig/pkg/event"
	"github.com/mhchlib/mconfig/pkg/mconfig"
	"github.com/mhchlib/mconfig/pkg/store"
	"google.golang.org/grpc"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

var (
	cli     clientv3.Client
	kv      clientv3.KV
	watcher clientv3.Watcher
)

// PREFIX_CONFIG ...
const PREFIX_CONFIG = "/mconfig/"

// EtcdStore ...
type EtcdStore struct {
	cancelFunc context.CancelFunc
}

func (e *EtcdStore) WatchConfigVal(consumers *store.Consumer) error {
	var watchChan clientv3.WatchChan
	ctx, cancelFunc := context.WithCancel(context.Background())
	e.cancelFunc = cancelFunc
	watchChan = watcher.Watch(ctx, Prefix(PREFIX_CONFIG, ""), clientv3.WithPrefix())
	for {
		select {
		case v, ok := <-watchChan:
			if ok == false {
				log.Error("watcher err ...")
				return store.Error_WatchFail
			}
			if v.Canceled {
				log.Error("watcher err ..." + v.Err().Error())
				return store.Error_WatchFail
			}
			events := v.Events
			for _, e := range events {
				//log.Info("get event value : ", string(event.Kv.Value))
				switch e.Type {
				case mvccpb.PUT:
					err := consumers.AddEvent(&event.Event{
						EventDesc: event.EventDesc{
							EventType: event.Event_Update,
							EventKey:  config.EVENT_KEY,
						},
						Metadata: config.ConfigEventMetadata{
							AppKey:    "appKey",
							ConfigKey: "configKey",
							Env:       "dev",
							Val:       mconfig.ConfigVal("66666" + strconv.Itoa(rand.Intn(1000))),
						},
					})
					if err != nil {
						log.Error(err)
					}
					log.Info("etcd add event to consumers")

					//appConfigs, err := parseAppConfigsJSONStr((AppConfigsJSONStr)(event.Kv.Value))
					//if err != nil {
					//	log.Error("app key: ", PREFIX_CONFIG, " mvccpb.PUT ", err)
					//}
					//configChan <- &pkg.ConfigEvent{
					//	Key:        (mconfig.Appkey)(RemovePrefix(PREFIX_CONFIG, string(event.Kv.Key))),
					//	AppConfigs: appConfigs,
					//	EventType:  pkg.Event_Update,
					//}
				case mvccpb.DELETE:

				}
			}
		}
	}
	return nil
}

func (e EtcdStore) GetConfigVal(appKey mconfig.Appkey, configKey mconfig.ConfigKey, env mconfig.ConfigEnv) (mconfig.ConfigVal, error) {
	panic("implement me")
}

func (e EtcdStore) PutConfigVal(appKey mconfig.Appkey, configKey mconfig.ConfigKey, env mconfig.ConfigEnv, content mconfig.ConfigVal) error {
	panic("implement me")
}

func (e EtcdStore) NewAppMetaData(meta mconfig.AppMetaData) error {
	panic("implement me")
}

func (e EtcdStore) NewConfigMetaData(meta mconfig.ConfigMetaData) error {
	panic("implement me")
}

func (e EtcdStore) GetAppConfigs(appKey mconfig.Appkey) ([]mconfig.ConfigMetaData, error) {
	panic("implement me")
}

func (e EtcdStore) UpdateAppMetaData(meta mconfig.AppMetaData) error {
	panic("implement me")
}

func (e EtcdStore) UpdateConfigMetaData(meta mconfig.ConfigMetaData) error {
	panic("implement me")
}

func (e EtcdStore) DeleteApp(appKey mconfig.Appkey) error {
	panic("implement me")
}

func (e EtcdStore) DeleteConfig(appKey mconfig.Appkey, configKey mconfig.ConfigKey) error {
	panic("implement me")
}

func (e EtcdStore) ListAppMetaData(limit int, offset int, filter string) error {
	panic("implement me")
}

func (e EtcdStore) Close() error {
	panic("implement me")
}

func init() {
	store.RegisterStorePlugin("etcd", Init)
}

// Init ...
func Init(addressStr string) (store.MConfigStore, error) {
	address := strings.Split(addressStr, ",")
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   address,
		DialTimeout: time.Second * 5,
		DialOptions: []grpc.DialOption{grpc.WithBlock()},
	})
	if err != nil {
		log.Fatal("dial to store etcd err :", err, "addr: ", addressStr)
	}
	kv = clientv3.NewKV(cli)
	watcher = clientv3.NewWatcher(cli)
	var list *clientv3.MemberListResponse
	timeoutCtx, _ := context.WithTimeout(context.Background(), time.Second*5)
	list, err = cli.MemberList(timeoutCtx)
	if err != nil {
		log.Fatal("etcd member list error :", err)
	}
	log.Info("etcd member list : ", list.Members)
	return &EtcdStore{}, nil
}

//// GetAppConfigs ...
//func (e EtcdStore) GetAppConfigs(key mconfig.Appkey) (*config.AppConfigs, error) {
//	get, err := kv.Get(context.TODO(), Prefix(PREFIX_CONFIG, string(key)))
//	if err != nil {
//		log.Error(err)
//	}
//	if get.Count == 1 {
//		appConfigs, err := parseAppConfigsJSONStr(AppConfigsJSONStr(string(get.Kvs[0].Value)))
//		if err != nil {
//			return nil, err
//		}
//		return appConfigs, nil
//	} else {
//		return nil, errors.New(string("app id: " + key + " not found"))
//	}
//}
//
//// PutAppConfigs ...
//func (e EtcdStore) PutAppConfigs(key mconfig.Appkey, value *config.AppConfigs) error {
//	configJsonStr, err := json.Marshal(value)
//	if err != nil {
//		return err
//	}
//	_, err = kv.Put(context.TODO(), string(PREFIX_CONFIG+key), string(configJsonStr))
//	if err != nil {
//		return err
//	}
//	return nil
//}
//

// Prefix ...
func Prefix(prefix string, v string) string {
	return prefix + v
}

// RemovePrefix ...
func RemovePrefix(prefix string, v string) string {
	return strings.ReplaceAll(v, prefix, "")
}
