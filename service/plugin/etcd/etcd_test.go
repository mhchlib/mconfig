package etcd

import (
	"context"
	"github.com/mhchlib/mconfig/service"
	"log"
	"testing"
	"time"
)

func TestEtcdStore_GetConfig(t *testing.T) {
	e := &EtcdStore{}
	config, rev, err := e.GetConfig("1000")
	log.Println(config, rev, err)
}

func TestEtcdStore_PutConfig(t *testing.T) {
	e := &EtcdStore{}
	err := e.PutConfig("1000", "{'aaa':'ccc'}")
	log.Println(err)
}

func TestEtcdStore_WatchConfig(t *testing.T) {
	e := &EtcdStore{}
	_, rev, err := e.GetConfig("1000")
	log.Println("get mconfig rev: ", rev)
	if err != nil {
		log.Fatal(err)
	}
	ctx, cancelFunc := context.WithCancel(context.Background())
	configChan, err := e.WatchConfig("1000", rev, ctx)
	go func(configChan chan *service.ConfigEvent) {
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

func TestEtcdStore_WatchConfigWithPrefix(t *testing.T) {
	e := &EtcdStore{}

	ctx, cancelFunc := context.WithCancel(context.Background())
	configChan, _ := e.WatchConfigWithPrefix(ctx)
	go func(configChan chan *service.ConfigEvent) {
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
