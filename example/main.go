package main

import (
	"context"
	"github.com/mhchlib/mconfig/api"
	"github.com/micro/go-micro/v2"
	log "github.com/micro/go-micro/v2/util/log"
	"sync"
)

func main() {
	group := sync.WaitGroup{}
	group.Add(1000)
	for i := 0; i < 100; i++ {
		go func(a int) {
			mService := micro.NewService()
			mService.Init()
			mConfigService := api.NewMConfigService("", mService.Client())
			resp, err := mConfigService.GetVStream(context.Background(), &api.GetVRequest{Configid: "1000"})
			if err != nil {
				log.Fatal(err)
			}
			log.Info(a)
			for {
				config, err := resp.Recv()
				if err != nil {
					log.Fatal(err)
				}
				log.Info(config)
			}
			group.Done()
		}(i)
	}
	group.Wait()
}
