package pkg

import (
	log "github.com/mhchlib/logger"
	"sync"
)

type ClientChanMap struct {
	sync.RWMutex
	m map[AppId]map[ClientId]chan interface{}
}

var (
	//用于收到store配置改变事件后通知更新客户端
	configChangeChan chan AppId
	clientChanMap    *ClientChanMap
)

func init() {
	configChangeChan = make(chan AppId, 10)
	clientChanMap = &ClientChanMap{
		m: make(map[AppId]map[ClientId]chan interface{}),
	}
}

func (ch *ClientChanMap) AddClient(clientId ClientId, appid AppId, clientMsgChan chan interface{}) {
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
	log.Info("Add Client: ", clientId, " Appid: ", appid)
}

func (ch *ClientChanMap) RemoveClient(clientId ClientId, appid AppId) {
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
		log.Info("Remove Client: ", clientId, " Appid: ", appid)
		return
	}
	return
}

func (ch *ClientChanMap) GetClientsChan(appid AppId) []chan interface{} {
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
