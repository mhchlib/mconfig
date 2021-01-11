package pkg

import (
	"context"
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
	ConfigStore AppConfigStore
	// Cancel ...
	Cancel context.CancelFunc
)

// ConfigEvent ...
type ConfigEvent struct {
	Key        Appkey
	AppConfigs *AppConfigs
	EventType  EventType
}
