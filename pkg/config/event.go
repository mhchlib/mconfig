package config

import (
	"errors"
	"fmt"
	log "github.com/mhchlib/logger"
	"github.com/mhchlib/mconfig/pkg/client"
	"github.com/mhchlib/mconfig/pkg/event"
	"github.com/mhchlib/mconfig/pkg/mconfig"
)

const EVENT_KEY = "config"

type ConfigEventMetadata struct {
	AppKey    mconfig.Appkey
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
	err = PutConfigToCache(eventMetadata.AppKey, eventMetadata.ConfigKey, eventMetadata.Env, eventMetadata.Val)
	if err != nil {
		log.Error(err)
	}
	clientSet := client.GetOnlineClientSet(eventMetadata.AppKey, eventMetadata.ConfigKey, eventMetadata.Env)
	if clientSet == nil {
		return
	}
	clients := clientSet.GetClients()
	val, err := GetConfigFromCache(eventMetadata.AppKey, eventMetadata.ConfigKey, eventMetadata.Env)
	if err != nil {
		log.Error(err)
		return
	}
	for _, c := range clients {
		err := c.SendMsg(val)
		if err != nil {
			log.Error(err)
			return
		}
	}
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
