package pkg

import (
	"context"
	log "github.com/mhchlib/logger"
)

// AppConfigStore ...
type AppConfigStore interface {
	GetAppConfigs(key Appkey) (*AppConfigs, error)
	PutAppConfigs(key Appkey, value *AppConfigs) error
	WatchAppConfigs(ctx context.Context) (chan *ConfigEvent, error)
}

// StorePlugin ...
type StorePlugin struct {
	Name string
	Init func(address string) (AppConfigStore, error)
	//...
}

// NewStorePlugin ...
func NewStorePlugin(name string, init func(address string) (AppConfigStore, error)) *StorePlugin {
	return &StorePlugin{Name: name, Init: init}
}

var StorePluginMap map[string]*StorePlugin

var storePluginNames []string

// RegisterStorePlugin ...
func RegisterStorePlugin(name string, init func(address string) (AppConfigStore, error)) {
	if StorePluginMap == nil {
		StorePluginMap = make(map[string]*StorePlugin)
	}
	if storePluginNames == nil {
		storePluginNames = []string{}
	}

	if _, ok := StorePluginMap[name]; ok {
		log.Fatal("repeated register same name store plugin ...")
	}
	StorePluginMap[name] = NewStorePlugin(name, init)
	storePluginNames = append(storePluginNames, name)
}

// InitStore ...
func InitStore(storeType string, storeAddress string) {
	plugin, ok := StorePluginMap[storeType]
	if !ok {
		log.Fatal("store type: ", storeType, " can not be supported, you can choose: ", storePluginNames)
	}
	store, err := plugin.Init(storeAddress)
	if err != nil {
		log.Fatal(err)
	}
	ConfigStore = store
	//测试连接
	log.Info("store init success with", storeType, storeAddress)
}
