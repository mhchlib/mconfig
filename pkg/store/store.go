package store

import (
	"context"
	log "github.com/mhchlib/logger"
	"github.com/mhchlib/mconfig/pkg"
	"github.com/mhchlib/mconfig/pkg/event"
	"github.com/mhchlib/mconfig/pkg/mconfig"
)

// MConfigStore ...
type MConfigStore interface {
	GetConfigVal(appKey mconfig.Appkey, configKey mconfig.ConfigKey, env mconfig.ConfigEnv) (mconfig.ConfigVal, error)
	PutConfigVal(appKey mconfig.Appkey, configKey mconfig.ConfigKey, env mconfig.ConfigEnv, content mconfig.ConfigVal) error
	WatchConfigVal(ctx context.Context, customer event.MConfigEventCustomer) error

	NewAppMetaData(meta mconfig.AppMetaData) error
	NewConfigMetaData(meta mconfig.ConfigMetaData) error
	GetAppConfigs(appKey mconfig.Appkey) ([]mconfig.ConfigMetaData, error)
	UpdateAppMetaData(meta mconfig.AppMetaData) error
	UpdateConfigMetaData(meta mconfig.ConfigMetaData) error
	DeleteApp(appKey mconfig.Appkey) error
	DeleteConfig(appKey mconfig.Appkey, configKey mconfig.ConfigKey) error
	ListAppMetaData(limit int, offset int, filter string) error
}

//CurrentMConfigStore
var CurrentMConfigStore MConfigStore

// InitStore ...
func InitStore(storeType string, storeAddress string) {
	plugin, ok := pkg.StorePluginMap[storeType]
	if !ok {
		log.Fatal("store type: ", storeType, " can not be supported, you can choose: ", pkg.storePluginNames)
	}
	store, err := plugin.Init(storeAddress)
	if err != nil {
		log.Fatal(err)
	}
	CurrentMConfigStore = store
	//测试连接
	log.Info("store init success with", storeType, storeAddress)
}
