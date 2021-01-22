package config

import (
	log "github.com/mhchlib/logger"
	"github.com/mhchlib/mconfig/pkg/event"
)

const EVENT_KEY = "config"

func init() {
	err := event.RegisterEventBus(EVENT_KEY, event.Event_ADD, configChange)
	if err != nil {
		log.Error(err)
	}
	err = event.RegisterEventBus(EVENT_KEY, event.Event_Delete, configDelete)
	if err != nil {
		log.Error(err)
	}
}

func configChange(metadata event.MetaData, changVal interface{}) {

}

func configDelete(metadata event.MetaData, changVal interface{}) {

}
