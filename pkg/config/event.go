package config

import (
	"errors"
	"fmt"
	log "github.com/mhchlib/logger"
	"github.com/mhchlib/mconfig/pkg/event"
	"github.com/mhchlib/mconfig/pkg/mconfig"
)

const EVENT_KEY event.EventKey = "config"

type ConfigEventMetadata struct {
	AppKey    mconfig.AppKey
	ConfigKey mconfig.ConfigKey
	Env       mconfig.ConfigEnv
	Val       mconfig.ConfigVal
}

func initEvent() {
	//config event 2 config center
	err := event.RegisterMultiEventBus(EVENT_KEY, []event.EventType{event.Event_Add, event.Event_Update}, configChange)
	if err != nil {
		log.Error(err)
	}
	err = event.RegisterEventBus(EVENT_KEY, event.Event_Delete, configDelete)
	if err != nil {
		log.Error(err)
	}
}

func configChange(metadata event.Metadata) {
	//sync config change to cache
	eventMetadata, err := parseConfigEventMetadata(metadata)
	if err != nil {
		log.Error(err)
		return
	}
	//avoid use a lot of memory，so here we just put cache what we need
	//cacheVal,err := GetConfigFromCache(eventMetadata.AppKey, eventMetadata.ConfigKey, eventMetadata.Env)
	//if cacheVal == "" && errors.As(err, &cache.ERROR_CACHE_NOT_FOUND) {
	//	log.Info("reject put config to cache key",eventMetadata.ConfigKey)
	//	return
	//}
	err = PutConfigToCache(eventMetadata.AppKey, eventMetadata.ConfigKey, eventMetadata.Env, eventMetadata.Val)
	if err != nil {
		log.Error(err)
	}

	_ = event.AddEvent(&event.Event{
		EventDesc: event.EventDesc{
			EventType: event.Event_Change,
			EventKey:  mconfig.EVENT_KEY_CLIENT_NOTIFY,
		},
		Metadata: mconfig.ClientNotifyEventMetadata{
			AppKey:    eventMetadata.AppKey,
			ConfigKey: eventMetadata.ConfigKey,
			Env:       eventMetadata.Env,
			Type:      mconfig.Event_Type_Config,
		},
	})
}

func configDelete(metadata event.Metadata) {

}

func parseConfigEventMetadata(metadata event.Metadata) (*ConfigEventMetadata, error) {
	eventMetadata, ok := metadata.(ConfigEventMetadata)
	if !ok {
		return nil, errors.New(fmt.Sprintf("parse config event metadata error, metadata : %+v", metadata))
	}
	return &eventMetadata, nil
}
