package pkg

import (
	"context"
	log "github.com/mhchlib/logger"
	"github.com/mhchlib/mconfig/pkg/config"
	"github.com/mhchlib/mconfig/pkg/event"
	"github.com/mhchlib/mconfig/pkg/store"
)

var (
	// Cancel ...
	Cancel context.CancelFunc
)

// InitMconfig ...
func InitMconfig(mconfig *MConfig) func() {
	ctx, cancelFunc := context.WithCancel(context.Background())
	Cancel = cancelFunc
	store.InitStore(mconfig.StoreType, mconfig.StoreAddress)
	err := store.CurrentMConfigStore.WatchConfigVal(ctx, event.NewMConfigEventCustomer())
	if err != nil {
		log.Fatal(err)
	}
	go event.StartMConfigStoreEventBus(ctx)
	go config.StratMconfigConfigManagement(ctx)
	return EndMconfig()
}

// EndMconfig ...
func EndMconfig() func() {
	return func() {
		Cancel()
	}
}
