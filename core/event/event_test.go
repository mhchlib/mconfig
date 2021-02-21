package event

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRegisterEventBus(t *testing.T) {
	InitEventBus()
	defer CloseEventBus()
	err := RegisterEventBus("test", Event_Change, func(metadata Metadata) {
	})
	assert.Nil(t, err)
}

func TestAddEvent(t *testing.T) {
	InitEventBus()
	defer CloseEventBus()
	err := RegisterEventBus("test", Event_Change, func(metadata Metadata) {
	})
	assert.Nil(t, err)
	var metadata Metadata
	err = AddEvent(&Event{
		EventDesc: EventDesc{
			EventType: Event_Change,
			EventKey:  "test",
		},
		Metadata: metadata,
	})
	assert.Nil(t, err)
}

func TestAll(t *testing.T) {
	InitEventBus()
	defer CloseEventBus()
	err := RegisterEventBus("test", Event_Change, f)
	assert.Nil(t, err)
	metadata := make(map[string]interface{})
	c := make(chan interface{})
	metadata["key"] = "test key..."
	metadata["chan"] = c
	err = AddEvent(&Event{
		EventDesc: EventDesc{
			EventType: Event_Change,
			EventKey:  "test",
		},
		Metadata: metadata,
	})
	assert.Nil(t, err)
	<-c
}

func f(data Metadata) {
	d := data.(map[string]interface{})
	i := d["chan"]
	c := i.(chan interface{})
	c <- struct{}{}
}
