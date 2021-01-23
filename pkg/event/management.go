package event

import (
	"errors"
	"sync"
)

type EventType string

var Event_Change EventType = "event_change"
var Event_Add EventType = "event_add"
var Event_Update EventType = "event_update"
var Event_Delete EventType = "event_delete"

type EventKey string

type EventInvoke func(metadata Metadata)

type Metadata interface{}

type EventManagement struct {
	sync.RWMutex
	record map[EventDesc]EventInvoke
}

type EventDesc struct {
	EventType EventType
	EventKey  EventKey
}

type Event struct {
	EventDesc EventDesc
	Metadata  Metadata
}

const LENGTH_MAX_EVENT = 20

func (management *EventManagement) registerEvent(eventKey EventKey, eventType EventType, handle EventInvoke) error {
	management.Lock()
	defer management.Unlock()
	eventDesc := &EventDesc{
		EventType: eventType,
		EventKey:  eventKey,
	}
	if management.record == nil {
		management.record = make(map[EventDesc]EventInvoke)
	}
	management.record[*eventDesc] = handle
	return nil
}

func (management *EventManagement) getEventInvoke(desc EventDesc) (EventInvoke, error) {
	management.RLock()
	defer management.RUnlock()
	if management.record == nil {
		return nil, errors.New("not found registed invoke method")
	}
	invoke, ok := management.record[desc]
	if !ok {
		return nil, errors.New("not found registed invoke method")
	}
	return invoke, nil
}

func (management *EventManagement) handleEvent(event *Event) error {
	eventDesc := event.EventDesc
	invoke, err := management.getEventInvoke(eventDesc)
	if err != nil {
		return err
	}
	go func() {
		invoke(event.Metadata)
	}()
	return nil
}
