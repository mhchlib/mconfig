package main

import (
	log "github.com/mhchlib/logger"
	"github.com/mhchlib/mconfig-api/api/v1/cli"
	"github.com/mhchlib/mconfig-api/api/v1/sdk"
	"github.com/mhchlib/mconfig/pkg"
	"github.com/micro/go-micro/v2"
	"time"
)

func init() {
}

func main() {
	defer pkg.InitMconfig()()
	mService := micro.NewService(Opt_RegistryTimeout)
	mService.Init()
	initRpc(mService)
	err := mService.Run()
	if err != nil {
		log.Fatal(err)
	}
}

func initRpc(mService micro.Service) {
	err := sdk.RegisterMConfigHandler(mService.Server(), pkg.NewMConfigSDK())
	if err != nil {
		log.Fatal(err)
	}
	err = cli.RegisterMConfigCliHandler(mService.Server(), pkg.NewMConfigCLI())
	if err != nil {
		log.Fatal(err)
	}
}

func Opt_RegistryTimeout(options *micro.Options) {
	resOptions := options.Registry.Options()
	resOptions.Timeout = 30 * time.Second
}
