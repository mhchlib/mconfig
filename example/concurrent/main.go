package main

import (
	"context"
	log "github.com/mhchlib/logger"
	"github.com/mhchlib/mconfig-api/api/v1/common"
	"github.com/mhchlib/mconfig-api/api/v1/sdk"
	"github.com/micro/go-micro/v2"
	"strconv"
	"sync"
)

func main() {
	count := 100
	group := sync.WaitGroup{}
	group.Add(count)
	for i := 0; i < count; i++ {
		go func(a int) {
			appid := strconv.Itoa(1000 + i%9)
			mService := micro.NewService()
			mService.Init()
			mConfigService := sdk.NewMConfigService("", mService.Client())
			resp, err := mConfigService.GetVStream(context.Background(), &sdk.GetVRequest{
				AppId:    appid,
				ConfigId: []string{"100", "101"},
				ExtraData: []*common.ExtraData{
					{
						Key:   "ip",
						Value: "10.92.12.3",
					},
				},
			})

			if err != nil {
				log.Error(err)
				group.Done()
				return
			}

			defer func() {
				log.Info("close stream")
				if resp != nil {
					_ = resp.Close()
				}
			}()
			//log.Info(a)
			for {
				//config, err := resp.Recv()
				_, err := resp.Recv()
				if err != nil {
					log.Error(err)
					return
				}
				log.Info(appid, " get msg")
				log.Info(" ------------------- ")
				//log.Info(config.Configs)
			}
			group.Done()
		}(i)
	}
	group.Wait()
}
