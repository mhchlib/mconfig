package etcd

import (
	"context"
	"github.com/coreos/etcd/clientv3"
	log "github.com/mhchlib/logger"
	"google.golang.org/grpc"
	"testing"
	"time"
)

func TestEtcdStore_GetSyncData(t *testing.T) {
	var err error
	cli, err = clientv3.New(clientv3.Config{
		Endpoints:   []string{"etcd.u.hcyang.top:31770"},
		DialTimeout: time.Second * 5,
		DialOptions: []grpc.DialOption{grpc.WithBlock()},
	})
	if err != nil {
		log.Fatal(err)
	}
	kv = clientv3.NewKV(cli)
	Response, err := kv.Get(context.Background(), "/com", clientv3.WithPrefix())
	if err != nil {
		log.Fatal(err)
	}
	log.Info(Response.Kvs)
}
