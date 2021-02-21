package core

import (
	"github.com/mhchlib/mconfig/core/client"
	"github.com/mhchlib/mconfig/core/config"
	"github.com/mhchlib/mconfig/core/event"
	"github.com/mhchlib/mconfig/core/filter"
	"github.com/mhchlib/mconfig/core/mconfig"
	"github.com/mhchlib/mconfig/core/store"
)

// InitMconfig ...
func InitMconfig(mconfig *mconfig.MConfig) func() {
	config.InitConfigCenter()
	filter.InitFilterEngine()
	client.InitClientManagement()
	store.InitStore(mconfig.StoreType, mconfig.StoreAddress)
	go event.InitEventBus()
	return EndMconfig()
}

// EndMconfig ...
func EndMconfig() func() {
	return func() {
	}
}
