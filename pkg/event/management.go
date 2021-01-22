package event

import (
	"errors"
	"sync"
)

type EventType int

var Event_Change EventType = 0
var Event_ADD EventType = 1
var Event_Update EventType = 2
var Event_Delete EventType = 3

type EventKey string

type MetaData map[string]interface{}

type ChangeVal interface{}

type EventManagement struct {
	sync.RWMutex
	record map[EventDesc]func(metaData MetaData, changVal interface{})
}

type EventDesc struct {
	EventType EventType
	EventKey  EventKey
}

type ChangeEvent struct {
	EventDesc EventDesc
	Metadata  MetaData
	val       ChangeVal
}

const LENGTH_MAX_EVENT = 20

func (management *EventManagement) registerEvent(eventKey EventKey, eventType EventType, handle func(metaData MetaData, changVal interface{})) error {
	management.Lock()
	defer management.Unlock()
	eventDesc := &EventDesc{
		EventType: eventType,
		EventKey:  eventKey,
	}
	if management.record == nil {
		management.record = make(map[EventDesc]func(metaData MetaData, changVal interface{}))
	}
	management.record[*eventDesc] = handle
	return nil
}

func (management *EventManagement) getEventInvoke(desc EventDesc) (func(metaData MetaData, changVal interface{}), error) {
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

func (management *EventManagement) handleEvent(event ChangeEvent) error {
	eventDesc := event.EventDesc
	invoke, err := management.getEventInvoke(eventDesc)
	if err != nil {
		return err
	}
	go func() {
		invoke(event.Metadata, event.val)
	}()
	return nil
}
