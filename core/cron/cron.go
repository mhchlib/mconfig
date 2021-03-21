package cron

import (
	log "github.com/mhchlib/logger"
	"github.com/mhchlib/mconfig/core/config"
	"github.com/mhchlib/mconfig/core/filter"
	"time"
)

func InitCron() {
	initSyncWithStoreCron(5 * 60 * time.Second)
}

func initSyncWithStoreCron(t time.Duration) {
	go func() {
		for {
			<-time.After(t)
			syncWithStore()
		}
	}()
}

func syncWithStore() {
	err := config.CheckCacheUpToDateWithStore()
	if err != nil {
		log.Error("config sync with store error")
	}
	err = filter.CheckCacheUpToDateWithStore()
	if err != nil {
		log.Error("filter sync with store error")
	}
}
