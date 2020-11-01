package service

import (
	"context"
	"github.com/mhchlib/mconfig-api/api/v1/common"
)

type AppId string
type ConfigJSONStr string

type EventType int

var Event_Update EventType = 0
var Event_Delete EventType = 1

type ConfigEvent struct {
	Key       AppId
	Value     ConfigJSONStr
	EventType EventType
}

var (
	ConfigStore Store
	Cancel      context.CancelFunc
)

type ConfigEntity struct {
	Id         string
	Schema     string
	Config     string
	Status     common.ConfigStatus
	Desc       string
	CreateTime int64
	UpdateTime int64
}
