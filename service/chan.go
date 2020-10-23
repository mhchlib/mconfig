package service

import "sync"

type ClientChanMap struct {
	sync.RWMutex
	m map[ConfigId]map[ClientId]chan interface{}
}

var (
	//用于收到store配置改变事件后通知更新客户端
	configChangeChan chan ConfigId
	clientChanMap    *ClientChanMap
)

func init() {
	configChangeChan = make(chan ConfigId, 10)
	clientChanMap = &ClientChanMap{
		m: make(map[ConfigId]map[ClientId]chan interface{}),
	}
}

func (ch *ClientChanMap) AddClient(clientId ClientId, configId ConfigId, clientMsgChan chan interface{}) {
	ch.Lock()
	defer ch.Unlock()
	v, ok := ch.m[configId]
	if ok == false {
		v = make(map[ClientId]chan interface{})
		ch.m[configId] = v
	}
	msgChan, ok := v[clientId]
	if ok == true {
		close(msgChan)
	}
	v[clientId] = clientMsgChan
}

func (ch *ClientChanMap) RemoveClient(clientId ClientId, configId ConfigId) {
	ch.Lock()
	defer ch.Unlock()
	v, ok := ch.m[configId]
	if ok == false {
		return
	}
	msgChan, ok := v[clientId]
	if ok == true {
		close(msgChan)
		delete(v, clientId)
		return
	}
	return
}

func (ch *ClientChanMap) GetClientsChan(configId ConfigId) []chan interface{} {
	chs := []chan interface{}{}
	ch.RLock()
	v, ok := ch.m[configId]
	if ok == false {
		return nil
	}
	for _, ch := range v {
		chs = append(chs, ch)
	}
	defer ch.RUnlock()
	return chs
}
