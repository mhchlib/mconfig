package client

import (
	"errors"
	"fmt"
	log "github.com/mhchlib/logger"
	"github.com/mhchlib/mconfig/pkg/config"
	"github.com/mhchlib/mconfig/pkg/event"
	"github.com/mhchlib/mconfig/pkg/mconfig"
)

//avoid to referencing package loop move to package common
//const EVENT_NOTIFY_KEY event.EventKey = "client_notify"

//type Event_Type int
//
//var(
//	Event_Type_Config Event_Type = 0
//	Event_Type_Filter Event_Type = 1
//)
//
//type ClientNotifyEventMetadata struct {
//	AppKey    mconfig.AppKey
//	ConfigKey mconfig.ConfigKey
//	Env       mconfig.ConfigEnv
//	Type      Event_Type
//}

func initEvent() {
	err := event.RegisterMultiEventBus(mconfig.EVENT_KEY_CLIENT_NOTIFY, []event.EventType{event.Event_Change}, notifyClient)
	if err != nil {
		log.Error(err)
	}
}

func notifyClient(metadata event.Metadata) {
	eventMetadata, err := parseClientNotifyEventMetadata(metadata)
	if err != nil {
		log.Error(err)
	}
	switch eventMetadata.Type {
	case mconfig.Event_Type_Config:
		notifyClientConfigChange(eventMetadata.AppKey, eventMetadata.ConfigKey, eventMetadata.Env)
	case mconfig.Event_Type_Filter:
		notifyClientFilterChange(eventMetadata.AppKey)
	default:
		log.Error("not support client notify event type", eventMetadata.Type)
	}
}

func notifyClientConfigChange(appKey mconfig.AppKey, configKey mconfig.ConfigKey, env mconfig.ConfigEnv) {
	clientSet := getOnlineClientSetByConfigRealtion(appKey, configKey, env)
	if clientSet == nil {
		return
	}
	val, err := config.GetConfigFromCache(appKey, configKey, env)
	if err != nil {
		log.Error(err)
		return
	}
	err = clientSet.SendMsg(&mconfig.ConfigChangeNotifyMsg{
		Key: configKey,
		Val: val.Val,
	})
	if err != nil {
		log.Error(err)
	}
}

func notifyClientFilterChange(appKey mconfig.AppKey) {
	clientSet := getOnlineClientSetByAppRealtion(appKey)
	if clientSet == nil {
		return
	}
	err := clientSet.ReCalEffectEnv()
	if err != nil {
		log.Error(err)
	}
}

func parseClientNotifyEventMetadata(metadata event.Metadata) (*mconfig.ClientNotifyEventMetadata, error) {
	eventMetadata, ok := metadata.(mconfig.ClientNotifyEventMetadata)
	if !ok {
		return nil, errors.New(fmt.Sprintf("parse config event metadata error, metadata : %+v", metadata))
	}
	return &eventMetadata, nil
}
