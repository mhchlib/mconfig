package cron

import (
	log "github.com/mhchlib/logger"
	"github.com/mhchlib/mconfig/core/config"
	"github.com/mhchlib/mconfig/core/filter"
	"github.com/mhchlib/mconfig/core/store"
	"time"
)

func InitCron() {
	initSyncWithStoreCron(5 * 60 * time.Second)
	if store.CheckNeedSyncData() {
		initSyncWithNodeCron(5 * 60 * time.Second)
	}
}

func initSyncWithStoreCron(t time.Duration) {
	go func() {
		for {
			<-time.After(t)
			syncWithStore()
		}
	}()
}

func initSyncWithNodeCron(t time.Duration) {
	go func() {
		for {
			<-time.After(t)
			syncWithNode()
		}
	}()
}

func syncWithNode() {
	store.SyncOtherMconfigDataCron()
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
