package event

import (
	"context"
	"errors"
	log "github.com/mhchlib/logger"
)

type EventLoop interface {
	addEvent(event ChangeEvent) error
	handleEvent(ctx context.Context) error
}

type EventLoopImpl struct {
	management *EventManagement
	eventBus   chan ChangeEvent
}

func newEventCustomer(eventBusLength int, management *EventManagement) EventLoop {
	loop := &EventLoopImpl{
		eventBus:   make(chan ChangeEvent, eventBusLength),
		management: management,
	}
	return loop
}

func (e EventLoopImpl) addEvent(event ChangeEvent) error {
	e.eventBus <- event
	return nil
}

func (e EventLoopImpl) handleEvent(ctx context.Context) error {
	for {
		select {
		case event, ok := <-e.eventBus:
			if !ok {
				return errors.New("event bus channel is closed")
			}
			err := e.management.handleEvent(event)
			if err != nil {
				log.Error(err)
			}
		case <-ctx.Done():
			return nil
		}
	}
}
