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
	//...
}

// StorePlugin ...
type StorePlugin struct {
	Name string
	Init func() (AppConfigStore, error)
	//...
}

// NewStorePlugin ...
func NewStorePlugin(name string, init func() (AppConfigStore, error)) *StorePlugin {
	return &StorePlugin{Name: name, Init: init}
}

var storePluginMap map[string]*StorePlugin

var storePluginNames []string

// RegisterStorePlugin ...
func RegisterStorePlugin(name string, init func() (AppConfigStore, error)) {
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
func InitStore(storeType string) {
	plugin, ok := storePluginMap[storeType]
	if !ok {
		log.Fatal("store type: ", storeType, " can not be supported, you can choose: ", storePluginNames)
	}
	store, err := plugin.Init()
	if err != nil {
		log.Fatal(err)
	}
	appConfigStore = store
	//测试连接

	log.Info("store init success... with  ", storeType)
}
