package main

import (
	"context"
	log "github.com/mhchlib/logger"
	"github.com/mhchlib/mconfig-api/api/v1/common"
	"github.com/mhchlib/mconfig-api/api/v1/sdk"
	"github.com/micro/go-micro/v2"
	"sync"
)

func main() {
	group := sync.WaitGroup{}
	group.Add(1000)
	for i := 0; i < 100; i++ {
		go func(a int) {
			mService := micro.NewService()
			mService.Init()
			mConfigService := sdk.NewMConfigService("", mService.Client())
			resp, err := mConfigService.GetVStream(context.Background(), &sdk.GetVRequest{
				AppId:    "1000",
				ConfigId: []string{"100", "101"},
				ExtraData: []*common.ExtraData{
					{
						Key:   "ip",
						Value: "10.92.12.3",
					},
				},
			})
			if err != nil {
				log.Fatal(err)
				return
			}
			//log.Info(a)
			for {
				config, err := resp.Recv()
				if err != nil {
					log.Fatal(err)
					return
				}
				log.Println(config.Configs)
			}
			group.Done()
		}(i)
	}
	group.Wait()
}
