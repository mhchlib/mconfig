package main

import (
	"fmt"
	log "github.com/mhchlib/logger"
	"github.com/mhchlib/mconfig/core/event"
)

func main() {
	event.InitEventBus()
	defer event.CloseEventBus()

	err := event.RegisterEventBus("example", event.Event_Change, f)
	if err != nil {
		fmt.Errorf(err.Error())
		return
	}
	metadata := make(map[string]interface{})
	c := make(chan interface{})
	metadata["key"] = "example key..."
	metadata["chan"] = c
	err = event.AddEvent(&event.Event{
		EventDesc: event.EventDesc{
			EventType: event.Event_Change,
			EventKey:  "example",
		},
		Metadata: metadata,
	})
	if err != nil {
		fmt.Errorf(err.Error())
		return
	}
	<-c
}

func f(data event.Metadata) {
	d := data.(map[string]interface{})
	log.Info(d["key"])
	i := d["chan"]
	c := i.(chan interface{})
	c <- struct{}{}
}
