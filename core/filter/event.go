package filter

import (
	"errors"
	"fmt"
	log "github.com/mhchlib/logger"
	"github.com/mhchlib/mconfig/core/config"
	"github.com/mhchlib/mconfig/core/event"
	"github.com/mhchlib/mconfig/core/mconfig"
)

// EVENT_KEY ...
const EVENT_KEY event.EventKey = "filter"

// FilterEventMetadata ...
type FilterEventMetadata struct {
	AppKey mconfig.AppKey
	Env    mconfig.ConfigEnv
	Val    *mconfig.StoreVal
}

func initEvent() {
	//config event 2 config center
	err := event.RegisterMultiEventBus(EVENT_KEY, []event.EventType{event.Event_Add, event.Event_Update}, filterChange)
	if err != nil {
		log.Error(err)
	}
	err = event.RegisterEventBus(EVENT_KEY, event.Event_Delete, filterDelete)
	if err != nil {
		log.Error(err)
	}
}

func filterChange(metadata event.Metadata) {
	//sync config change to cache
	eventMetadata, err := parseConfigEventMetadata(metadata)
	if err != nil {
		log.Error(err)
		return
	}
	//avoid use a lot of memory，so here we just put cache what we need
	exist := config.CheckRegisterAppNotifyExist(eventMetadata.AppKey)
	if !exist {
		log.Info("reject put config to cache with app", eventMetadata.AppKey)
		return
	}
	_ = deleteEffectFilterCacheWithAppKey(eventMetadata.AppKey)
	err = PutFilterToCache(eventMetadata.AppKey, eventMetadata.Env, eventMetadata.Val)
	if err != nil {
		log.Error(err)
	}
}

func filterDelete(metadata event.Metadata) {

}

func parseConfigEventMetadata(metadata event.Metadata) (*FilterEventMetadata, error) {
	eventMetadata, ok := metadata.(FilterEventMetadata)
	if !ok {
		return nil, errors.New(fmt.Sprintf("parse config event metadata error, metadata : %+v", metadata))
	}
	return &eventMetadata, nil
}
