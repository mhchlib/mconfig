package cron

import (
	log "github.com/mhchlib/logger"
	"github.com/mhchlib/mconfig/core/config"
	"github.com/mhchlib/mconfig/core/filter"
	"time"
)

// InitCron ...
func InitCron() {
	initSyncWithStoreCron(5 * 60 * time.Second)
}

func initSyncWithStoreCron(t time.Duration) {
	go func() {
		timer := time.NewTimer(t)
		defer timer.Stop()
		for {
			<-timer.C
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
