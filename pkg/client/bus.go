package client

import (
	"sync"
)

type ClientSendFunc func(interface{}) error
type ClientRecvFunc func(c *Client) error

type ClientMsgBus struct {
	sync.RWMutex
	SendFunc      ClientSendFunc
	RecvFunc      ClientRecvFunc
	willBeRemoved bool
}

func (bus *ClientMsgBus) sendMsg(data interface{}) error {
	//get config from cache and send data
	err := bus.SendFunc(data)
	if err != nil {
		return err
	}
	return nil
}

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
