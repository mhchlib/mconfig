package event

import (
	log "github.com/mhchlib/logger"
)

var eventLoop EventLoop
var management *EventManagement

func init() {
	management = &EventManagement{}
}

func InitEventBus() {
	log.Info("start event bus")
	eventLoop = newEventCustomer(LENGTH_MAX_EVENT, management)
	go func() {
		err := eventLoop.handleEvent()
		if err != nil {
			log.Error(err)
		}
	}()
}

func CloseEventBus() {
	eventLoop.close()
	log.Info("close event bus")
}

func AddEvent(event *Event) error {
	return eventLoop.addEvent(event)
}

func RegisterEventBus(eventKey EventKey, eventType EventType, handle EventInvoke) error {
	return management.registerEvent(eventKey, eventType, handle)
}

func RegisterMultiEventBus(eventKey EventKey, eventTypes []EventType, handle EventInvoke) error {
	for _, eventType := range eventTypes {
		err := RegisterEventBus(eventKey, eventType, handle)
		if err != nil {
			return err
		}
	}
	return nil
}
