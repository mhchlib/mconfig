package etcd

import (
	"context"
	"github.com/mhchlib/mconfig/pkg"
	"log"
	"testing"
	"time"
)

func TestEtcdStore_GetAppConfigs(t *testing.T) {
	e := &EtcdStore{}
	config, rev, err := e.GetAppConfigs("1000")
	log.Println(config, rev, err)
}

func TestEtcdStore_PutAppConfigs(t *testing.T) {
	e := &EtcdStore{}
	err := e.PutAppConfigs("1000", "{'aaa':'ccc'}")
	log.Println(err)
}

func TestEtcdStore_WatchAppConfigs(t *testing.T) {
	e := &EtcdStore{}
	_, rev, err := e.GetAppConfigs("1000")
	log.Println("get mconfig rev: ", rev)
	if err != nil {
		log.Fatal(err)
	}
	ctx, cancelFunc := context.WithCancel(context.Background())
	configChan, err := e.WatchAppConfigs("1000", rev, ctx)
	go func(configChan chan *pkg.ConfigEvent) {
		for {
			select {
			case v, ok := <-configChan:
				if !ok {
					return
				}
				log.Println(v)
			}

		}
	}(configChan)
	time.Sleep(time.Second * 10)
	cancelFunc()
	time.Sleep(time.Second * 3)
	log.Println("over...")
}

func TestEtcdStore_WatchAppConfigsWithPrefix(t *testing.T) {
	e := &EtcdStore{}

	ctx, cancelFunc := context.WithCancel(context.Background())
	configChan, _ := e.WatchAppConfigsWithPrefix(ctx)
	go func(configChan chan *pkg.ConfigEvent) {
		for {
			select {
			case v, ok := <-configChan:
				if !ok {
					return
				}
				log.Println(v)
			}

		}
	}(configChan)
	time.Sleep(time.Second * 60)
	cancelFunc()
	time.Sleep(time.Second * 3)
	log.Println("over...")
}

func TestRemovePrefix(t *testing.T) {
	new := RemovePrefix("config/", "config/1111")
	log.Println(new)
	log.Println(new == "1111")
}
