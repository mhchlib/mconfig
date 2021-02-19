package mconfig

import "github.com/mhchlib/mconfig/pkg/event"

const EVENT_KEY_CLIENT_NOTIFY event.EventKey = "client_notify"

type Event_Type int

var (
	Event_Type_Config Event_Type = 0
	Event_Type_Filter Event_Type = 1
)

type ClientNotifyEventMetadata struct {
	AppKey    AppKey
	ConfigKey ConfigKey
	Env       ConfigEnv
	Type      Event_Type
}
