package etcd

import (
	"context"
	"github.com/mhchlib/mconfig/pkg"
	"github.com/mhchlib/mconfig/pkg/config"
	"github.com/mhchlib/mconfig/pkg/mconfig"
	"log"
	"testing"
	"time"
)

func TestEtcdStore_GetAppConfigs(t *testing.T) {
	e := &EtcdStore{}
	config, err := e.GetAppConfigs(mconfig.Appkey("1000"))
	log.Println(config, err)
}

func TestEtcdStore_PutAppConfigs(t *testing.T) {
	e := &EtcdStore{}
	err := e.PutAppConfigs("1000", &config.AppConfigs{})
	log.Println(err)
}

func TestEtcdStore_WatchAppConfigsWithPrefix(t *testing.T) {
	e := &EtcdStore{}

	ctx, cancelFunc := context.WithCancel(context.Background())
	configChan, _ := e.WatchAppConfigs(ctx)
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
