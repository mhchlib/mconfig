package event

import (
	"context"
	log "github.com/mhchlib/logger"
)

var eventLoop EventLoop
var management *EventManagement

func init() {
	management = &EventManagement{}
}

func StartEventBus(ctx context.Context) {
	log.Info("start event bus")
	eventLoop := newEventCustomer(LENGTH_MAX_EVENT, management)
	err := eventLoop.handleEvent(ctx)
	if err != nil {
		log.Error(err)
	}
	log.Info("event bus is closed")
}

func AddEvent(event ChangeEvent) error {
	return eventLoop.addEvent(event)
}

func RegisterEventBus(eventDataType EventKey, eventType EventType, handle func(metadata MetaData, changVal interface{})) error {
	return management.registerEvent(eventDataType, eventType, handle)
}

func RegisterMultiEventBus(eventDataType EventKey, eventTypes []EventType, handle func(metadata MetaData, changVal interface{})) error {
	for _, eventType := range eventTypes {
		err := RegisterEventBus(eventDataType, eventType, handle)
		if err != nil {
			return err
		}
	}
	return nil
}
