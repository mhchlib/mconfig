package etcd

import (
	"context"
	"errors"
	"github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/mvcc/mvccpb"
	log "github.com/mhchlib/logger"
	"github.com/mhchlib/mconfig/service"
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
	//Init etcd link
	initEtcd()
	//init store
	service.ConfigStore = &EtcdStore{}
}

func (e EtcdStore) GetConfig(key string) (service.ConfigJSONStr, int64, error) {
	get, err := kv.Get(context.TODO(), Prefix(PREFIX_CONFIG, key))
	if err != nil {
		log.Error(err)
	}
	if get.Count == 1 {
		return service.ConfigJSONStr(string(get.Kvs[0].Value)), get.Header.Revision, nil
	} else {
		return "", 0, errors.New("configid: " + key + " no value")
	}
}

func (e EtcdStore) PutConfig(key string, value service.ConfigJSONStr) error {
	//TODO:version control
	_, err := kv.Put(context.TODO(), PREFIX_CONFIG+key, string(value))
	if err != nil {
		return err
	}
	return nil
}

func (e EtcdStore) WatchConfig(key string, rev int64, ctx context.Context) (chan *service.ConfigEvent, error) {
	var watchChan clientv3.WatchChan
	if rev != 0 {
		watchChan = watcher.Watch(ctx, Prefix(PREFIX_CONFIG, key), clientv3.WithRev(rev))
	} else {
		watchChan = watcher.Watch(ctx, Prefix(PREFIX_CONFIG, key))
	}
	configChan := make(chan *service.ConfigEvent)
	go func(ctx context.Context, watchChan <-chan clientv3.WatchResponse, configChan chan<- *service.ConfigEvent) {
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
					log.Println("get event value : ", string(event.Kv.Value))
					switch event.Type {
					case mvccpb.PUT:
						configChan <- &service.ConfigEvent{
							Key:   (service.AppId)(RemovePrefix(PREFIX_CONFIG, string(event.Kv.Key))),
							Value: (service.ConfigJSONStr)(event.Kv.Value),
						}
					case mvccpb.DELETE:
						configChan <- &service.ConfigEvent{
							Key:       (service.AppId)(event.Kv.Key),
							EventType: service.Event_Delete,
						}
					}
				}
			case <-ctx.Done():
				log.Println("watcher done ...")
				return
			}
		}
	}(ctx, watchChan, configChan)
	return configChan, nil
}

func (e EtcdStore) WatchConfigWithPrefix(ctx context.Context) (chan *service.ConfigEvent, error) {
	var watchChan clientv3.WatchChan
	watchChan = watcher.Watch(ctx, Prefix(PREFIX_CONFIG, ""), clientv3.WithPrefix())
	configChan := make(chan *service.ConfigEvent)
	go func(ctx context.Context, watchChan <-chan clientv3.WatchResponse, configChan chan<- *service.ConfigEvent) {
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
					log.Println("get event value : ", string(event.Kv.Value))
					switch event.Type {
					case mvccpb.PUT:
						configChan <- &service.ConfigEvent{
							Key:       (service.AppId)(RemovePrefix(PREFIX_CONFIG, string(event.Kv.Key))),
							Value:     (service.ConfigJSONStr)(event.Kv.Value),
							EventType: service.Event_Update,
						}
					case mvccpb.DELETE:
						configChan <- &service.ConfigEvent{
							Key:       (service.AppId)(event.Kv.Key),
							EventType: service.Event_Delete,
						}
					}
				}
			case <-ctx.Done():
				log.Println("watcher done ...")
				return
			}
		}
	}(ctx, watchChan, configChan)
	return configChan, nil
}

func initEtcd() {
	cli, err := clientv3.New(clientv3.Config{
		//TODO：后面需要从flag中获取
		Endpoints: []string{"127.0.0.1:2379"},
		// Endpoints: []string{"localhost:2379", "localhost:22379", "localhost:32379"}
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		log.Fatal(err)
	}
	kv = clientv3.NewKV(cli)
	watcher = clientv3.NewWatcher(cli)
}

func Prefix(prefix string, v string) string {
	return prefix + v
}

func RemovePrefix(prefix string, v string) string {
	return strings.ReplaceAll(v, prefix, "")
}
