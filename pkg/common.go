package pkg

import (
	"context"
	"sync"
)

type Appkey string
type AppConfigsJSONStr string
type EventType int

var Event_Update EventType = 0
var Event_Delete EventType = 1

var (
	appConfigStore AppConfigStore
	Cancel         context.CancelFunc
)

type ConfigEvent struct {
	Key       Appkey
	Value     AppConfigsJSONStr
	EventType EventType
}

type Config struct {
	Schema     string `json:"schema"`
	Config     string `json:"config"`
	CreateTime int64  `json:"create_time"`
	UpdateTime int64  `json:"update_time"`
}

type Configs struct {
	Configs    ConfigsMap        `json:"configs"`
	Desc       string            `json:"desc"`
	CreateTime int64             `json:"create_time"`
	UpdateTime int64             `json:"update_time"`
	ABFilters  map[string]string `json:"ABFilters"`
}

type AppConfigsMap struct {
	mutex sync.RWMutex
	AppConfigs
}

type ConfigsMap struct {
	mutex sync.RWMutex
	Entry map[string]*Config `json:"entry"`
}

type AppConfigs map[string]*Configs
