package client

type ClientMsgBus chan interface{}

const LENGTH_MAX_BUS = 20

func newClientMsgBus() ClientMsgBus {
	bus := make(chan interface{}, LENGTH_MAX_BUS)
	return bus
}
