package client

import (
	"sync"
)

// ClientSendFunc ...
type ClientSendFunc func(interface{}) error

// ClientRecvFunc ...
type ClientRecvFunc func(c *Client) error

// ClientMsgBus ...
type ClientMsgBus struct {
	sync.RWMutex
	SendFunc      ClientSendFunc
	RecvFunc      ClientRecvFunc
	willBeRemoved bool
}

func (bus *ClientMsgBus) sendMsg(data interface{}) error {
	err := bus.SendFunc(data)
	if err != nil {
		return err
	}
	return nil
}

// Close ...
func (bus *ClientMsgBus) Close() {
	bus.willBeRemoved = true
}

func newClientMsgBus(send ClientSendFunc, recv ClientRecvFunc) *ClientMsgBus {
	bus := &ClientMsgBus{
		SendFunc: send,
		RecvFunc: recv,
	}
	return bus
}
