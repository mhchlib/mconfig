package etcd

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/mvcc/mvccpb"
	log "github.com/mhchlib/logger"
	"github.com/mhchlib/mconfig/pkg"
	"google.golang.org/grpc"
	"strings"
	"time"
)

var (
	cli           clientv3.Client
	kv            clientv3.KV
	watcher       clientv3.Watcher
	etcdAddresses *string
)

// PREFIX_CONFIG ...
const PREFIX_CONFIG = "/mconfig/"

type AppConfigsJSONStr string

// EtcdStore ...
type EtcdStore struct {
}

func init() {
	pkg.RegisterStorePlugin("etcd", Init)
	initFlag()
}

func initFlag() {
	etcdAddresses = flag.String("store_etcd", "127.0.0.1:2389", "input your etcd address, multiple ip commas separate")
}

// Init ...
func Init() (pkg.AppConfigStore, error) {
	address := strings.Split(*etcdAddresses, ",")
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   address,
		DialTimeout: time.Second * 5,
		DialOptions: []grpc.DialOption{grpc.WithBlock()},
	})
	if err != nil {
		log.Fatal("dial to store etcd err :", err, "addr: ", *etcdAddresses)
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

// GetAppConfigs ...
func (e EtcdStore) GetAppConfigs(key pkg.Appkey) (*pkg.AppConfigs, error) {
	get, err := kv.Get(context.TODO(), Prefix(PREFIX_CONFIG, string(key)))
	if err != nil {
		log.Error(err)
	}
	if get.Count == 1 {
		appConfigs, err := parseAppConfigsJSONStr(AppConfigsJSONStr(string(get.Kvs[0].Value)))
		if err != nil {
			return nil, err
		}
		return appConfigs, nil
	} else {
		return nil, errors.New(string("app id: " + key + " not found"))
	}
}

// PutAppConfigs ...
func (e EtcdStore) PutAppConfigs(key pkg.Appkey, value *pkg.AppConfigs) error {
	configJsonStr, err := json.Marshal(value)
	if err != nil {
		return err
	}
	_, err = kv.Put(context.TODO(), string(PREFIX_CONFIG+key), string(configJsonStr))
	if err != nil {
		return err
	}
	return nil
}

// WatchAppConfigs ...
func (e EtcdStore) WatchAppConfigs(ctx context.Context) (chan *pkg.ConfigEvent, error) {
	var watchChan clientv3.WatchChan
	watchChan = watcher.Watch(ctx, Prefix(PREFIX_CONFIG, ""), clientv3.WithPrefix())
	configChan := make(chan *pkg.ConfigEvent)
	go func(ctx context.Context, watchChan <-chan clientv3.WatchResponse, configChan chan<- *pkg.ConfigEvent) {
		defer func() {
			close(configChan)
		}()
		for {
			select {
			case v, err := <-watchChan:
				if err == false {
					log.Error("watcher err ...")
				}
				if v.Canceled {
					log.Error("watcher err ..." + v.Err().Error())
				}
				events := v.Events
				for _, event := range events {
					//log.Info("get event value : ", string(event.Kv.Value))
					switch event.Type {
					case mvccpb.PUT:
						appConfigs, err := parseAppConfigsJSONStr((AppConfigsJSONStr)(event.Kv.Value))
						if err != nil {
							log.Error("app key: ", PREFIX_CONFIG, " mvccpb.PUT ", err)
						}
						configChan <- &pkg.ConfigEvent{
							Key:        (pkg.Appkey)(RemovePrefix(PREFIX_CONFIG, string(event.Kv.Key))),
							AppConfigs: appConfigs,
							EventType:  pkg.Event_Update,
						}
					case mvccpb.DELETE:
						configChan <- &pkg.ConfigEvent{
							Key:       (pkg.Appkey)(event.Kv.Key),
							EventType: pkg.Event_Delete,
						}
					}
				}
			case <-ctx.Done():
				log.Info("watcher done ...")
				return
			}
		}
	}(ctx, watchChan, configChan)
	return configChan, nil
}

// Prefix ...
func Prefix(prefix string, v string) string {
	return prefix + v
}

// RemovePrefix ...
func RemovePrefix(prefix string, v string) string {
	return strings.ReplaceAll(v, prefix, "")
}
