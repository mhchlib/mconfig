package pkg

import (
	log "github.com/mhchlib/logger"
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

func (ch *ClientChanMap) AddClient(clientId ClientId, appid Appkey, clientMsgChan chan interface{}) {
	ch.Lock()
	defer ch.Unlock()
	v, ok := ch.m[appid]
	if ok == false {
		v = make(map[ClientId]chan interface{})
		ch.m[appid] = v
	}
	msgChan, ok := v[clientId]
	if ok == true {
		close(msgChan)
	}
	v[clientId] = clientMsgChan
	log.Info("add client: ", clientId, " listen app : ", appid)
}

func (ch *ClientChanMap) RemoveClient(clientId ClientId, appid Appkey) {
	ch.Lock()
	defer ch.Unlock()
	v, ok := ch.m[appid]
	if ok == false {
		return
	}
	msgChan, ok := v[clientId]
	if ok == true {
		close(msgChan)
		delete(v, clientId)
		log.Info("remove client: ", clientId, " listen app ", appid)
		return
	}
	return
}

func (ch *ClientChanMap) GetClientsChan(appid Appkey) []chan interface{} {
	chs := []chan interface{}{}
	ch.RLock()
	v, ok := ch.m[appid]
	if ok == false {
		return nil
	}
	for _, ch := range v {
		chs = append(chs, ch)
	}
	defer ch.RUnlock()
	return chs
}
