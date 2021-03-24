package store

import (
	log "github.com/mhchlib/logger"
)

// StoreMode ...
type StoreMode string

const (
	// MODE_SHARE ...
	MODE_SHARE StoreMode = "share"
	// MODE_LOCAL ...
	MODE_LOCAL StoreMode = "local"
)

// StorePlugin ...
type StorePlugin struct {
	Name string
	Mode StoreMode
	Init func(address string) (MConfigStore, error)
	//...
}

// NewStorePlugin ...
func NewStorePlugin(name string, mode StoreMode, init func(address string) (MConfigStore, error)) *StorePlugin {
	return &StorePlugin{Name: name, Mode: mode, Init: init}
}

var storePluginMap map[string]*StorePlugin

var storePluginNames []string

// RegisterStorePlugin ...
func RegisterStorePlugin(name string, mode StoreMode, init func(address string) (MConfigStore, error)) {
	if storePluginMap == nil {
		storePluginMap = make(map[string]*StorePlugin)
	}
	if storePluginNames == nil {
		storePluginNames = []string{}
	}

	if _, ok := storePluginMap[name]; ok {
		log.Fatal("repeated register same name store plugin ...")
	}
	storePluginMap[name] = NewStorePlugin(name, mode, init)
	storePluginNames = append(storePluginNames, name)
}
