package service

import (
	"sync/atomic"
)

var n int32 = 1000

type ClientId int32

type Client struct {
	Id      ClientId
	MsgChan chan interface{}
}

func NewClient() (*Client, error) {
	id := atomic.AddInt32(&n, 1)
	return &Client{
		Id:      ClientId(id),
		MsgChan: make(chan interface{}, 10),
	}, nil
}
