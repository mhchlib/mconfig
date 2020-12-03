package etcd

import (
	"context"
	"errors"
	"github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/mvcc/mvccpb"
	log "github.com/mhchlib/logger"
	"github.com/mhchlib/mconfig/pkg"
	"google.golang.org/grpc"
	"strings"
	"time"
)

var (
	cli     clientv3.Client
	kv      clientv3.KV
	watcher clientv3.Watcher
)

const PREFIX_CONFIG = "/mconfig/"

type EtcdStore struct {
}

func init() {
	pkg.RegisterStorePlugin("etcd", Init)
}

func Init(addr string) (pkg.AppConfigStore, error) {
	address := strings.Split(addr, ",")
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   address,
		DialTimeout: time.Second * 5,
		DialOptions: []grpc.DialOption{grpc.WithBlock()},
	})
	if err != nil {
		log.Fatal("dial to store etcd err :", err, "addr: ", addr)
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

func (e EtcdStore) GetAppConfigs(key string) (pkg.AppConfigsJSONStr, int64, error) {
	get, err := kv.Get(context.TODO(), Prefix(PREFIX_CONFIG, key))
	if err != nil {
		log.Error(err)
	}
	if get.Count == 1 {
		return pkg.AppConfigsJSONStr(string(get.Kvs[0].Value)), get.Header.Revision, nil
	} else {
		return "", 0, errors.New("app id: " + key + " not found")
	}
}

func (e EtcdStore) PutAppConfigs(key string, value pkg.AppConfigsJSONStr) error {
	_, err := kv.Put(context.TODO(), PREFIX_CONFIG+key, string(value))
	if err != nil {
		return err
	}
	return nil
}

func (e EtcdStore) WatchAppConfigs(key string, rev int64, ctx context.Context) (chan *pkg.ConfigEvent, error) {
	var watchChan clientv3.WatchChan
	if rev != 0 {
		watchChan = watcher.Watch(ctx, Prefix(PREFIX_CONFIG, key), clientv3.WithRev(rev))
	} else {
		watchChan = watcher.Watch(ctx, Prefix(PREFIX_CONFIG, key))
	}
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
						configChan <- &pkg.ConfigEvent{
							Key:   (pkg.Appkey)(RemovePrefix(PREFIX_CONFIG, string(event.Kv.Key))),
							Value: (pkg.AppConfigsJSONStr)(event.Kv.Value),
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

func (e EtcdStore) WatchAppConfigsWithPrefix(ctx context.Context) (chan *pkg.ConfigEvent, error) {
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
						configChan <- &pkg.ConfigEvent{
							Key:       (pkg.Appkey)(RemovePrefix(PREFIX_CONFIG, string(event.Kv.Key))),
							Value:     (pkg.AppConfigsJSONStr)(event.Kv.Value),
							EventType: pkg.Event_Update,
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

func Prefix(prefix string, v string) string {
	return prefix + v
}

func RemovePrefix(prefix string, v string) string {
	return strings.ReplaceAll(v, prefix, "")
}
