package pkg

import (
	"context"
	log "github.com/mhchlib/logger"
)

// AppConfigStore ...
type AppConfigStore interface {
	GetAppConfigs(key string) (AppConfigsJSONStr, int64, error)
	PutAppConfigs(key string, value AppConfigsJSONStr) error
	WatchAppConfigs(key string, rev int64, ctx context.Context) (chan *ConfigEvent, error)
	WatchAppConfigsWithPrefix(ctx context.Context) (chan *ConfigEvent, error)
	//...
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

var storePluginMap map[string]*StorePlugin

var storePluginNames []string

// RegisterStorePlugin ...
func RegisterStorePlugin(name string, init func(address string) (AppConfigStore, error)) {
	if storePluginMap == nil {
		storePluginMap = make(map[string]*StorePlugin)
	}
	if storePluginNames == nil {
		storePluginNames = []string{}
	}

	if _, ok := storePluginMap[name]; ok {
		log.Fatal("Repeated  register same name store plugin ...")
	}
	storePluginMap[name] = NewStorePlugin(name, init)
	storePluginNames = append(storePluginNames, name)
}

// InitStore ...
func InitStore(store_type, store_address string) {
	plugin, ok := storePluginMap[store_type]
	if !ok {
		log.Fatal("store type: ", store_type, " can not be supported, you can choose: ", storePluginNames)
	}
	store, err := plugin.Init(store_address)
	if err != nil {
		log.Fatal(err)
	}
	appConfigStore = store
	//测试连接

	log.Info("store init success...")
}
