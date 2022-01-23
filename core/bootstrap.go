package core

import (
	log "github.com/mhchlib/logger"
	"github.com/mhchlib/mconfig/core/client"
	"github.com/mhchlib/mconfig/core/config"
	"github.com/mhchlib/mconfig/core/cron"
	"github.com/mhchlib/mconfig/core/event"
	"github.com/mhchlib/mconfig/core/filter"
	"github.com/mhchlib/mconfig/core/mconfig"
	"github.com/mhchlib/mconfig/core/store"
)

// InitMconfig ...
func InitMconfig(mconfig *mconfig.MConfigConfig) func() {
	event.InitEventBus()
	config.InitConfigCenter()
	filter.InitFilterEngine()
	client.InitClientManagement()
	storeType, storeGracefulStopFunc, err := store.InitStore(mconfig.StoreAddress)
	if err != nil {
		log.Fatal(err)
	}
	mconfig.StoreType = storeType
	cron.InitCron()
	log.Info("mconfig core init success")
	return func() {
		storeGracefulStopFunc()
	}
}
