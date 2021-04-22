package store

import (
	log "github.com/mhchlib/logger"
)

// StoreMode ...
type StoreMode string

//PluginInitFunc ...
type PluginInitFunc func(address string) (MConfigStore, error)

//PluginCloseFunc ...
type PluginGracefulStopFunc func() error

const (
	// MODE_SHARE ...
	MODE_SHARE StoreMode = "share"
	// MODE_LOCAL ...
	MODE_LOCAL StoreMode = "local"
)

// StorePlugin ...
type StorePlugin struct {
	Name         string
	Mode         StoreMode
	Init         PluginInitFunc
	GracefulStop PluginGracefulStopFunc
	//...
}

// NewStorePlugin ...
func NewStorePlugin(name string, mode StoreMode, init PluginInitFunc, gracefulStop PluginGracefulStopFunc) *StorePlugin {
	return &StorePlugin{Name: name, Mode: mode, Init: init, GracefulStop: gracefulStop}
}

var storePluginMap map[string]*StorePlugin

var storePluginNames []string

// RegisterStorePlugin ...
func RegisterStorePlugin(name string, mode StoreMode, init PluginInitFunc, gracefulStop PluginGracefulStopFunc) {
	if storePluginMap == nil {
		storePluginMap = make(map[string]*StorePlugin)
	}
	if storePluginNames == nil {
		storePluginNames = []string{}
	}

	if _, ok := storePluginMap[name]; ok {
		log.Fatal("repeated register same name store plugin ...")
	}
	storePluginMap[name] = NewStorePlugin(name, mode, init, gracefulStop)
	storePluginNames = append(storePluginNames, name)
}
