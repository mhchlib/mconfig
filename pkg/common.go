package pkg

import (
	"context"
)

type AppId string
type AppConfigsJSONStr string
type EventType int

var Event_Update EventType = 0
var Event_Delete EventType = 1

var (
	appConfigStore AppConfigStore
	Cancel         context.CancelFunc
)

type ConfigEvent struct {
	Key       AppId
	Value     AppConfigsJSONStr
	EventType EventType
}

type Config struct {
	Schema     string
	Config     string
	CreateTime int64 `json:"create_time"`
	UpdateTime int64 `json:"update_time"`
}

type Configs struct {
	Configs    map[string]Config
	Desc       string
	CreateTime int64             `json:"create_time"`
	UpdateTime int64             `json:"update_time"`
	ABFilters  map[string]string `json:"ABFilters"`
}

type AppConfigs map[string]Configs
