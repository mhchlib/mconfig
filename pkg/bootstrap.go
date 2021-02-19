package pkg

import (
	"github.com/mhchlib/mconfig/pkg/client"
	"github.com/mhchlib/mconfig/pkg/config"
	"github.com/mhchlib/mconfig/pkg/event"
	"github.com/mhchlib/mconfig/pkg/filter"
	"github.com/mhchlib/mconfig/pkg/mconfig"
	"github.com/mhchlib/mconfig/pkg/store"
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
