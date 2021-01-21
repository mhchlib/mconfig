package pkg

import (
	log "github.com/mhchlib/logger"
	"github.com/mhchlib/mconfig/pkg/store"
)

// StorePlugin ...
type StorePlugin struct {
	Name string
	Init func(address string) (store.MConfigStore, error)
	//...
}

// NewStorePlugin ...
func NewStorePlugin(name string, init func(address string) (store.MConfigStore, error)) *StorePlugin {
	return &StorePlugin{Name: name, Init: init}
}

var StorePluginMap map[string]*StorePlugin

var storePluginNames []string

// RegisterStorePlugin ...
func RegisterStorePlugin(name string, init func(address string) (store.MConfigStore, error)) {
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
