package main

import (
	"context"
	log "github.com/mhchlib/logger"
	"github.com/mhchlib/mconfig-api/api/v1/sdk"
	"github.com/micro/go-micro/v2"
	"strconv"
)

func main() {
	//appid := 1000 + rand.Intn(4)
	appId := 1000
	configIds := []string{"1000-100", "1000-103"}
	extreData := map[string]string{
		"ip": "192.168.1.1",
	}
	log.Info("client listen app ", appId, " config ", configIds, " with data ", extreData)
	mService := micro.NewService()
	mService.Init()
	mConfigService := sdk.NewMConfigService("", mService.Client())
	resp, err := mConfigService.GetVStream(context.Background(), &sdk.GetVRequest{
		AppId: strconv.Itoa(appId),
		Filters: &sdk.ConfigFilters{
			ConfigIds: configIds,
			ExtraData: extreData,
		},
	})
	if err != nil {
		log.Fatal(err)
		return
	}
	defer func() {
		log.Info("close stream")
		_ = resp.Close()
	}()
	//log.Info(a)
	for {
		config, err := resp.Recv()
		//_, err := resp.Recv()
		if err != nil {
			log.Fatal(err)
			return
		}
		log.Info(appId, " get msg")
		log.Info(" ------------------- ")
		log.Info(config.Configs)
	}
}
