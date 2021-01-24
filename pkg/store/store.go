package store

import (
	log "github.com/mhchlib/logger"
	"github.com/mhchlib/mconfig/pkg/mconfig"
)

// MConfigStore ...
type MConfigStore interface {
	GetConfigVal(appKey mconfig.Appkey, configKey mconfig.ConfigKey, env mconfig.ConfigEnv) (mconfig.ConfigVal, error)
	PutConfigVal(appKey mconfig.Appkey, configKey mconfig.ConfigKey, env mconfig.ConfigEnv, content mconfig.ConfigVal) error
	WatchDynamicVal(customer *Consumer) error

	NewAppMetaData(meta *mconfig.AppMetaData) error
	NewConfigMetaData(meta *mconfig.ConfigMetaData) error
	GetAppConfigs(appKey mconfig.Appkey) ([]*mconfig.ConfigMetaData, error)
	UpdateAppMetaData(meta *mconfig.AppMetaData) error
	UpdateConfigMetaData(meta *mconfig.ConfigMetaData) error
	DeleteApp(appKey mconfig.Appkey) error
	DeleteConfig(appKey mconfig.Appkey, configKey mconfig.ConfigKey) error
	ListAppMetaData(limit int, offset int, filter string) error

	Close() error
}

//CurrentMConfigStore
var currentMConfigStore MConfigStore

// InitStore ...
func InitStore(storeType string, storeAddress string) {
	plugin, ok := StorePluginMap[storeType]
	if !ok {
		log.Fatal("store type:", storeType, "does not be supported, you can choose:", storePluginNames)
	}
	store, err := plugin.Init(storeAddress)
	if err != nil {
		log.Fatal(err)
	}
	currentMConfigStore = store
	log.Info("store init success with", storeType, storeAddress)
	go func() {
		err = currentMConfigStore.WatchDynamicVal(newConsumer())
		if err != nil {
			log.Error(err)
		}
	}()
}
