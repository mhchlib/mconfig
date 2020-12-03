package main

import (
	"context"
	log "github.com/mhchlib/logger"
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
			AppKey := strconv.Itoa(1000 + i%5)
			mService := micro.NewService()
			mService.Init()
			mConfigService := sdk.NewMConfigService("com.github.mhchlib.mconfig", mService.Client())
			resp, err := mConfigService.GetVStream(context.Background(), &sdk.GetVRequest{
				AppKey:  appkey,
				Filters: &sdk.ConfigFilters{},
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
				log.Info(appkey, " get msg")
				log.Info(" ------------------- ")
				//log.Info(config.Configs)
			}
			group.Done()
		}(i)
	}
	group.Wait()
}
