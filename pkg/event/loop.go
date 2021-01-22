package event

import (
	"context"
	"errors"
	log "github.com/mhchlib/logger"
)

type EventLoop interface {
	addEvent(event *Event) error
	handleEvent() error
	close()
}

type EventLoopImpl struct {
	management *EventManagement
	eventBus   chan *Event
	cancelFunc func()
}

func newEventCustomer(eventBusLength int, management *EventManagement) EventLoop {
	loop := &EventLoopImpl{
		eventBus:   make(chan *Event, eventBusLength),
		management: management,
	}
	return loop
}

func (e *EventLoopImpl) addEvent(event *Event) error {
	e.eventBus <- event
	return nil
}

func (e *EventLoopImpl) handleEvent() error {
	ctx, cancelFunc := context.WithCancel(context.Background())
	e.cancelFunc = cancelFunc
	for {
		select {
		case <-ctx.Done():
			{
				return nil
			}
		case event, ok := <-e.eventBus:
			if !ok {
				return errors.New("event bus channel is closed")
			}
			err := e.management.handleEvent(event)
			if err != nil {
				log.Error(err)
			}
		}
	}
}

func (e *EventLoopImpl) close() {
	if e.cancelFunc != nil {
		e.cancelFunc()
	}
	close(e.eventBus)
}
