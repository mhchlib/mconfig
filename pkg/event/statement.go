package event

import (
	"context"
	"github.com/mhchlib/mconfig/pkg/mconfig"
)

// EventType ...
type EventType int

// Event_Update ...
var Event_Update EventType = 0

// Event_Delete ...
var Event_Delete EventType = 1

// ConfigEvent ...
type MConfigValEvent struct {
	EventType EventType
	appKey    mconfig.Appkey
	configKey mconfig.ConfigKey
	val       mconfig.ConfigVal
}

type MConfigEventCustomer interface {
	AddEvent(event MConfigValEvent) error
	handleEvent(ctx context.Context)
}

var customer MConfigEventCustomer

const LENGTH_MAX_EVENT = 20
