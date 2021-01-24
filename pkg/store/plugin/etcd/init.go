package etcd

import (
	"context"
	"github.com/coreos/etcd/clientv3"
	log "github.com/mhchlib/logger"
	"github.com/mhchlib/mconfig/pkg/store"
	"google.golang.org/grpc"
	"strings"
	"time"
)

func init() {
	store.RegisterStorePlugin(PLUGIN_NAME, Init)
}

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
