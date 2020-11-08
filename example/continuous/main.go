package main

import (
	"context"
	log "github.com/mhchlib/logger"
	"github.com/mhchlib/mconfig-api/api/v1/common"
	"github.com/mhchlib/mconfig-api/api/v1/sdk"
	"github.com/micro/go-micro/v2"
	"math/rand"
	"strconv"
)

func main() {
	appid := 1000 + rand.Intn(10)
	mService := micro.NewService()
	mService.Init()
	mConfigService := sdk.NewMConfigService("", mService.Client())
	resp, err := mConfigService.GetVStream(context.Background(), &sdk.GetVRequest{
		AppId:    strconv.Itoa(appid),
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
	defer func() {
		log.Info("close stream")
		_ = resp.Close()
	}()
	//log.Info(a)
	for {
		//config, err := resp.Recv()
		_, err := resp.Recv()
		if err != nil {
			log.Fatal(err)
			return
		}
		log.Info(appid, " get msg")
		log.Info(" ------------------- ")
		//log.Info(config.Configs)
	}
}
