package pkg

import (
	"context"
	"sync"
)

// Appkey ...
type Appkey string

// EventType ...
type EventType int

// Event_Update ...
var Event_Update EventType = 0

// Event_Delete ...
var Event_Delete EventType = 1

var (
	appConfigStore AppConfigStore
	// Cancel ...
	Cancel context.CancelFunc
)

// ConfigEvent ...
type ConfigEvent struct {
	Key        Appkey
	AppConfigs *AppConfigs
	EventType  EventType
}

// Config ...
type Config struct {
	Schema     string `json:"schema"`
	Config     string `json:"config"`
	CreateTime int64  `json:"create_time"`
	UpdateTime int64  `json:"update_time"`
}

// Configs ...
type Configs struct {
	Configs    ConfigsMap        `json:"configs"`
	Desc       string            `json:"desc"`
	CreateTime int64             `json:"create_time"`
	UpdateTime int64             `json:"update_time"`
	ABFilters  map[string]string `json:"ABFilters"`
}

// AppConfigsMap ...
type AppConfigsMap struct {
	mutex      sync.RWMutex
	AppConfigs *AppConfigs
}

// ConfigsMap ...
type ConfigsMap struct {
	mutex sync.RWMutex
	Entry map[string]*Config `json:"entry"`
}

// AppConfigs ...
type AppConfigs map[string]*Configs
