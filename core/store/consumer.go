package store

import "github.com/mhchlib/mconfig/core/event"

// Consumer ...
type Consumer struct {
}

func newConsumer() *Consumer {
	return &Consumer{}
}

// AddEvent ...
func (consumer Consumer) AddEvent(e *event.Event) error {
	err := event.AddEvent(e)
	return err
}
