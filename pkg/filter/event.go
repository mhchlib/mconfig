package filter

import (
	"errors"
	"fmt"
	log "github.com/mhchlib/logger"
	"github.com/mhchlib/mconfig/pkg/event"
	"github.com/mhchlib/mconfig/pkg/mconfig"
)

const EVENT_KEY event.EventKey = "filter"

type FilterEventMetadata struct {
	AppKey mconfig.AppKey
	Env    mconfig.ConfigEnv
	Val    *mconfig.FilterStoreVal
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
	//cacheVal,_ := GetFilterFromCache(eventMetadata.AppKey)
	//if cacheVal == nil {
	//	log.Info("reject put config to cache key",eventMetadata.Env)
	//	return
	//}
	err = PutFilterToCache(eventMetadata.AppKey, eventMetadata.Env, eventMetadata.Val)
	if err != nil {
		log.Error(err)
	}
	_ = event.AddEvent(&event.Event{
		EventDesc: event.EventDesc{
			EventType: event.Event_Change,
			EventKey:  mconfig.EVENT_KEY_CLIENT_NOTIFY,
		},
		Metadata: mconfig.ClientNotifyEventMetadata{
			AppKey: eventMetadata.AppKey,
			Type:   mconfig.Event_Type_Filter,
		},
	})
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
