package main

import (
	"fmt"
	"github.com/mhchlib/mconfig/pkg/event"
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
	fmt.Println(data["key"])
	i := data["chan"]
	c := i.(chan interface{})
	c <- struct{}{}
}
