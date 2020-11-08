package main

import (
	log "github.com/mhchlib/logger"
	"github.com/mhchlib/mconfig-api/api/v1/sdk"
	"github.com/mhchlib/mconfig/pkg"
	"github.com/micro/go-micro/v2"
)

func init() {
}

func main() {
	defer pkg.InitMconfig()()
	mService := micro.NewService()
	mService.Init()
	err := sdk.RegisterMConfigHandler(mService.Server(), &pkg.MConfig{})
	if err != nil {
		log.Fatal(err)
	}
	err = mService.Run()
	if err != nil {
		log.Fatal(err)
	}
}
